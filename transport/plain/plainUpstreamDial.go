package plain

import (
	"net"

	"github.com/artyomturkin/nprxy"
)

func init() {
	nprxy.UpstreamDialFactory["plain"] = buildPlainUpstreamDialer
}

func buildPlainUpstreamDialer(c nprxy.ServiceConfig) (nprxy.DialUpstream, error) {
	return net.Dial, nil
}
