package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestBcryptAPIKey(t *testing.T) {
	type testCase struct {
		name   string
		system string
		key    string
		result int
	}

	cases := []testCase{
		testCase{name: "success", system: "test-system", key: "api-key", result: 200},
		testCase{name: "unauth", system: "test-system", key: "api-key-2", result: 401},
	}

	e := echo.New()

	h := BCryptAPIKey(map[string]string{"test-system": "$2a$10$0ZYFiKcYonvy.y/P4jAzJOr79AQoeO1LGO2hyj27QS5pTx/1nyzRm"})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Add("X-NPRXY-Client", cs.system)
			req.Header.Add("X-NPRXY-Key", cs.key)

			res := httptest.NewRecorder()

			c := e.NewContext(req, res)
			err := h(c)

			if err != nil {
				if errObj, ok := err.(*echo.HTTPError); ok {
					if errObj.Code != cs.result {
						t.Errorf("expected %d, got: %d", cs.result, errObj.Code)
					}
				} else {
					t.Error(err)
				}
			} else if c.Response().Status != cs.result {
				t.Errorf("expected %d, got: %d", cs.result, c.Response().Status)
			}
		})
	}

}
