package auth

import (
	"testing"
)

func TestInScope(t *testing.T) {
	if IsInPathScope("/wc24/hi", "/wc24") != true {
		t.Fatalf(`Wanted true for '/wc24/hi inside /wc24'?`)
	}
	if IsInPathScope("/hello/there", "/") != true {
		t.Fatalf(`Wanted true for '/hello/there inside /'?`)
	} 
	if IsInPathScope("/there/hi/there", "/hi/there") != false {
		t.Fatalf(`Wanted false for '/there/hi/there inside /hi/there'?`)
	} 
}