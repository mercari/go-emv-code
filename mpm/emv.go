/*
Package mpm implements encoding and decoding of EMV as defined in EMV Payment Code.
*/
package mpm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mercari/go-emv-code/crc16"
	"github.com/mercari/go-emv-code/tlv"
)

// Code represents EMV Payment Code payload structure.
type Code struct {
	PayloadFormatIndicator          string                    `emv:"00"` // The first data object
	PointOfInitiationMethod         PointOfInitiationMethod   `emv:"01"`
	MerchantAccountInformation      []tlv.TLV                 `emv:"MerchantAccountInformation"`
	MerchantCategoryCode            string                    `emv:"52"`
	TransactionCurrency             string                    `emv:"53"`
	TransactionAmount               NullString                `emv:"54"`
	TipOrConvenienceIndicator       TipOrConvenienceIndicator `emv:"55"`
	ValueOfConvenienceFeeFixed      NullString                `emv:"56"`
	ValueOfConvenienceFeePercentage NullString                `emv:"57"`
	CountryCode                     string                    `emv:"58"`
	MerchantName                    string                    `emv:"59"`
	MerchantCity                    string                    `emv:"60"`
	PostalCode                      string                    `emv:"61"`
	AdditionalDataFieldTemplate     string                    `emv:"62"`
	// CRC                             string  `emv:"63"` // The last object under the root. But useless for value.
	MerchantInformation NullMerchantInformation `emv:"64"`
}

const (
	tagName   = "emv"
	tagLength = 2
	lenLength = 2

	// MaxSize represents max size of EMV Payment Code payload.
	MaxSize = 512
)

const (
	payloadFormatIndicatorID  = "00"
	payloadFormatIndicator    = payloadFormatIndicatorID + "0201"
	payloadFormatIndicatorLen = len(payloadFormatIndicator)

	crcID           = "63"
	crcIDLengthRepr = crcID + "04"
	crcValueLen     = 4
	crcLen          = len(crcIDLengthRepr) + crcValueLen

	merchantAccountInformationIDFrom  = 2
	merchantAccountInformationIDTo    = 51
	merchantAccountInformationTagName = "MerchantAccountInformation"
)

func merchantAccountInformation(tag, length []rune) ([]rune, []rune) {
	id, _ := strconv.Atoi(string(tag))
	if (id >= merchantAccountInformationIDFrom) && (id <= merchantAccountInformationIDTo) {
		return []rune(merchantAccountInformationTagName), length
	}
	return tag, length
}

// ValidatorFunc is an adapter for functions as validator.
type ValidatorFunc func(*Code) error

// Decode decodes payload and validates as EMV MPM.
func Decode(payload []byte, vfs ...ValidatorFunc) (*Code, error) {
	l := len(payload)
	if l < crcLen {
		return nil, NewInvalidFormat("mpm: too short payload")
	}

	if string(payload[:payloadFormatIndicatorLen]) != payloadFormatIndicator {
		return nil, NewInvalidFormat(fmt.Sprintf("mpm: first %d bytes should be match %s", payloadFormatIndicatorLen, payloadFormatIndicator))
	}

	if string(payload[l-crcLen:l-crcValueLen]) != crcIDLengthRepr {
		return nil, NewInvalidFormat(fmt.Sprintf("mpm: last %d bytes should be represents CRC. got %s", crcLen, string(payload[l-crcLen:])))
	}

	crc := crc16.ChecksumCCITTFalse([]byte(string(payload[:l-crcValueLen])))
	if got, _ := strconv.ParseUint(string(payload[l-crcValueLen:l]), 16, 64); uint16(got) != crc {
		return nil, NewInvalidCRC(crc, uint16(got))
	}

	var c Code
	if err := tlv.NewDecoder(bytes.NewReader(payload), tagName, MaxSize, tagLength, lenLength, tlv.TagLengthTranslatorFunc(merchantAccountInformation)).Decode(&c); err != nil {
		switch e := err.(type) {
		case *tlv.MalformedPayloadError:
			return nil, NewInvalidFormat(fmt.Sprintf("mpm: %s", e.Error()))
		}
		return nil, err
	}

	vfs = append(vfs, validateMerchantInformation)
	for _, f := range vfs {
		if err := f(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func merchantAccountInformationTagLengthTranslator(tag, length []rune) ([]rune, []rune) {
	if string(tag) == merchantAccountInformationTagName {
		tag = []rune{}
		length = []rune{}
	}
	return tag, length
}

// Encode encodes to EMV Payment Code payload.
func Encode(c *Code, vfs ...ValidatorFunc) ([]byte, error) {
	if c == nil {
		return nil, errors.New("mpm: nil is not allowed")
	}

	vfs = append(vfs, validateMerchantInformation)
	for _, f := range vfs {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	hash := crc16.NewCCITTFalse()

	var buf bytes.Buffer
	buf.Grow(MaxSize)

	w := io.MultiWriter(&buf, hash)

	if _, err := w.Write([]byte(payloadFormatIndicator)); err != nil {
		return nil, fmt.Errorf("mpm: failed to write PayloadFormatIndicator: %s", err)
	}

	if err := tlv.NewEncoder(w, tagName, []string{payloadFormatIndicatorID, crcID}, tlv.TagLengthTranslatorFunc(merchantAccountInformationTagLengthTranslator)).Encode(c); err != nil {
		return nil, fmt.Errorf("mpm: failed to encode: %s", err)
	}

	// To calculate CRC, we need the ID and Length of the CRC itself.
	if _, err := w.Write([]byte(crcIDLengthRepr)); err != nil {
		return nil, fmt.Errorf("mpm: failed to write CRC header: %s", err)
	}

	crc := strings.ToUpper(fmt.Sprintf("%04x", hash.Sum16()))
	if _, err := w.Write([]byte(crc)); err != nil {
		return nil, fmt.Errorf("mpm: failed to write CRC: %s", err)
	}

	return buf.Bytes(), nil
}

func validateMerchantInformation(c *Code) error {
	if !c.MerchantInformation.Valid {
		return nil
	}
	if len(c.MerchantInformation.LanguagePreference) != 2 {
		return NewInvalidFormat("mpm: length of MerchantInformation.LanguagePreference should be 2")
	}
	if c.MerchantInformation.Name == "" || 25 < utf8.RuneCountInString(c.MerchantInformation.Name) {
		return NewInvalidFormat("mpm: length of MerchantInformation.Name should be between 1 and 25")
	}
	if 15 < utf8.RuneCountInString(c.MerchantInformation.City) {
		return NewInvalidFormat("mpm: length of MerchantInformation.City should be less than 15")
	}
	return nil
}
