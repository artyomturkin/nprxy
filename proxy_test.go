package nprxy_test

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo"

	"github.com/labstack/echo/middleware"

	"github.com/artyomturkin/nprxy"
)

func TestHTTPProxy(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	pUrl := "http://" + l.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	u, _ := url.Parse(ts.URL)

	p := &nprxy.HTTPProxy{
		Upstream:    u,
		Grace:       time.Second * 30,
		Middlewares: []echo.MiddlewareFunc{middleware.Logger()},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		err = p.Serve(ctx, l, net.Dial)
		wg.Done()
	}()

	resp, _ := http.Get(pUrl + "/api")
	body, _ := ioutil.ReadAll(resp.Body)

	cancel()

	if resp.StatusCode != 200 {
		t.Errorf("Wrong status code: %d, expected 200", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("wrong Content-Type: %s, expected: text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	}
	if string(body) != "<html><body>Hello World!</body></html>" {
		t.Errorf("Wrong body: %s, expected: <html><body>Hello World!</body></html>", string(body))
	}

	wg.Wait()
	if err != http.ErrServerClosed {
		t.Errorf("Serve failed: %v", err)
	}
}

func BenchmarkHTTPProxy(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	pUrl := "http://" + l.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	u, _ := url.Parse(ts.URL)

	p := &nprxy.HTTPProxy{
		Upstream: u,
		Grace:    time.Second * 30,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		p.Serve(ctx, l, net.Dial)
		wg.Done()
	}()

	for n := 0; n < b.N; n++ {
		resp, _ := http.Get(pUrl + "/api")
		ioutil.ReadAll(resp.Body)
	}

	cancel()
	wg.Wait()
}
