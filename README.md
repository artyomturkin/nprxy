# nprxy

## Service Configuration

|Key|Required|Default|Purpose|
|---|--------|-------|-------|
|Name|yes||Name of the service|
|Listen|yes||Endpoint for proxy to listen on. [ip]:port|
|Upstream|yes||Endpoint to forward data to. Schema determines proxy kind (HTTP, TCP)|
|Grace|no|5s|Grace period for proxy to terminate existing connections|

Example:
```json
{
    "services": [
        {
            "name": "testService",
            "listen": ":8080",
            "upstream": "http://localhost",
            "grace": "30s",
            "timeout": "50h"
        }
    ]
}
```
```yaml
services:
- name: testService
  listen: :8080
  upstream: http://localhost
  grace: 30s
  timeout: 50h
```

## HTTP Proxy

### Configurations

|Key|Required|Default|Purpose|
|---|--------|-------|-------|
|Timeout|no|5s|HTTP Request timeout|