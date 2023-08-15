package tlv

import (
	"errors"
	"io"
	"reflect"
	"strconv"
)

// MalformedPayloadError indicates given payload is malformed.
type MalformedPayloadError struct {
	msg string
}

func (e *MalformedPayloadError) Error() string {
	return e.msg
}

// Decoder reads and decodes TLV payload from an input stream.
type Decoder struct {
	r   io.RuneReader
	buf []rune

	tagName   string
	tagLength int
	lenLength int
	f         TagLengthTranslator
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.RuneReader, tagName string, bufSize, tagLength, lenLength int, f TagLengthTranslator) *Decoder {
	return &Decoder{
		r:         r,
		buf:       make([]rune, bufSize),
		tagName:   tagName,
		tagLength: tagLength,
		lenLength: lenLength,
		f:         f,
	}
}

// Decode reads the next TLV value from its input and stores it in the value pointed to by dst.
func (d *Decoder) Decode(dst interface{}) error {
	v := reflect.ValueOf(dst)

	if v.Kind() != reflect.Ptr {
		return errors.New("dst should be a pointer, not a value")
	}
	if deref(v.Type()).Kind() != reflect.Struct {
		return errors.New("dst should be a struct")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed")
	}

	indexes := tagIndexMap(v, d.tagName)

	var n int
	var errs []error
	for {
		nn, er := readChunk(d.r, d.buf[n:], d.tagLength, d.lenLength)
		if er != nil {
			if er != io.EOF {
				errs = append(errs, er)
			}
			break
		}

		if er := scan(v, indexes, d.buf[n:n+nn], d.tagLength, d.lenLength, d.f); er != nil {
			errs = append(errs, er)
		}
		n += nn
	}
	if len(errs) != 0 {
		for _, er := range errs {
			if _, ok := er.(*FieldMissingErr); !ok {
				return er
			}
		}
	}

	return nil
}

func readChunk(r io.RuneReader, b []rune, tagLength, lenLength int) (n int, err error) {
	// read Tag
	if len(b) < n+tagLength {
		return n, &MalformedPayloadError{msg: "cannot read tag"}
	}
	nn, err := readRunes(r, b[:tagLength], tagLength)
	if err != nil {
		return
	}
	n += nn

	// read Length
	if len(b) < n+lenLength {
		return n, &MalformedPayloadError{msg: "cannot read value length"}
	}
	nn, err = readRunes(r, b[n:n+lenLength], lenLength)
	if err != nil {
		return
	}
	length, err := strconv.Atoi(string(b[n : n+lenLength]))
	if err != nil {
		return n, &MalformedPayloadError{msg: err.Error()}
	}
	n += nn

	// read Value
	if len(b) < n+length {
		return n, &MalformedPayloadError{msg: "cannot read value"}
	}
	nn, err = readRunes(r, b[n:n+length], length)
	if err != nil {
		return
	}
	n += nn

	return
}

func readRunes(r io.RuneReader, b []rune, n int) (int, error) {
	for i := 0; i < n; i++ {
		chr, _, err := r.ReadRune()
		if err != nil {
			return 0, err
		}
		b[i] = chr
	}
	return n, nil
}
