package nprxy

import (
	"crypto/tls"
	"net"
)

func init() {
	listenerFactory["tls"] = buildTLSListener
}

func buildTLSListener(c ServiceConfig) (net.Listener, error) {
	cer, err := tls.LoadX509KeyPair(c.Listen.TLSCert, c.Listen.TLSKey)
	if err != nil {
		return nil, err
	}
	tc := &tls.Config{Certificates: []tls.Certificate{cer}}
	return tls.Listen("tcp", c.Listen.Address, tc)
}
