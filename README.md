# perf_exporter
`perf_exporter` is a Prometheus exporter that exposes metrics from the perf
subsytem in Linux.

# Configuration
The configuration format allows you to specific specific profilers at the
subsytem level. For each subsytem individual events can be configured. Note
that configuring a subsystem event for a specific processor isn't supported as
of now.

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

# Building
