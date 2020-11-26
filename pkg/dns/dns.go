package dns

import (
	"github.com/google/gopacket/layers"
	"net"
)

type ZoneType uint16

const (
	DNSForwardLookupZone ZoneType = 1
	DNSReverseLookupZone ZoneType = 2
)

type Handler interface {
	serveDNS(*udpConnection, *layers.DNS)
}

type Server struct {
	port    int
	handler Handler
}

func NewServer(port int) *Server {

}

type serveMux struct {
	handler map[string]Handler
}

func NewServerMux() *serveMux {
	h := make(map[string]Handler)
	return &serveMux{handler: h}
}

type udpConnection struct {
	conn net.PacketConn
	addr net.Addr
}
