package nprxy_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/artyomturkin/nprxy"
	_ "github.com/artyomturkin/nprxy/protocol/http"
	_ "github.com/artyomturkin/nprxy/transport/plain"
	_ "github.com/artyomturkin/nprxy/transport/tls"
)

func testServer() *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

// Test and benchmark plain listener, http proxy and plain upstream
func TestPlainProxy(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	service := nprxy.ServiceConfig{
		Name:       "test",
		DisableLog: true,
		Listen: nprxy.ListenerConfig{
			Address: "127.0.0.1:59010",
		},
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

func BenchmarkPlainProxy(b *testing.B) {
	ts := testServer()
	defer ts.Close()

	service := nprxy.ServiceConfig{
		Name:       "test",
		DisableLog: true,
		Listen: nprxy.ListenerConfig{
			Address: "127.0.0.1:59010",
		},
		Upstream: ts.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		nprxy.ProxyService(ctx, service)
		wg.Done()
	}()

	b.Run("native", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			resp, _ := http.Get(ts.URL)
			ioutil.ReadAll(resp.Body)
		}
	})

	b.Run("proxy", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			resp, _ := http.Get("http://127.0.0.1:59010/api")
			ioutil.ReadAll(resp.Body)
		}
	})

	cancel()
	wg.Wait()
}

// Test and benchmark tls listener, http proxy and plain upstream
type tbHandler interface {
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
}

func createCertKey(t tbHandler) (string, string) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		t.Fatalf("create ca failed: %v", err)
		return "", ""
	}

	// Public key
	certOut, err := ioutil.TempFile("", "cert")
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()

	// Private key
	keyOut, err := ioutil.TempFile("", "key")
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return certOut.Name(), keyOut.Name()
}

func TestTLSProxy(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	cert, key := createCertKey(t)

	service := nprxy.ServiceConfig{
		Name:       "test",
		DisableLog: true,
		Listen: nprxy.ListenerConfig{
			Address: "127.0.0.1:59010",
			Kind:    "tls",
			TLSCert: cert,
			TLSKey:  key,
		},
		Upstream: ts.URL,
	}
	// Disable TLS verification for tests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	var err error
	go func() {
		err = nprxy.ProxyService(ctx, service)
		wg.Done()
	}()

	resp, _ := http.Get("https://127.0.0.1:59010/api")
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

func BenchmarkTLSProxy(b *testing.B) {
	ts := testServer()
	defer ts.Close()

	cert, key := createCertKey(b)

	service := nprxy.ServiceConfig{
		Name:       "test",
		DisableLog: true,
		Listen: nprxy.ListenerConfig{
			Address: "127.0.0.1:59010",
			Kind:    "tls",
			TLSCert: cert,
			TLSKey:  key,
		},
		Upstream: ts.URL,
	}
	// Disable TLS verification for tests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		nprxy.ProxyService(ctx, service)
		wg.Done()
	}()

	b.Run("native", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			resp, _ := http.Get(ts.URL)
			ioutil.ReadAll(resp.Body)
		}
	})

	b.Run("proxy", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			resp, _ := http.Get("https://127.0.0.1:59010/api")
			ioutil.ReadAll(resp.Body)
		}
	})

	cancel()
	wg.Wait()
}
