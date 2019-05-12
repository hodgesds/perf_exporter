# perf_exporter
`perf_exporter` is a Prometheus exporter that exposes metrics from the perf
subsystem in Linux. It can read any kernel tracepoints and expose them as
Prometheus compatible metrics.

## Configuration
The configuration format allows you to specific specific profilers at the
subsytem level. For each subsytem individual events can be configured. Note
that configuring a subsystem event for a specific processor isn't supported as
of now. To find available events for your system you can use the perf tooling
(i.e. `perf list`) or you can read directly from tracefs `available_events` in
combination with the `tools/tracepoint2yaml` script.  Here is a rough example
of a configuration file to get started (note this is ***highly*** system
specific).

```
kmem:
  events:
    - mm_page_alloc_extfrag
    - mm_page_pcpu_drain
    - mm_page_alloc_zone_locked
    - mm_page_alloc
    - mm_page_free_batched
    - mm_page_free
    - kmem_cache_free
    - kfree
    - kmem_cache_alloc_node
    - kmalloc_node
    - kmem_cache_alloc
    - kmalloc
net:
  events:
    - netif_rx_ni_entry
    - netif_rx_entry
    - netif_receive_skb_entry
    - napi_gro_receive_entry
    - napi_gro_frags_entry
    - netif_rx
    - netif_receive_skb
    - net_dev_queue
    - net_dev_xmit
    - net_dev_start_xmit
```

**Note** that the proper value for `perf_event_paranoid` should be set, in this
case it should be set to **0** becuase the exporter runs on all processors. For
more info see `man perf_event_open`.

## Building
This repo uses make for the build system, to build the binary just type `make`.
It is assumed that you are using go 1.11+.

## Example
Here is an example of some of the events that can be exposed:
![](https://github.com/hodgesds/dev_pics/blob/master/events.png)

## FAQ
- How is perf being used? You may want to see this
  [library](https://github.com/hodgesds/perf-utils) which is where most of the
  perf related utilities are.
- Should I use this in production? Probably not yet, this is pretty experimental software.
- I don't see values for my perf events, is the collector broken? This is
  difficult to debug due to a large number of factors at play. Everything from
  the way your kernel was configured to debugfs mount points can cause an
  issue, please file an issue so that datapoints can be collected.
- Is there a max number of events that can be profiled? Yes, it is dependent on
  your kernel configuration, originally there was a `--yolo` flag to trace
  everything but that didn't work out so well.
