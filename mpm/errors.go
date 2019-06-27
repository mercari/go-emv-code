package mpm

import "fmt"

type errorCode int

const (
	// InvalidFormat represents given payload has invalid format.
	InvalidFormat errorCode = iota + 1
	// InvalidCRC represents CRC is invalid.
	InvalidCRC
)

type genericError struct {
	code errorCode
	msg  string
}

func (e *genericError) Error() string {
	return e.msg
}

// InvalidFormat returns true if code is InvalidFormat.
func (e *genericError) InvalidFormat() bool {
	return e.code == InvalidFormat
}

// NewInvalidFormat creates a new NewInvalidFormat error.
func NewInvalidFormat(msg string) error {
	return &genericError{
		code: InvalidFormat,
		msg:  msg,
	}
}

// InvalidCRC returns true if code is InvalidCRC.
func (e *genericError) InvalidCRC() bool {
	return e.code == InvalidCRC
}

// NewInvalidCRC creates a new NewInvalidCRC error.
func NewInvalidCRC(expected, got uint16) error {
	return &genericError{
		code: InvalidCRC,
		msg:  fmt.Sprintf("mpm: expected CRC is %x not %x", expected, got),
	}
}
