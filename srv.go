package srv

import (
	"crypto/tls"
	"net"
	"net/http"
	"reflect"
)

type GMux interface {
	AppendHandler(handlers ...interface{})
	Serve() error
	Close()
}

type Matcher interface {
	Match(*http.Request) bool
}
type Handler interface {
	ServeHTP(http.ResponseWriter, *http.Request)
}

type Server struct {
	instances []interface{}
	proxys    []Proxy
	addr      string
	listener  net.Listener
}
type Proxy interface {
	Matcher
}

func (s *Server) AppendHandler(ifs ...interface{}) {
	for _, i := range ifs {
		reflectVal := reflect.ValueOf(i)
		t := reflect.Indirect(reflectVal).Type()
		newObj := reflect.New(t)
		handler, ok := newObj.Interface().(Handler)
		if ok {
			matcher, ok := newObj.Interface().(Matcher)
			if ok {

			}
			s.proxys = append(s.proxys, Proxy)
		}
	}

}
func (s *Server) Close() {
	s.Close()
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range s.handlers {
		h.ServeHTTP(w, r)
	}
}
func (s *Server) Serve() error {
	server := &http.Server{
		Handler: s,
		Addr:    s.addr,
	}
	return server.Serve(s.listener)
}

func New(addr string) (GMux, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		addr:     addr,
		listener: l,
	}
	return s, nil
}

func getTLSConfig() *tls.Config {
	tlsconf := &tls.Config{
		// GetCertificate: p.GetCertificate,
		// NextProtos: []string{ACMETLS1Protocol},
	}
	tlsconf.InsecureSkipVerify = true

	tlsconf.PreferServerCipherSuites = true

	tlsconf.CipherSuites = []uint16{
		//tls.TLS_AES_128_GCM_SHA256,
		//tls.TLS_CHACHA20_POLY1305_SHA256,
		//tls.TLS_AES_256_GCM_SHA384,
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
	}
	return tlsconf
}
