// Copyright 2015 - 2016 Square Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cached

import (
	"sync"
	"time"

	"github.com/square/metrics/api"
	"github.com/square/metrics/metric_metadata"
	"github.com/square/metrics/util"
)

// BackgroundAPI is a MetadataAPI that also supports background cache updates.
type BackgroundAPI interface {
	metadata.MetricAPI
	// GetBackgroundAction returns a function to be called to execute a background cache update.
	GetBackgroundAction() func(metadata.Context) error
	// CurrentLiveRequests returns the number of requests currently in the queue
	CurrentLiveRequests() int
	// MaximumLiveRequests returns the maximum number of requests that can be in the queue
	MaximumLiveRequests() int
}

// metricMetadataAPI caches some of the metadata associated with the API to reduce latency.
// However, it does not reduce total QPS: whenever it reads from the cache, it performs an update
// in the background by launching a new goroutine.
type metricMetadataAPI struct {
	metricMetadataAPI metadata.MetricAPI // The internal MetricAPI that performs the actual queries.
	clock             util.Clock         // Here so we can mock out in tests

	// Cached items
	getAllTagsCache      map[api.MetricKey]*TagSetList // The cache of metric -> tags
	getAllTagsCacheMutex sync.Mutex                    // Mutex for getAllTagsCache

	// Cache Config
	freshness  time.Duration // How long until cache entries become stale
	timeToLive time.Duration // How long until cache entries become expired

	// Queue
	backgroundQueue chan func(metadata.Context) error // A channel that holds background requests.
	queueMutex      sync.Mutex                        // Synchronizing mutex for the queue
}

// metricUpdateAPI is a wrapper for when the underlying metadata.MetricAPI is also a metadata.MetricUpdateAPI.
type metricUpdateAPI struct {
	metricMetadataAPI
}

func (c *metricMetadataAPI) AddMetric(metric api.TaggedMetric, context metadata.Context) error {
	return c.metricMetadataAPI.(metadata.MetricUpdateAPI).AddMetric(metric, context)
}

func (c *metricMetadataAPI) AddMetrics(metrics []api.TaggedMetric, context metadata.Context) error {
	return c.metricMetadataAPI.(metadata.MetricUpdateAPI).AddMetrics(metrics, context)
}

// Config stores data needed to instantiate a CachedMetricMetadataAPI.
type Config struct {
	Freshness    time.Duration
	RequestLimit int
	TimeToLive   time.Duration
}

// TagSetList is an item in the cache.
type TagSetList struct {
	TagSets []api.TagSet // The tagsets for this metric
	Err     error
	Expiry  time.Time // The time at which the cache entry expires
	Stale   time.Time // The time at which the cache entry becomes stale

	Enqueued bool

	sync.Mutex // Synchronizing mutex
}

// NewMetricMetadataAPI creates a cached API given configuration and an underlying API object.
func NewMetricMetadataAPI(apiInstance metadata.MetricAPI, config Config) BackgroundAPI {
	requests := make(chan func(metadata.Context) error, config.RequestLimit)
	if config.Freshness == 0 {
		config.Freshness = config.TimeToLive
	}
	result := metricMetadataAPI{
		metricMetadataAPI: apiInstance,
		clock:             util.RealClock{},
		getAllTagsCache:   map[api.MetricKey]*TagSetList{},
		freshness:         config.Freshness,
		timeToLive:        config.TimeToLive,
		backgroundQueue:   requests,
	}
	if _, ok := apiInstance.(metadata.MetricUpdateAPI); ok {
		return &metricUpdateAPI{result}
	}
	return &result
}

// GetBackgroundAction is a blocking method that runs one queued cache update.
// It will block until an update is available.
func (c *metricMetadataAPI) GetBackgroundAction() func(metadata.Context) error {
	return <-c.backgroundQueue
}

// GetAllMetrics waits for a slot to be open, then queries the underlying API.
func (c *metricMetadataAPI) GetAllMetrics(context metadata.Context) ([]api.MetricKey, error) {
	return c.metricMetadataAPI.GetAllMetrics(context)
}

// GetMetricsForTag wwaits for a slot to be open, then queries the underlying API.
func (c *metricMetadataAPI) GetMetricsForTag(tagKey, tagValue string, context metadata.Context) ([]api.MetricKey, error) {
	return c.metricMetadataAPI.GetMetricsForTag(tagKey, tagValue, context)
}

// CheckHealthy checks if the underlying MetricAPI is healthy
func (c *metricMetadataAPI) CheckHealthy() error {
	return c.metricMetadataAPI.CheckHealthy()
}

// GetAllTags uses the cache to serve tag data for the given metric.
// If the cache entry is missing or out of date, it uses the results of a query
// to the underlying API to return to the caller. Even if the cache entry is
// up-to-date, this method may enqueue a background request to the underlying API
// to keep the cache fresh.
func (c *metricMetadataAPI) GetAllTags(metricKey api.MetricKey, context metadata.Context) ([]api.TagSet, error) {
	defer context.Profiler.Record("CachedMetricMetadataAPI_GetAllTags")()

	c.getAllTagsCacheMutex.Lock()
	item, ok := c.getAllTagsCache[metricKey]
	if !ok {
		// Create and store the item in the cache.
		item = &TagSetList{}
		c.getAllTagsCache[metricKey] = item
	}
	c.getAllTagsCacheMutex.Unlock()

	item.Lock()
	defer item.Unlock()
	// For the rest of this function, we have exclusive access to the item.

	if item.Expiry.IsZero() || item.Expiry.Before(c.clock.Now()) || item.Err != nil {
		// The item is expired (not just stale).
		// The lock exclusion means everyone will wait until I'm done updating.
		item.TagSets, item.Err = c.metricMetadataAPI.GetAllTags(metricKey, context)
		item.Expiry = c.clock.Now().Add(c.timeToLive)
		item.Stale = c.clock.Now().Add(c.freshness)
		return item.TagSets, item.Err
	}

	if item.Stale.Before(c.clock.Now()) {
		// The item is stale, but not expired.
		// If it hasn't been enqueued, enqueue it (if we can).
		if !item.Enqueued {
			item.Enqueued = true
			select {
			case c.backgroundQueue <- func(context metadata.Context) error {
				// It is guaranteed that the item has been enqueued at this point;
				// channel send happens after item.Enqueued becomes true.
				// Acquire the lock here, and repeat.
				item.Lock() // acquire the item
				defer item.Unlock()
				// NOTE: if a background function is removed from the queue, but never called, the item will be "enqueued" forever.
				// There is no effective way to solve this problem, other than demanding that callers do eventually call background functions.

				item.Enqueued = false
				if !item.Stale.Before(c.clock.Now()) || item.Err != nil {
					return nil // nothing to update, it's not stale (because it was updated by a calling goroutine before it left the background queue)
				}
				item.TagSets, item.Err = c.metricMetadataAPI.GetAllTags(metricKey, context)
				item.Expiry = c.clock.Now().Add(c.timeToLive)
				item.Stale = c.clock.Now().Add(c.freshness)
				return item.Err
			}:
			// nothing to do
			default:
				item.Enqueued = false
				// We couldn't actually add it to the queue.
			}
		}
	}

	return item.TagSets, nil
}

// CurrentLiveRequests returns the number of requests currently in the queue
func (c *metricMetadataAPI) CurrentLiveRequests() int {
	return len(c.backgroundQueue)
}

// MaximumLiveRequests returns the maximum number of requests that can be in the queue
func (c *metricMetadataAPI) MaximumLiveRequests() int {
	return cap(c.backgroundQueue)
}
