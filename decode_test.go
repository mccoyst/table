// © 2014 Steve McCoy.

package table

import (
	"fmt"
	"io"
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func ExampleDecoder_Decode() {
	type X struct {
		A int
		B string
		c int
	}
	lines := `
1,blonde
2,on
3,blonde
`
	dec := NewDecoder(csv.NewReader(strings.NewReader(lines)))
	for {
		var x X
		if err := dec.Decode(&x); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "oops:", err)
			return
		}
		fmt.Println(x.A, x.B, x.c)
	}

	// output: 1 blonde 0
	// 2 on 0
	// 3 blonde 0
}

func TestDecodeNonstruct(t *testing.T) {
	lines := `
1,blonde
2,on
3,blonde
`
	dec := NewDecoder(csv.NewReader(strings.NewReader(lines)))
	var x int
	err := dec.Decode(&x)
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if x != 0 {
		t.Error("Something touched x:", x)
	}
}
