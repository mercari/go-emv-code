package mpm_test

import (
	"fmt"
	"testing"

	"go.mercari.io/go-emv-code/mpm"
)

func TestNewInvalidFormat(t *testing.T) {
	msg := "testing"
	err := mpm.NewInvalidFormat(msg)
	if err == nil {
		t.Fatal("should not be nil")
	}
	if err.Error() != msg {
		t.Errorf("unexpexted value expected: %s, give: %s", msg, err.Error())
	}
	type tester interface {
		InvalidFormat() bool
	}
	i, ok := err.(tester)
	if !ok {
		t.Fatal("unimplemented InvalidFormat method")
	}
	if !i.InvalidFormat() {
		t.Errorf("unexpexted value expected: %t, give: %t", true, i.InvalidFormat())
	}
}

func TestNewInvalidCRC(t *testing.T) {
	var expected, got uint16 = 1, 2
	err := mpm.NewInvalidCRC(expected, got)
	if err == nil {
		t.Fatal("should not be nil")
	}
	msg := fmt.Sprintf("mpm: expected CRC is %x not %x", expected, got)
	if err.Error() != msg {
		t.Errorf("unexpexted error message expected: %s, give: %s", msg, err.Error())
	}
	type tester interface {
		InvalidCRC() bool
	}
	i, ok := err.(tester)
	if !ok {
		t.Fatal("unimplemented InvalidFormat method")
	}
	if !i.InvalidCRC() {
		t.Errorf("unexpexted value expected: %t, give: %t", true, i.InvalidCRC())
	}
}
