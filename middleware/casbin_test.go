package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin"
	"github.com/labstack/echo"
)

func TestCasbinAllowed(t *testing.T) {
	ce := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")

	e := echo.New()

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set("system", "alice")
	c.Set("operation", "data1")

	h := CasbinEnforcer(ce, ValueFromContext("system"), ValueFromContext("operation"))(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := h(c)

	if err != nil {
		if errObj, ok := err.(*echo.HTTPError); ok {
			if errObj.Code != 200 {
				t.Errorf("expected 200, got: %d", errObj.Code)
			}
		} else {
			t.Error(err)
		}
	} else if c.Response().Status != 200 {
		t.Errorf("expected 200, got: %d", c.Response().Status)
	}
}

func TestCasbinForbidden(t *testing.T) {
	ce := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")

	e := echo.New()

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set("system", "bob")
	c.Set("operation", "data1")

	h := CasbinEnforcer(ce, ValueFromContext("system"), ValueFromContext("operation"))(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := h(c)

	if err != nil {
		if errObj, ok := err.(*echo.HTTPError); ok {
			if errObj.Code != 403 {
				t.Errorf("expected 403, got: %d", errObj.Code)
			}
		} else {
			t.Error(err)
		}
	} else if c.Response().Status != 403 {
		t.Errorf("expected 403, got: %d", c.Response().Status)
	}
}

func BenchmarkCasbin(t *testing.B) {
	ce := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")
	e := echo.New()

	t.Run("allowed", func(b *testing.B) {
		for index := 0; index < b.N; index++ {
			req := httptest.NewRequest("GET", "/", nil)
			res := httptest.NewRecorder()

			c := e.NewContext(req, res)
			c.Set("system", "alice")
			c.Set("operation", "data1")

			h := CasbinEnforcer(ce, ValueFromContext("system"), ValueFromContext("operation"))(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			err := h(c)

			if err != nil {
				if errObj, ok := err.(*echo.HTTPError); ok {
					if errObj.Code != 200 {
						b.Errorf("expected 200, got: %d", errObj.Code)
					}
				} else {
					b.Error(err)
				}
			} else if c.Response().Status != 200 {
				b.Errorf("expected 200, got: %d", c.Response().Status)
			}
		}
	})

	t.Run("forbidden", func(b *testing.B) {
		for index := 0; index < b.N; index++ {
			req := httptest.NewRequest("GET", "/", nil)
			res := httptest.NewRecorder()

			c := e.NewContext(req, res)
			c.Set("system", "bob")
			c.Set("operation", "data1")

			h := CasbinEnforcer(ce, ValueFromContext("system"), ValueFromContext("operation"))(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			err := h(c)

			if err != nil {
				if errObj, ok := err.(*echo.HTTPError); ok {
					if errObj.Code != 403 {
						b.Errorf("expected 403, got: %d", errObj.Code)
					}
				} else {
					b.Error(err)
				}
			} else if c.Response().Status != 403 {
				b.Errorf("expected 403, got: %d", c.Response().Status)
			}
		}
	})

}
