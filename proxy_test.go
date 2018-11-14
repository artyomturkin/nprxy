package nprxy_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/artyomturkin/nprxy"
)

func testServer() *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestProxy(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	service := nprxy.ServiceConfig{
		Name:     "test",
		Listen:   "127.0.0.1:59010",
		Upstream: ts.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	var err error
	go func() {
		err = nprxy.ProxyService(ctx, service)
		wg.Done()
	}()

	resp, _ := http.Get("http://127.0.0.1:59010/api")
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

func BenchmarkProxy(b *testing.B) {
	ts := testServer()
	defer ts.Close()

	service := nprxy.ServiceConfig{
		Name:     "test",
		Listen:   "127.0.0.1:59010",
		Upstream: ts.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		nprxy.ProxyService(ctx, service)
		wg.Done()
	}()

	for n := 0; n < b.N; n++ {
		resp, _ := http.Get("http://127.0.0.1:59010/api")
		ioutil.ReadAll(resp.Body)
	}

	cancel()
	wg.Wait()
}
