// Copyright © 2019 Daniel Hodges <hodges.daniel.scott@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter

import (
	"fmt"
	"runtime"
	"strings"

	perf "github.com/hodgesds/perf-utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

const (
	namespace            = "perf"
	exporter             = "perf_exporter"
	eventConfigNamespace = "events"
)

// PerfCollector is an interface that embeds the prometheus Collector interface
// as well as the perf GroupProfiler.
type PerfCollector interface {
	prometheus.Collector
	Start() error
	Stop() error
	Reset() error
}

type perfCollector struct {
	profilers       map[int]perf.GroupProfiler
	descs           map[string]map[string]*prometheus.Desc
	collectionOrder []string
}

// NewPerfCollector implements collectorConstructor.
func NewPerfCollector(config *viper.Viper) (PerfCollector, error) {
	eventAttrs := []unix.PerfEventAttr{}
	descs := map[string]map[string]*prometheus.Desc{}
	profilers := map[int]perf.GroupProfiler{}
	collectionOrder := []string{}

	for _, subsystem := range config.AllKeys() {
		subsystem = strings.Replace(subsystem, ".events", "", -1)
		subConfig := config.Sub(subsystem)
		if subConfig == nil {
			continue
		}
		subsystemEvents := subConfig.GetStringSlice(eventConfigNamespace)
		for _, event := range subsystemEvents {
			eventAttr, err := perf.TracepointEventAttr(subsystem, event)
			if err != nil {
				return nil, fmt.Errorf(
					"Unabled to configure PerfEventAttr: %s %s %v",
					subsystem,
					event,
					err,
				)
			}
			collectionOrder = append(collectionOrder, fmt.Sprintf(
				"%s:%s",
				subsystem,
				event,
			))
			eventAttrs = append(eventAttrs, *eventAttr)
			if _, ok := descs[subsystem]; !ok {
				descs[subsystem] = map[string]*prometheus.Desc{}
			}
			descs[subsystem][event] = prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					event,
				),
				fmt.Sprintf(
					"perf event %s %s",
					subsystem,
					event,
				),
				[]string{"cpu"},
				nil,
			)
		}
	}

	for cpu := 0; cpu < runtime.NumCPU(); cpu++ {
		profiler, err := perf.NewGroupProfiler(-1, cpu, 0, eventAttrs...)
		if err != nil {
			return nil, fmt.Errorf(
				"Unable to configure GroupProfiler: %v", err)
		}
		profilers[cpu] = profiler
	}

	return &perfCollector{
		profilers:       profilers,
		descs:           descs,
		collectionOrder: collectionOrder,
	}, nil
}

// Start starts the profilers.
func (c *perfCollector) Start() error {
	for _, profiler := range c.profilers {
		if err := profiler.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Stop stops the profilers.
func (c *perfCollector) Stop() error {
	for _, profiler := range c.profilers {
		if err := profiler.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// Reset resets the profilers.
func (c *perfCollector) Reset() error {
	for _, profiler := range c.profilers {
		if err := profiler.Reset(); err != nil {
			return err
		}
	}
	return nil
}

// Describe implements the prometheus.Collector interface.
func (c *perfCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, subsystemDescs := range c.descs {
		for _, desc := range subsystemDescs {
			ch <- desc
		}
	}
}

// Collect implements prometheus.Collector interface.
func (c *perfCollector) Collect(ch chan<- prometheus.Metric) {
	for cpu := range c.profilers {
		c.collect(cpu, ch)
	}
}

func (c *perfCollector) collect(cpu int, ch chan<- prometheus.Metric) {
	cpuStr := fmt.Sprintf("%d", cpu)
	profiler := c.profilers[cpu]
	p, err := profiler.Profile()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	for i, value := range p.Values {
		// get the Desc from the ordered group value.
		descKey := c.collectionOrder[i]
		descKeySlice := strings.Split(descKey, ":")
		ch <- prometheus.MustNewConstMetric(
			c.descs[descKeySlice[0]][descKeySlice[1]],
			prometheus.CounterValue,
			float64(value),
			cpuStr,
		)
	}
}
