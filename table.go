// © 2014 Steve McCoy.

/*
Package table is used to decode CSV streams into arbitrary structs.

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
			fmt.Fprintf(os.Stderr, "oops:", err)
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

type DecodeError string

func (d DecodeError) Error() string {
	return string(d) + " is not decodable"
}

type FieldReader interface {
	Read() ([]string, error)
}

type Decoder struct {
	Modify map[reflect.Kind]func(*reflect.Value, string)error
	r FieldReader
}

func NewDecoder(r FieldReader) Decoder {
	return Decoder{defaultMods, r}
}

func (d *Decoder) Decode(s interface{}) error {
	fields, err := d.r.Read()
	if err != nil {
		return err
	}

	t := reflect.TypeOf(s)
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	val := reflect.ValueOf(s)

	j := 0 // j is the index of val's field i in the fields slice
	for i := 0; i < t.NumField(); i++ {
		if j >= len(fields) {
			return RowError{ len(fields), j, t.Field(i).Name }
		}

		f := t.Field(i)
		fv := val.Field(i)
		if !fv.CanSet() {
			continue
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
