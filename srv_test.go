package srv

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSrv(*testing.T) {
	addr := ":8080"
	// new server
	server := New(addr)

	// set head middleware
	server.AddHeadHandler(&zlog{}, &auth{})

	// set tail middleware
	server.AddTailHandler(&zlog{})

	// set your handler
	server.Handle("/", &yourHandler{})
	server.Handle("/ws/", &yourWebsocketHandler{})
	server.Handle("/static/", &yourStaticHandler{})

	/*
	* if hava tls certificate data
	 */
	// keyData := []byte("xxxx")
	// certData := []byte("xxxx")
	// server.SetTLSConfigFromBytes(certData, keyData)

	/*
	* if hava tls config
	 */
	// tlsconfig := &tls.Config{}
	// server.SetTLSConfig(tlsconfig)

	// listen server
	/*
	* method 1
	* run directly
	 */
	server.Run()

	/*
	* method 2
	* if you have own listener
	 */
	// http.ListenAndServe(addr, server)
}

type yourHandler struct {
}

func (s *yourHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("your handler:", r.URL.String())
	fmt.Fprint(w, "your handler")
}

type yourWebsocketHandler struct {
}

func (s *yourWebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("your websocket handler:", r.URL.String())
	fmt.Fprint(w, "your websocket handler")
}

type yourStaticHandler struct {
}

func (s *yourStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("your static handler:", r.URL.String())
	fmt.Fprint(w, "your static handler")
}

// log middleware
type zlog struct {
}

func (s *zlog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("log handle:", r.URL.String())
}

// auth middleware
type auth struct {
}

func (s *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: auth code
	fmt.Println("auth handle", r.URL.String())
}
