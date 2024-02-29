# Mirror Proxy

I wanna create a proxy setup that allow me asses different proxy strategies withing a real time scenario and produce data as output to benchmark those solutions using a Load Balancer to mirroring traffic to each of the proposes solutions and fetch data from some source in real time or historically as we which.

## Proxies

- [ ] (Go) Low level TCP proxy
- [ ] (Go) HTTP simple proxy
- [ ] (Go) HTTP reverse proxy [httputil](https://pkg.go.dev/net/http/httputil#NewSingleHostReverseProxy) implementation
- [ ] (Go) fast HTTP proxy
- [ ] (Go) HTTP3 support proxy
- [ ] (Rust) low level TCP proxy using Tokio or something similar

## Criatiria to be evaluated

1. Support for multiple protocols specially on top of HTTP
2. Throughput: req/s 
3. Drop request rate and errors
4. Timeouts: req latency time

## Output format & metrics

**TODO**

### Visualization tools

[Datadash](https://github.com/keithknott26/datadash)
