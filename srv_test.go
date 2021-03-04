package srv

import (
	"fmt"
	"testing"
)

func TestSrv(*testing.T) {
	s, err := New(":8080")
	fmt.Println("err:", err)
	s.AppendServer()
}
