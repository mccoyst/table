// © 2014 Steve McCoy.

/*
Package table is used to decode CSV-like streams into arbitrary structs.

For example:

	type X struct {
		A int
		B string
		c int
	}

	dec := table.NewDecoder(csvReader)
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
*/
package table

import (
	"reflect"
	"strconv"
)

// RowError is returned from Decode when the number of fields in a row
// does not equal the number of exported fields in the destination struct.
// If there are more row fields than struct fields, MissingField will contain
// the name of the next available field.
type RowError struct {
	RowLen int
	StructLen int
	MissingField string
}

func (r RowError) Error() string {
	msg := "row mismatch: row length = " + strconv.Itoa(r.RowLen) +
		", but struct length = " + strconv.Itoa(r.StructLen)
	if r.MissingField != "" {
		msg += " (field " + r.MissingField + ")"
	}
	return msg
}

// DecodeError is returned from Decode if a field is of a Kind that
// does not have an associated function in Modify.
type DecodeError string

func (d DecodeError) Error() string {
	return string(d) + " is not decodable"
}

// FieldReader represents anything that behaves similar to
// encoding/csv's Reader type. Any errors encoundered
// by the reader will be immediately returned by Decode.
type FieldReader interface {
	Read() ([]string, error)
}

// Decoder contains a map of functions from reflect.Kinds to 
// functions that should set a *reflect.Value of the associated Kind
// with the value represented by a provided string.
type Decoder struct {
	Modify map[reflect.Kind]func(*reflect.Value, string)error
	r FieldReader
}

// NewDecoder returns a Decoder that reads from r and has a default
// Modify map that can set values for bool, int types, float types, and strings.
func NewDecoder(r FieldReader) Decoder {
	return Decoder{defaultMods, r}
}

// Decode sets the exported fields of the struct s with the values
// represented by the fields in the next row provided by d's FieldReader.
// Fields are parsed and set using the functions in d.Modify.
//
// Any errors from Read are returned immediately.
// If s is not a pointer to a struct, Decode returns nil and *s is not modified.
// A DecodeError is returned for the first field whose Kind has
// no entry in d.Modify. A RowError is returned when the row has too many
// or too few fields for s.
func (d *Decoder) Decode(s interface{}) error {
	fields, err := d.r.Read()
	if err != nil {
		return err
	}

	t := reflect.TypeOf(s)
	if t == nil || (t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Struct) {
		return nil
	}
	t = t.Elem()

	val := reflect.ValueOf(s).Elem()

	j := 0 // j is the index of val's field i in the fields slice
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)
		if !fv.CanSet() {
			continue
		}

		if j >= len(fields) {
			return RowError{ len(fields), j, t.Field(i).Name }
		}

		m, ok := d.Modify[f.Type.Kind()]
		if !ok {
			return DecodeError(f.Type.Kind().String())
		}
		m(&fv, fields[j])
		j++
	}

	if j < len(fields) {
		return RowError{ len(fields), j, "" }
	}

	return nil
}

func modInt(v *reflect.Value, f string, bitSize int) error {
	n, err := strconv.ParseInt(f, 10, bitSize)
	v.SetInt(n)
	return err
}

func modUint(v *reflect.Value, f string, bitSize int) error {
	n, err := strconv.ParseUint(f, 10, bitSize)
	v.SetUint(n)
	return err
}

var defaultMods = map[reflect.Kind]func(*reflect.Value, string)error {
	reflect.Bool: func(v *reflect.Value, f string) error {
		b, err := strconv.ParseBool(f)
		v.SetBool(b)
		return err
	},
	reflect.Int: func(v *reflect.Value, f string) error {
		return modInt(v, f, strconv.IntSize)
	},
	reflect.Int8: func(v *reflect.Value, f string) error {
		return modInt(v, f, 8)
	},
	reflect.Int16: func(v *reflect.Value, f string) error {
		return modInt(v, f, 16)
	},
	reflect.Int32: func(v *reflect.Value, f string) error {
		return modInt(v, f, 32)
	},
	reflect.Int64: func(v *reflect.Value, f string) error {
		return modInt(v, f, 64)
	},
	reflect.Uint: func(v *reflect.Value, f string) error {
		return modUint(v, f, strconv.IntSize)
	},
	reflect.Uint8: func(v *reflect.Value, f string) error {
		return modUint(v, f, 8)
	},
	reflect.Uint16: func(v *reflect.Value, f string) error {
		return modUint(v, f, 16)
	},
	reflect.Uint32: func(v *reflect.Value, f string) error {
		return modUint(v, f, 32)
	},
	reflect.Uint64: func(v *reflect.Value, f string) error {
		return modUint(v, f, 64)
	},
	reflect.Float32: func(v *reflect.Value, f string) error {
		n, err := strconv.ParseFloat(f, 32)
		v.SetFloat(n)
		return err
	},
	reflect.Float64: func(v *reflect.Value, f string) error {
		n, err := strconv.ParseFloat(f, 64)
		v.SetFloat(n)
		return err
	},
	reflect.String: func(v *reflect.Value, f string) error {
		v.SetString(f)
		return nil
	},
}
