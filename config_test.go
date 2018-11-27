package nprxy_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/artyomturkin/nprxy"
	"github.com/spf13/viper"
)

func TestConfigUnmarshalJSON(t *testing.T) {
	cs := `
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
	}`

	c := &nprxy.Config{}
	viper.SetConfigType("json")
	viper.ReadConfig(bytes.NewBuffer([]byte(cs)))
	err := viper.Unmarshal(c)

	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if len(c.Services) != 1 {
		t.Fatalf("No services definitions: %+v", c)
	}

	if c.Services[0].Grace != 30*time.Second {
		t.Errorf("Wrong service grace: %v, expected: 30s", c.Services[0].Grace)
	}
	if c.Services[0].Timeout != 50*time.Hour {
		t.Errorf("Wrong service timeout: %v, expected: 50h", c.Services[0].Timeout)
	}
	if c.Services[0].Name != "testService" {
		t.Errorf("Wrong service name: %v, expected: testService", c.Services[0].Name)
	}
	if c.Services[0].Listen.Address != ":8080" {
		t.Errorf("listen endpoint is incorrect: %s, expected :8080\n%+v", c.Services[0].Listen, viper.Get("service"))
	}
	if c.Services[0].Listen.TLSCert != "example.crt" {
		t.Errorf("listen cert is incorrect: %s, expected example.crt\n%+v", c.Services[0].Listen, viper.Get("service"))
	}
}

func TestConfigUnmarshalYAML(t *testing.T) {
	cs := `
services:
- name: testService
  listen:
    address: :8080
  upstream: http://localhost
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
        parameters: [client, operation]`

	c := &nprxy.Config{}
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer([]byte(cs)))
	err := viper.Unmarshal(c)

	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if len(c.Services) != 1 {
		t.Fatalf("No services definitions: %+v", c)
	}

	if c.Services[0].Grace != 30*time.Second {
		t.Errorf("Wrong service grace: %v, expected: 30s", c.Services[0].Grace)
	}
	if c.Services[0].Timeout != 50*time.Hour {
		t.Errorf("Wrong service timeout: %v, expected: 50h", c.Services[0].Timeout)
	}
	if c.Services[0].Name != "testService" {
		t.Errorf("Wrong service name: %v, expected: testService", c.Services[0].Name)
	}
	if c.Services[0].Listen.Address != ":8080" {
		t.Errorf("listen endpoint is incorrect: %s, expected :8080\n%+v", c.Services[0].Listen, viper.Get("service"))
	}
	if c.Services[0].HTTP.Authz.Params["parameters"].([]interface{})[0].(string) != "client" {
		t.Errorf("listen endpoint is incorrect: %v, expected client\n%+v", c.Services[0].HTTP.Authz.Params["parameters"].([]interface{})[0], viper.Get("service"))
	}
}
