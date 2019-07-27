package tlv

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unicode/utf8"
)

const tlvEntityFormat = "%s%s%s"

// Encoder writes EMV Payment Code payload to an output stream.
type Encoder struct {
	w          io.Writer
	tagName    string
	ignoreTags map[string]struct{}
	f          TagLengthTranslator
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, tagName string, ignoreTags []string, f TagLengthTranslator) *Encoder {
	m := make(map[string]struct{}, len(ignoreTags))
	for _, v := range ignoreTags {
		m[v] = struct{}{}
	}
	return &Encoder{
		w:          w,
		tagName:    tagName,
		ignoreTags: m,
		f:          f,
	}
}

// Encode writes TLV payload of src to the stream.
func (e *Encoder) Encode(src interface{}) error {
	v := reflect.ValueOf(src)

	if v.IsNil() {
		return errors.New("nil pointer passed")
	}

	indexes := tagIndexMap(v, e.tagName)

	v = reflect.Indirect(v)
	for id, index := range indexes {
		if _, ok := e.ignoreTags[id]; ok {
			continue
		}

		f := v.Field(index)
		if isTokenizable(f.Type()) {
			var res []reflect.Value
			if m, ok := reflect.PtrTo(f.Type()).MethodByName("Tokenize"); ok {
				res = m.Func.Call([]reflect.Value{f.Addr()})
			}
			if res == nil {
				return errors.New("unexpected value passed")
			}

			err := res[1].Interface()
			if err == nil {
				switch nv := res[0].Interface().(type) {
				case string:
					f = reflect.ValueOf(nv)
				default:
					return fmt.Errorf("unexpected Tokenizer return type id: %v type: %s", id, nv)
				}
			} else {
				if e, ok := err.(error); ok {
					return e
				}
				return errors.New("unexpected value returned")
			}
		}

		v, err := fieldToString(f)
		if err != nil {
			return fmt.Errorf("failed to convert field value to string: %s", err)
		}
		if len(v) < 1 {
			continue // value should be non-zero length
		}

		length := fmt.Sprintf("%02d", utf8.RuneCountInString(v))

		if e.f != nil {
			strID, strLength := e.f.Translate([]rune(id), []rune(length))
			id = string(strID)
			length = string(strLength)
		}

		if _, err := e.w.Write([]byte(fmt.Sprintf(tlvEntityFormat, id, length, v))); err != nil {
			return fmt.Errorf("failed to write body: %s", err)
		}
	}

	return nil
}

func fieldToString(v reflect.Value) (ret string, err error) {
	switch v.Kind() {
	case reflect.String:
		ret = v.String()
	case reflect.Float64:
		ret = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Slice:
		typ := v.Type().Elem()

		switch typ {
		case reflect.TypeOf(TLV{}):
			for i := 0; i < v.Cap(); i++ {
				y := v.Index(i).Interface().(TLV)
				ret = ret + y.token()
			}
		default:
			return "", fmt.Errorf("unsupported slice element type %s passed", typ.Kind())
		}
	default:
		return "", fmt.Errorf("unsupported field type %s passed", v.Kind())
	}
	return
}

// Tokenizer is the interface providing the Tokenize method.
type Tokenizer interface {
	Tokenize() (string, error)
}

var _tokenizerInterface = reflect.TypeOf((*Tokenizer)(nil)).Elem()

func isTokenizable(t reflect.Type) bool {
	return reflect.PtrTo(t).Implements(_tokenizerInterface)
}
