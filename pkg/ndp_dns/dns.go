package ndp_dns

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type ZoneType uint16

const (
	DNSForwardLookupZone ZoneType = 1
	//DNSReverseLookupZone ZoneType = 2
)

type Handler interface {
	serveDNS(*udpConnection, *layers.DNS)
}

type Server struct {
	ip      string
	port    int
	handler Handler
}

type customHandler func(string) (string, error)

func generateHandler(records map[string]string, lookupFunc customHandler) func(w *udpConnection, r *layers.DNS) {
	return func(w *udpConnection, r *layers.DNS) {
		// In this example, the handler will only handle the DNS request with type A.
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

func (srv *Server) LoadConf() {
	conf := Conf{}
	records, err := conf.Read()
	if err != nil {
		fmt.Printf("Can't read conf from \"%s\", %s.\n", conf.path, err)
		return
	}
	for _, record := range records {
		srv.AddZoneData(record.DomainName, record.Map, nil, DNSForwardLookupZone)
	}
}

func (srv *Server) StartAndServe() {
	if srv.ip == "" {
		srv.ip = "0.0.0.0"
	}
	addr := net.UDPAddr{
		Port: srv.port,
		IP: net.ParseIP(srv.ip),
	}
	l, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Start to serve and listen at %s:%d.", addr.IP, addr.Port)
	udpConnection := &udpConnection{conn: l}
	srv.serve(udpConnection)
}

func (srv *Server) serve(u *udpConnection) {
	for {
		tmp := make([]byte, 1024)
		_, addr, err := u.conn.ReadFrom(tmp)
		if err != nil {
			panic(err)
		}
		u.addr = addr
		// NewPacket creates a new Packet object from a set of bytes.
		// layers.LayerTypeDNS: DNS Decoder
		// gopacket.Default: Decode Options
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		// Layer returns the first layer in this packet of the given type, or nil
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		srv.handler.serveDNS(u, tcp)
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
	// Is handlerConvert() an explicit type conversion?
	srv.handler[pattern] = handlerConvert(f)
}

func (srv *ServeMux) serveDNS(u *udpConnection, request *layers.DNS) {
	var h Handler
	if len(request.Questions) < 1 { // allow more than one question
		return
	}
	if h = srv.match(string(request.Questions[0].Name), request.Questions[0].Type); h == nil {
		fmt.Println("No handler found for", string(request.Questions[0].Name))
		return
	}
	h.serveDNS(u, request)
}

func (srv *ServeMux) match(q string, t layers.DNSType) Handler {
	var handler Handler
	b := make([]byte, len(q))
	off := 0
	end := false
	for {
		l := len(q[off:])
		for i := 0; i < l; i++ {
			b[i] = q[off+i]
			if b[i] >= 'A' && b[i] <= 'Z' {
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
		if end {
			break
		}
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
			return i + 1, false
		default:
			quote = false
		}
	}
	return i + 1, true
}

type udpConnection struct {
	conn net.PacketConn
	addr net.Addr
}

func (udp *udpConnection) Write(b []byte) error {
	_, _ = udp.conn.WriteTo(b, udp.addr)
	return nil
}

type handlerConvert func(*udpConnection, *layers.DNS)

func (f handlerConvert) serveDNS(w *udpConnection, r *layers.DNS) {
	f(w, r)
	//panic("serveDNS() has not been implemented.")
}

func handleATypeQuery(w *udpConnection, r *layers.DNS, records map[string]string, lookupFunc customHandler) {
	replyMess := r
	var ip string
	var err error
	var ok bool
	if lookupFunc == nil {
		// what is the questions
		ip, ok = records[string(r.Questions[0].Name)]
		if !ok {
			fmt.Println("No IP found in records.")
		}
	} else {
		ip, err = lookupFunc(string(r.Questions[0].Name))
	}
	a, _, _ := net.ParseCIDR(ip + "/24")
	// DNSResourceRecord, see https://en.wikipedia.org/wiki/Domain_Name_System#Resource_records
	var answer layers.DNSResourceRecord
	// TYPE: specifies the type of query.
	answer.Type = layers.DNSTypeA
	// RDATA
	answer.IP = a
	// Name: Name of the node to which this record pertains
	answer.Name = r.Questions[0].Name
	answer.Class = layers.DNSClassIN

	// DNS Message, see https://en.wikipedia.org/wiki/Domain_Name_System#DNS_message_format
	// QR: A one-bit field that specifies whether this message is a query (0), or a response (1).
	replyMess.QR = true
	// ANCOUNT: 16-bit field, denoting the number of answers in response
	replyMess.ANCount = 1
	// OPCODE: The type can be QUERY (standard query, 0), IQUERY (inverse query, 1), or STATUS (server status request, 2)
	replyMess.OpCode = layers.DNSOpCodeNotify
	// AA: Authoritative Answer, in a response, indicates if the DNS server is authoritative for the queried hostname
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, answer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	_ = w.Write(buf.Bytes())
}
