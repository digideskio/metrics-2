// Copyright 2015 Square Inc.
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

// Package api holds common data types and public interface exposed by the indexer library.
package api

import "github.com/square/metrics/inspect"

type MetricMetadataAPIContext struct {
	Profiler *inspect.Profiler
}

type MetricMetadataAPI interface {
	// AddMetric adds the metric to the system.
	AddMetric(metric TaggedMetric, context MetricMetadataAPIContext) error
	// Bulk metrics addition
	AddMetrics(metric []TaggedMetric, context MetricMetadataAPIContext) error
	// For a given MetricKey, retrieve all the tagsets associated with it.
	GetAllTags(metricKey MetricKey, context MetricMetadataAPIContext) ([]TagSet, error)
	// GetAllMetrics returns all metrics managed by the system.
	GetAllMetrics(context MetricMetadataAPIContext) ([]MetricKey, error)
	// For a given tag key-value pair, obtain the list of all the MetricKeys
	// associated with them.
	GetMetricsForTag(tagKey, tagValue string, context MetricMetadataAPIContext) ([]MetricKey, error)
	// CheckHealthy checks if this MetricMetadataAPI is healthy, returning a possible error
	CheckHealthy() error
}
