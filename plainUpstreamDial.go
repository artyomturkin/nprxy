package nprxy

import "net"

func init() {
	upstreamDialFactory["plain"] = buildPlainUpstreamDialer
}

func buildPlainUpstreamDialer(c ServiceConfig) (DialUpstream, error) {
	return net.Dial, nil
}
