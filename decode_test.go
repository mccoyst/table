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

func TestShortRow(t *testing.T) {
	type X struct {
		A int
		B string
		C int
	}
	lines := `
1,blonde
2,on
3,blonde
`

	dec := NewDecoder(csv.NewReader(strings.NewReader(lines)))
	var x X
	err := dec.Decode(&x)
	if err == nil {
		t.Error("Expected an error", x)
	}
	if re, ok := err.(RowError); ok {
		if re.RowLen != 2 {
			t.Error("Expected RowLen of 2, got", re.RowLen)
		}
		if re.StructLen != 3 {
			t.Error("Expected StructLen of 3, got", re.StructLen)
		}
		if re.MissingField != "C" {
			t.Error("Expected MissingField of C, got", re.MissingField)
		}
	} else {
		t.Error("Expected a RowError, got", err)
	}
}

func TestLongRow(t *testing.T) {
	type X struct {
		A int
		B string
		c int
	}
	lines := `
1,blonde,6
2,on,6
3,blonde,6
`

	dec := NewDecoder(csv.NewReader(strings.NewReader(lines)))
	var x X
	err := dec.Decode(&x)
	if err == nil {
		t.Error("Expected an error", x)
	}
	if re, ok := err.(RowError); ok {
		if re.RowLen != 3 {
			t.Error("Expected RowLen of 3, got", re.RowLen)
		}
		if re.StructLen != 2 {
			t.Error("Expected StructLen of 2, got", re.StructLen)
		}
		if re.MissingField != "" {
			t.Error("Expected empty MissingField, got", re.MissingField)
		}
	} else {
		t.Error("Expected a RowError, got", err)
	}
}

