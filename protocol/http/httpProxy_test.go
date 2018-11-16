package http

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	gohttp "net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestHTTPProxy(t *testing.T) {
	handler := func(w gohttp.ResponseWriter, r *gohttp.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	ts := httptest.NewServer(gohttp.HandlerFunc(handler))
	defer ts.Close()

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	pu := "http://" + l.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	u, _ := url.Parse(ts.URL)

	p := &httpProxy{
		Upstream:   u,
		Grace:      time.Second * 30,
		DisableLog: true,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		err = p.Serve(ctx, l, net.Dial)
		wg.Done()
	}()

	resp, _ := gohttp.Get(pu + "/api")
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
	if err != gohttp.ErrServerClosed {
		t.Errorf("Serve failed: %v", err)
	}
}

func BenchmarkHTTPProxy(b *testing.B) {
	handler := func(w gohttp.ResponseWriter, r *gohttp.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	ts := httptest.NewServer(gohttp.HandlerFunc(handler))
	defer ts.Close()

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	pu := "http://" + l.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	u, _ := url.Parse(ts.URL)

	p := &httpProxy{
		Upstream:   u,
		Grace:      time.Second * 30,
		DisableLog: true,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		p.Serve(ctx, l, net.Dial)
		wg.Done()
	}()

	for n := 0; n < b.N; n++ {
		resp, _ := gohttp.Get(pu + "/api")
		ioutil.ReadAll(resp.Body)
	}

	cancel()
	wg.Wait()
}
