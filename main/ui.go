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

package main

import (
	"flag"

	"github.com/square/metrics/api/backend"
	"github.com/square/metrics/api/backend/blueflood"
	"github.com/square/metrics/main/common"
	"github.com/square/metrics/query"
	"github.com/square/metrics/ui"
)

func main() {
	flag.Parse()
	common.SetupLogger()

	apiInstance := common.NewAPI()
	bluefloodConfig := blueflood.BluefloodClientConfig{
		BaseUrl:  *common.BluefloodUrl,
		TenantId: *common.BluefloodTenantId,
	}
	blueflood := blueflood.NewBlueflood(bluefloodConfig)
	backend := backend.NewSequentialMultiBackend(blueflood)
	ui.Main(query.ExecutionContext{backend, apiInstance, 1000})
}
