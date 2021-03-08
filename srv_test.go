package srv

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSrv(*testing.T) {
	s, err := New(":8080")
	fmt.Println("err:", err, s)
	s.AppendHandler(&t1{})
	http.ListenAndServe()
	s.Serve()
}

type t1 struct {
}

func (t *t1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in-->:", r.URL.String())
	fmt.Fprint(w, "t1")
}
