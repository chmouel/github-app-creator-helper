package pkg

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fail()
	}
	route, err := c.GetFirstTLSRoute()
	if err != nil {
		t.Fail()
	}
	fmt.Println(route.Spec.Host)
	fmt.Println(route.Spec.TLS)
}
