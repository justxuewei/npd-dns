package ndp_dns

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

type customHandler func(string) (string, error)

func generateHandler(records map[string]string, lookupFunc customHandler) func (w *udpConnection, r *layers.DNS) {
	return func(w *udpConnection, r *layers.DNS) {
		switch r.Questions[0].Type {
		case layers.DNSTypeA:
			handleATypeQuery(w, r, records, lookupFunc)
		}
	}
}

func (srv *Server) AddZoneData(zone string, records map[string]string,
	lookupFunc func(string) (string, error), lookupZone ZoneType) {
	if lookupZone == DNSForwardLookupZone {
		serverMuxCurrent := srv.handler.(*ServeMux)
		serverMuxCurrent.handleFunc(zone, generateHandler(records, lookupFunc))
	}
}

func NewServer(port int) *Server {
	handler := NewServerMux()
	return &Server{port: port, handler: handler}
}

type ServeMux struct {
	handler map[string]Handler
}

func NewServerMux() *ServeMux {
	h := make(map[string]Handler)
	return &ServeMux{handler: h}
}

func (srv *ServeMux) handleFunc(pattern string, f func(*udpConnection, *layers.DNS)) {
	srv.handler[pattern] = handlerConvert(f)
}

func (srv *ServeMux) serveDNS(u *udpConnection, request *layers.DNS) {
	//var h Handler
	//if len(request.Questions) < 1 { // allow more than one question
	//	return
	//}
	//if h = srv.ma
}

func (srv *ServeMux) match(q string, t layers.DNSType) Handler {
	var handler Handler
	b := make([]byte, len(q))
	off := 0
	end := false
	for {
		l := len(q[off:])
		for i := 0; i < l; i++ {
			b[i] = q[off+1]
			if b[i] >= 'A' && b[i] <= 'Z' {
				// TODO: What is |=?
				b[i] |= 'a' - 'A'
			}
		}
		if h, ok := srv.handler[string(b[:l])]; ok {
			if uint16(t) != uint16(43) {
				return h
			}
			handler = h
		}
		off, end = nextLabel(q, off)
		if end { break }
	}
	if h, ok := srv.handler["."]; ok {
		return h
	}
	return handler
}

func nextLabel(s string, offset int) (i int, end bool) {
	quote := false
	for i = offset; i < len(s)-1; i++ {
		switch s[i] {
		case '\\':
			quote = !quote
		case '.':
			if quote {
				quote = !quote
				continue
			}
			return i+1, false
		default:
			quote = false
		}
	}
	return i+1, true
}

type udpConnection struct {
	conn net.PacketConn
	addr net.Addr
}

type handlerConvert func(*udpConnection, *layers.DNS)

// TODO: What is the purpose of this method for a func type?
func (f handlerConvert) serveDNS(w *udpConnection, r *layers.DNS) {
	f(w, r)
}

func handleATypeQuery(w *udpConnection, r *layers.DNS, records map[string]string, lookupFunc customHandler) {
	panic("handleATypeQuery() has not been implemented.")
}
