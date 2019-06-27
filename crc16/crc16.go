/*
Package crc16 implements the 16-bit cyclic redundancy check, or CRC-16, checksum.
*/
package crc16

import "hash"

// simpleMakeTable allocates and constructs a Table for the specified
// polynomial. The table is suitable for use with the simple algorithm
// (simpleUpdate).
func simpleMakeTable(poly uint16) *Table {
	t := new(Table)
	simplePopulateTable(poly, t)
	return t
}

// simplePopulateTable constructs a Table for the specified polynomial, suitable
// for use with simpleUpdate.
func simplePopulateTable(poly uint16, t *Table) {
	for i := 0; i < 256; i++ {
		crc := uint16(i) << 8
		for j := 0; j < 8; j++ {
			if crc&(1<<15) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t[i] = crc
	}
}

// simpleUpdate uses the simple algorithm to update the CRC, given a table that
// was previously computed using simpleMakeTable.
func simpleUpdate(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

// Predefined polynomials.
const (
	CCITTFalse = 0x1021
)

// CCITTFalseTable is the table for the CCITT-FALSE polynomial.
var CCITTFalseTable = simpleMakeTable(CCITTFalse)

// Table is a 256-word table representing the polynomial for efficient processing.
type Table [256]uint16

// Checksum returns the CRC-16 checksum of data
// using the polynomial represented by the Table.
func Checksum(data []byte, tab *Table) uint16 { return simpleUpdate(0xFFFF, tab, data) }

// ChecksumCCITTFalse returns the CRC-16 checksum of data using the CCITT-FALSE polynomial.
func ChecksumCCITTFalse(data []byte) uint16 { return Checksum(data, CCITTFalseTable) }

type digest struct {
	crc uint16
	tab *Table
}

// NewCCITTFalse creates a new Hash16 computing the CRC-16 checksum using the IEEE polynomial.
func NewCCITTFalse() Hash16 {
	return &digest{0xFFFF, CCITTFalseTable}
}

// The size of a CRC-16 checksum in bytes.
const Size = 4

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Reset() { d.crc = 0 }

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = simpleUpdate(d.crc, d.tab, p)
	return len(p), nil
}

func (d *digest) Sum16() uint16 { return d.crc }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum16()
	return append(in, byte(s>>16), byte(s>>8), byte(s))
}

// Hash16 is the common interface implemented by all 16-bit hash functions.
type Hash16 interface {
	hash.Hash
	Sum16() uint16
}
