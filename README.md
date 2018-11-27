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
- name: docker_hub_registry
  listen:
    address: :80
  upstream: https://registry-1.docker.io
  grace: 30s
  timeout: 50h
  http:
    kind: soap
    authn:
      kind: api-key
      params:
        path: example-keys.yaml
    authz:
      kind: casbin
      params:
        model: example_model.conf
        policy: example_policy.csv
        parameters: [client, operation]
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

benchmark                          iter        time/iter   bytes alloc          allocs
---------                          ----        ---------   -----------          ------
BenchmarkPlainProxy/native-4       2000     733.00 μs/op     5277 B/op    70 allocs/op
BenchmarkPlainProxy/proxy-4        1000    1575.29 μs/op    44419 B/op   163 allocs/op
BenchmarkTLSProxy/native-4         2000     896.47 μs/op     5268 B/op    70 allocs/op
BenchmarkTLSProxy/proxy-4          2000    1121.07 μs/op    44185 B/op   162 allocs/op
BenchmarkPlainSOAPProxy/native-4   2000     726.12 μs/op     6318 B/op    88 allocs/op
BenchmarkPlainSOAPProxy/proxy-4     500    2937.28 μs/op    57912 B/op   319 allocs/op
BenchmarkTLSSOAPProxy/native-4     2000     717.99 μs/op     6319 B/op    88 allocs/op
BenchmarkTLSSOAPProxy/proxy-4      1000    2451.32 μs/op    57728 B/op   318 allocs/op
```