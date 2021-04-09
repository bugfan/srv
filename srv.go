package srv

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
)

type Server struct {
	serveMux       *http.ServeMux
	headMiddleware []http.Handler
	tailMiddleware []http.Handler
	addr           string
	tls            *tls.Config
}

func New(addr string) *Server {
	s := &Server{
		serveMux: http.NewServeMux(),
		addr:     addr,
	}
	return s
}

func (s *Server) SetTLSConfig(tls *tls.Config) {
	if tls != nil {
		s.tls = tls
	}
}

func (s *Server) SetTLSConfigFromBytes(certPEM, keyPEM []byte) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Println("SetTLSConfigFromBytes error:", err)
		return
	}

	s.tls = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},
		// MinVersion:               tls.VersionTLS12,
		// MaxVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		Certificates:             []tls.Certificate{cert},
	}
}

func (s *Server) AddHeadHandler(handlers ...http.Handler) {
	s.headMiddleware = append(s.headMiddleware, handlers...)
}
func (s *Server) AddTailHandler(handlers ...http.Handler) {
	s.tailMiddleware = append(s.tailMiddleware, handlers...)
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.serveMux.Handle(pattern, s.wrap(handler))
}

func (s *Server) wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range s.headMiddleware {
			m.ServeHTTP(w, r) // serve
		}
		h.ServeHTTP(w, r)
		for _, m := range s.tailMiddleware {
			m.ServeHTTP(w, r)
		}
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	log.Printf("Listen at %s\n", s.addr)
	if s.tls != nil {
		server := &http.Server{
			Handler:   s,
			TLSConfig: s.tls,
		}
		l, err := net.Listen("tcp", s.addr)
		if err != nil {
			return err
		}
		return server.Serve(l)
	}
	return http.ListenAndServe(s.addr, s)
}
