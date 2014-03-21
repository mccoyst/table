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

func TestDecodeError(t *testing.T) {
	type X struct {
		A int
		B complex64
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
		t.Error("Expected an error")
	}
	if de, ok := err.(DecodeError); !ok {
		t.Error("Expected a DecodeError, got", err)
	} else if de != "complex64" {
		t.Error("Expected the error to be on complex64, got", de)
	}
}

func TestVariousParses(t *testing.T) {
	type X struct {
		A int
		B string
		C uint
		D int8
		E int16
		F int32
		G int64
		H uint8
		I uint16
		J uint32
		K uint64
		L float32
		M float64
		N bool
	}
	lines := `
-1,meow,2,3,4,5,6,7,8,9,10,11.1,12.6,true
`
	dec := NewDecoder(csv.NewReader(strings.NewReader(lines)))
	for {
		var x X
		if err := dec.Decode(&x); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal("Expected no error, got", err)
		}
		if x.A != -1 {
			t.Error("Expected A to be -1, got", x.A)
		}
		if x.B != "meow" {
			t.Error("Expected B to be meow, got", x.B)
		}
		if x.C != 2 {
			t.Error("Expected C to be 2, got", x.C)
		}
		if x.D != 3 {
			t.Error("Expected D to be 3, got", x.D)
		}
		if x.E != 4 {
			t.Error("Expected E to be 4, got", x.E)
		}
		if x.F != 5 {
			t.Error("Expected F to be 5, got", x.F)
		}
		if x.G != 6 {
			t.Error("Expected G to be 6, got", x.G)
		}
		if x.H != 7 {
			t.Error("Expected H to be 7, got", x.H)
		}
		if x.I != 8 {
			t.Error("Expected I to be 8, got", x.I)
		}
		if x.J != 9 {
			t.Error("Expected J to be 9, got", x.J)
		}
		if x.K != 10 {
			t.Error("Expected K to be 10, got", x.K)
		}
		if x.L != 11.1 {
			t.Error("Expected L to be 11, got", x.L)
		}
		if x.M != 12.6 {
			t.Error("Expected M to be 12, got", x.M)
		}
		if !x.N {
			t.Error("Expected N to be true, got", x.N)
		}
	}
}
