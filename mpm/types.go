package mpm

import (
	"fmt"
	"strings"

	"go.mercari.io/go-emv-code/tlv"
)

// NullMerchantInformation represents Data Objects for Merchant Information—Language Template.
type NullMerchantInformation struct {
	LanguagePreference string `emv:"00"`
	Name               string `emv:"01"`
	City               string `emv:"02"`
	Valid              bool
}

// Tokenize turns NullMerchantInformation into a string
func (m *NullMerchantInformation) Tokenize() (string, error) {
	if m == nil {
		return "", nil
	}
	if !m.Valid {
		return "", nil
	}
	var buf strings.Builder
	if err := tlv.NewEncoder(&buf, tagName, nil, nil).Encode(m); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (m *NullMerchantInformation) Scan(token []rune) error {
	var mm NullMerchantInformation
	if err := tlv.NewDecoder(strings.NewReader(string(token)), tagName, MaxSize, tagLength, lenLength, nil).Decode(&mm); err != nil {
		return err
	}
	mm.Valid = mm.LanguagePreference != "" && mm.Name != ""
	*m = mm
	return nil
}

// PointOfInitiationMethod represents Data Objects for Point of Initiation Method.
type PointOfInitiationMethod string

const (
	PointOfInitiationMethodStatic  PointOfInitiationMethod = "11"
	PointOfInitiationMethodDynamic PointOfInitiationMethod = "12"
)

// Tokenize turns PointOfInitiationMethod into a string
func (p *PointOfInitiationMethod) Tokenize() (string, error) {
	if p == nil {
		return "", nil
	}
	return string(*p), nil
}

func (p *PointOfInitiationMethod) Scan(token []rune) error {
	switch PointOfInitiationMethod(string(token)) {
	case PointOfInitiationMethodStatic:
		*p = PointOfInitiationMethodStatic
		return nil
	case PointOfInitiationMethodDynamic:
		*p = PointOfInitiationMethodDynamic
		return nil
	}
	return fmt.Errorf("passed value is invalid for PointOfInitiationMethod: %v", token)
}

// NullString represents a string that may be null.
type NullString struct {
	String string
	Valid  bool
}

// Tokenize turns NullString into a string
func (n *NullString) Tokenize() (string, error) {
	if n == nil || !n.Valid {
		return "", nil
	}
	return n.String, nil
}

func (n *NullString) Scan(token []rune) error {
	nn := NullString{
		String: string(token),
		Valid:  true,
	}
	*n = nn
	return nil
}

// TipOrConvenienceIndicator represents Data Objects for Tip or Convenience Indicator.
type TipOrConvenienceIndicator string

const (
	TipOrConvenienceIndicatorPrompt     TipOrConvenienceIndicator = "01"
	TipOrConvenienceIndicatorFixed      TipOrConvenienceIndicator = "02"
	TipOrConvenienceIndicatorPercentage TipOrConvenienceIndicator = "03"
)

// Tokenize turns TipOrConvenienceIndicator into a string
func (t *TipOrConvenienceIndicator) Tokenize() (string, error) {
	if t == nil {
		return "", nil
	}
	return string(*t), nil
}

func (t *TipOrConvenienceIndicator) Scan(token []rune) error {
	switch TipOrConvenienceIndicator(string(token)) {
	case TipOrConvenienceIndicatorPrompt:
		*t = TipOrConvenienceIndicatorPrompt
		return nil
	case TipOrConvenienceIndicatorFixed:
		*t = TipOrConvenienceIndicatorFixed
		return nil
	case TipOrConvenienceIndicatorPercentage:
		*t = TipOrConvenienceIndicatorPercentage
		return nil
	}
	return fmt.Errorf("passed value is invalid for TipOrConvenienceIndicator: %v", token)
}
