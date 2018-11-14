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
		        "listen": ":8080",
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
	if c.Services[0].Listen != ":8080" {
		t.Errorf("listen endpoint is incorrect: %s, expected :8080\n%+v", c.Services[0].Listen, viper.Get("service"))
	}
}

func TestConfigUnmarshalYAML(t *testing.T) {
	cs := `
services:
- name: testService
  listen: :8080
  upstream: http://localhost
  grace: 30s
  timeout: 50h`

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
	if c.Services[0].Listen != ":8080" {
		t.Errorf("listen endpoint is incorrect: %s, expected :8080\n%+v", c.Services[0].Listen, viper.Get("service"))
	}
}
