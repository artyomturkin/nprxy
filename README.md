# nprxy

[![Go Report Card](https://goreportcard.com/badge/github.com/artyomturkin/nprxy)](https://goreportcard.com/report/github.com/artyomturkin/nprxy)

## Service Configuration

|Key|Required|Default|Purpose|
|---|--------|-------|-------|
|Name|yes||Name of the service|
|Listen.Address|yes||Endpoint for proxy to listen on. [ip]:port|
|Listen.Kind|no|plain|Listen endpoint type: plain, tls|
|Listen.tlsCert|no||Path to TLS cert. Required if Kind=tls|
|Listen.tlsKey|no||Path to TLS key. Required if Kind=tls|
|Upstream|yes||Endpoint to forward data to. Schema determines proxy kind (HTTP, TCP)|
|Grace|no|5s|Grace period for proxy to terminate existing connections|

### Examples:

Configuration in json format with TLS listener
```json
{
    "services": [
        {
            "name": "testService",
            "listen": {
                "kind": "tls",
                "address": ":8080",
                "tlsCert": "example.crt",
                "tlsKey": "example.key"
            },
            "upstream": "http://localhost",
            "grace": "30s",
            "timeout": "50h"
        }
    ]
}
```

Configuration in yaml format with plain http listener
```yaml
services:
- name: testService
  listen:
    address: :8080
  upstream: http://localhost
  grace: 30s
  timeout: 50h
```


## HTTP Proxy

### Configurations

|Key|Required|Default|Purpose|
|---|--------|-------|-------|
|Timeout|no|5s|HTTP Request timeout|

## Benchmarks


```
goos: windows
goarch: amd64
pkg: github.com/artyomturkin/nprxy
PASS

benchmark                      iter        time/iter   bytes alloc          allocs
---------                      ----        ---------   -----------          ------
BenchmarkPlainProxy/native-4   2000     640.50 μs/op     4981 B/op    65 allocs/op
BenchmarkPlainProxy/proxy-4    1000    1331.99 μs/op    43872 B/op   153 allocs/op

BenchmarkTLSProxy/native-4     2000     635.00 μs/op     4982 B/op    65 allocs/op
BenchmarkTLSProxy/proxy-4      2000     859.00 μs/op    43878 B/op   155 allocs/op
BenchmarkTLSProxy/proxyLog-4   1000    1445.11 μs/op    44320 B/op   164 allocs/op
```