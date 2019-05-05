package exporter

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
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
}
