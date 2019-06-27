/*
Package tlv implements encoding and decoding of TLV (type-length-value or tag-length-value) as defined in EMV Payment Code.
*/
package tlv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// TLV represents a chunk of TLV payload.
type TLV struct {
	Tag    string
	Length string
	Value  string
}

func (t *TLV) token() string {
	return t.Tag + t.Length + t.Value
}

// FieldMissingErr represents error of field not found for tag.
type FieldMissingErr struct {
	Tag string
}

func (e *FieldMissingErr) Error() string {
	return fmt.Sprintf("missing field for tag %s", string(e.Tag))
}

func scan(v reflect.Value, m map[string]int, token []rune, tagLength, lenLength int, f TagLengthTranslator) error {
	v = reflect.Indirect(v)

	tag := token[:tagLength]
	orgTag := tag
	if f != nil {
		tag, _ = f.Translate(tag, token[tagLength:tagLength+lenLength])
	}

	val := token[tagLength+lenLength:]

	if i, ok := m[string(tag)]; ok {
		f := v.Field(i)

		if isScannable(f.Type()) {
			var res []reflect.Value
			if m, ok := reflect.PtrTo(f.Type()).MethodByName("Scan"); ok {
				res = m.Func.Call([]reflect.Value{f.Addr(), reflect.ValueOf(val)})
			}
			if res == nil {
				return errors.New("unexpected value passed")
			}
			err := res[0].Interface()
			if err == nil {
				return nil
			} else {
				if e, ok := err.(error); ok {
					return e
				}
				return fmt.Errorf("unexpected value returned id: %s", string(tag))
			}
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(string(val))
			return nil
		case reflect.Float64:
			vl, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				return fmt.Errorf(": %s", err)
			}
			f.SetFloat(vl)
			return nil
		case reflect.Slice:
			typ := f.Type().Elem()

			switch typ {
			case reflect.TypeOf(TLV{}):
				rv := reflect.New(typ).Elem()

				for i := 0; i < rv.NumField(); i++ {
					f := rv.Field(i)

					switch typ.Field(i).Name {
					case "Tag":
						f.SetString(string(orgTag))
					case "Length":
						f.SetString(string(token[tagLength : tagLength+lenLength]))
					case "Value":
						f.SetString(string(val))
					}
				}
				f.Set(reflect.Append(f, rv))
				return nil
			}
		}

		return fmt.Errorf("unsupported field type %s passed", f.Kind())
	}

	return &FieldMissingErr{Tag: string(tag)}
}

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func tagIndexMap(v reflect.Value, tagName string) map[string]int {
	v = reflect.Indirect(v)
	t := deref(v.Type())
	m := make(map[string]int, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if emvTag, ok := f.Tag.Lookup(tagName); ok {
			m[emvTag] = i
		}
	}

	return m
}

// TagLengthTranslator is a interface of Tag/Length value translator.
type TagLengthTranslator interface {
	Translate(srcTagName, srcLength []rune) ([]rune, []rune)
}

// TagLengthTranslatorFunc type is an adapter to allow the use of ordinary functions as TagLengthTranslator.
type TagLengthTranslatorFunc func(srcTagName, srcLength []rune) ([]rune, []rune)

// Translate calls f(srcTagName, srcLength).
func (f TagLengthTranslatorFunc) Translate(srcTagName, srcLength []rune) ([]rune, []rune) {
	return f(srcTagName, srcLength)
}

// Scanner is interface for parse various types
type Scanner interface {
	Scan([]rune) (err error)
}

var _scannerInterface = reflect.TypeOf((*Scanner)(nil)).Elem()

func isScannable(t reflect.Type) bool {
	return reflect.PtrTo(t).Implements(_scannerInterface)
}
