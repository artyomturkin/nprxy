package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestOperationResolverSOAP(t *testing.T) {
	type testCase struct {
		name   string
		params map[string]string
		result int
	}

	cases := []testCase{
		testCase{name: "success", params: map[string]string{"SOAPAction": "http://tempuri.org/test"}, result: 200},
		testCase{name: "fail", params: map[string]string{}, result: 400},
	}

	e := echo.New()

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range cs.params {
				req.Header.Set(k, v)
			}

			res := httptest.NewRecorder()

			c := e.NewContext(req, res)

			var op string
			h := OperationResolver("soap")(func(c echo.Context) error {
				op = c.Get("operation").(string)
				return c.JSON(http.StatusOK, "test")
			})

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
			} else if o, ok := cs.params["SOAPAction"]; !ok || o != op {
				if !ok {
					t.Errorf("expected SOAPAction to be set")
				} else {
					t.Errorf("expected SOAPAction to be '%s', got : '%s'", o, op)
				}
			}
		})
	}
}
