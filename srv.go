package srv

import (
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
)

type Server struct {
	handlers []http.Handler
}

func New(listen string) (interface{}, error) {
	// 单端口复用

	l, err := net.Listen("tcp", listen)
	if err != nil {
		return nil, err
	}
	m := cmux.New(l)
	// tlsL := m.Match(cmux.TLS())
	// http2 := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.HTTP1Fast())
	go p.GetHTTPServer(listen).Serve(httpL)
	// go p.GetHTTPSServer(listen).ServeTLS(http2, "", "")
	// go p.GetHTTPSServer(listen).ServeTLS(tlsL, "", "")
	return m, nil
}

func (s *Server) AppendServer(handlers ...http.Handler) {
	for _, h := range handlers {
		s.handlers = append(s.handlers, h)
	}
}
