package exporter

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func TestPerfExporter(t *testing.T) {
	paranoidBytes, err := ioutil.ReadFile("/proc/sys/kernel/perf_event_paranoid")
	if err != nil {
		t.Skip("Procfs not mounted, skipping perf tests")
	}
	paranoidStr := strings.Replace(string(paranoidBytes), "\n", "", -1)
	paranoid, err := strconv.Atoi(paranoidStr)
	if err != nil {
		t.Fatalf("Expected perf_event_paranoid to be an int, got: %s", paranoidStr)
	}
	if paranoid >= 1 {
		t.Skip("Skipping perf tests, set perf_event_paranoid to 0")
	}
	t.Skip()

	config := viper.GetViper()
	config.SetDefault("raw_syscalls.events", []string{"sys_enter", "sys_exit"})

	collector, err := NewPerfCollector(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := collector.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
	desc := make(chan *prometheus.Desc)
	collector.Describe(desc)
}

func TestPerfExporterCollect(t *testing.T) {
	t.Skip()
	config := viper.GetViper()
	config.SetDefault("raw_syscalls.events", []string{"sys_enter", "sys_exit"})

	collector, err := NewPerfCollector(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := collector.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		if err := collector.Start(); err != nil {
			t.Fatal(err)
		}
	}()
	ch := make(chan prometheus.Metric)
	collector.Collect(ch)
}
