package crc16_test

import (
	"testing"

	"go.mercari.io/go-emv-code/crc16"
)

// TestChecksumCCITTFalse tests the generic algorithm.
func TestChecksumCCITTFalse(t *testing.T) {
	tests := []struct {
		in  string
		out uint16
	}{
		// test goldens from https://golang.org/src/hash/crc32/crc32_test.go and http://www.sunshine2k.de/coding/javascript/crc/crc_js.html
		{"a", 0x9D77},
		{"ab", 0x69F0},
		{"abc", 0x514A},
		{"abcd", 0x2CF6},
		{"abcde", 0x2FED},
		{"abcdef", 0x34ED},
		{"abcdefg", 0x8796},
		{"abcdefgh", 0x9AC1},
		{"abcdefghi", 0x1E7C},
		{"abcdefghij", 0x4213},
		{"Discard medicine more than two years old.", 0xB869},
		{"He who has a shady past knows that nice guys finish last.", 0x3BF6},
		{"I wouldn't marry him with a ten foot pole.", 0xF3A8},
		{"Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave", 0x3929},
		{"The days of the digital watch are numbered.  -Tom Stoppard", 0x5DAE},
		{"Nepal premier won't resign.", 0xC724},
		{"For every action there is an equal and opposite government program.", 0x97E0},
		{"His money is twice tainted: 'taint yours and 'taint mine.", 0xC3BF},
		{"There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977", 0x8550},
		{"It's a tiny change to the code and not completely disgusting. - Bob Manchek", 0xC968},
		{"size:  a.out:  bad magic", 0x5FE9},
		{"The major problem is with sendmail.  -Mark Horton", 0xFE1D},
		{"Give me a rock, paper and scissors and I will move the world.  CCFestoon", 0xB2C5},
		{"If the enemy is within range, then so are you.", 0x8838},
		{"It's well we cannot hear the screams/That we create in others' dreams.", 0x5F97},
		{"You remind me of a TV show, but that's all right: I watch it anyway.", 0xD800},
		{"C is as portable as Stonehedge!!", 0xA30C},
		{"Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley", 0x0BF2},
		{"The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule", 0xCF9F},
		{"How can you write a big system without C++?  -Paul Glick", 0xCEEA},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			if crc := crc16.ChecksumCCITTFalse([]byte(tt.in)); crc != tt.out {
				t.Errorf("ChecksumCCITTFalse(%s) = 0x%x want 0x%x", tt.in, crc, tt.out)
			}
		})
	}
}

func TestNewCCITTFalse(t *testing.T) {

	h := crc16.NewCCITTFalse()
	if h == nil {
		t.Error("should not be nil")
	}

	switch h.(type) {
	case crc16.Hash16:
	default:
		t.Error("expect implement crc16.Hash16 interface")
	}
}

func TestDigest_BlockSize(t *testing.T) {

	const crc16BlockSize = 1
	bs := crc16.NewCCITTFalse().BlockSize()
	if bs != crc16BlockSize {
		t.Errorf("give: %d, want: %d", bs, crc16BlockSize)
	}
}

func TestDigest_Reset(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic, catch: %+v", r)
		}
	}()

	crc16.NewCCITTFalse().Reset()
}

func TestDigest_Size(t *testing.T) {

	const crc16Size = 4
	s := crc16.NewCCITTFalse().Size()
	if s != crc16Size {
		t.Errorf("give: %d, want: %d", s, crc16Size)
	}
}

func TestDigest_Sum16(t *testing.T) {
	tests := []struct {
		in  string
		out uint16
	}{
		// test goldens from https://golang.org/src/hash/crc32/crc32_test.go and http://www.sunshine2k.de/coding/javascript/crc/crc_js.html
		{"a", 0x9D77},
		{"ab", 0x69F0},
		{"abc", 0x514A},
		{"abcd", 0x2CF6},
		{"abcde", 0x2FED},
		{"abcdef", 0x34ED},
		{"abcdefg", 0x8796},
		{"abcdefgh", 0x9AC1},
		{"abcdefghi", 0x1E7C},
		{"abcdefghij", 0x4213},
		{"Discard medicine more than two years old.", 0xB869},
		{"He who has a shady past knows that nice guys finish last.", 0x3BF6},
		{"I wouldn't marry him with a ten foot pole.", 0xF3A8},
		{"Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave", 0x3929},
		{"The days of the digital watch are numbered.  -Tom Stoppard", 0x5DAE},
		{"Nepal premier won't resign.", 0xC724},
		{"For every action there is an equal and opposite government program.", 0x97E0},
		{"His money is twice tainted: 'taint yours and 'taint mine.", 0xC3BF},
		{"There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977", 0x8550},
		{"It's a tiny change to the code and not completely disgusting. - Bob Manchek", 0xC968},
		{"size:  a.out:  bad magic", 0x5FE9},
		{"The major problem is with sendmail.  -Mark Horton", 0xFE1D},
		{"Give me a rock, paper and scissors and I will move the world.  CCFestoon", 0xB2C5},
		{"If the enemy is within range, then so are you.", 0x8838},
		{"It's well we cannot hear the screams/That we create in others' dreams.", 0x5F97},
		{"You remind me of a TV show, but that's all right: I watch it anyway.", 0xD800},
		{"C is as portable as Stonehedge!!", 0xA30C},
		{"Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley", 0x0BF2},
		{"The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule", 0xCF9F},
		{"How can you write a big system without C++?  -Paul Glick", 0xCEEA},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			h := crc16.NewCCITTFalse()
			_, err := h.Write([]byte(tt.in))
			if err != nil {
				t.Errorf("unexpected error, err = %s", err)
			}
			if crc := h.Sum16(); crc != tt.out {
				t.Errorf("give: 0x%x, want: 0x%x", crc, tt.out)
			}
		})
	}
}
