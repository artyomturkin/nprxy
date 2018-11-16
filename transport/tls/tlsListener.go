package tls

import (
	"crypto/tls"
	"net"

	"github.com/artyomturkin/nprxy"
)

func init() {
	nprxy.ListenerFactory["tls"] = buildTLSListener
}

func buildTLSListener(c nprxy.ServiceConfig) (net.Listener, error) {
	cer, err := tls.LoadX509KeyPair(c.Listen.TLSCert, c.Listen.TLSKey)
	if err != nil {
		return nil, err
	}
	tc := &tls.Config{Certificates: []tls.Certificate{cer}}
	return tls.Listen("tcp", c.Listen.Address, tc)
}
