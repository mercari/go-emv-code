package jpqr

import (
	"errors"
	"fmt"
	"strings"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/tlv"
)

const (
	idPrefix = "jp.or.paymentsjapan"

	lv0Length = len(idPrefix)
	lv1Length = 13
	lv2Length = 4
	lv3Length = 6
	lv4Length = 6

	tokenFormat = "%s%02d%s"
)

// ID represents a parsed JPQR-ID.
type ID struct {
	Prefix string `lv:"00"`
	LV1    string `lv:"01"`
	LV2    string `lv:"02"`
	LV3    string `lv:"03"`
	LV4    string `lv:"04"`
}

// String returns the accumulated string.
func (i *ID) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(tokenFormat, "00", lv0Length, idPrefix))
	b.WriteString(fmt.Sprintf(tokenFormat, "01", lv1Length, i.LV1))
	b.WriteString(fmt.Sprintf(tokenFormat, "02", lv2Length, i.LV2))
	b.WriteString(fmt.Sprintf(tokenFormat, "03", lv3Length, i.LV3))
	b.WriteString(fmt.Sprintf(tokenFormat, "04", lv4Length, i.LV4))
	return b.String()
}

func validateIDLength(i *ID) error {
	if len(i.LV1) != lv1Length {
		return fmt.Errorf("len(LV1) should be %d", lv1Length)
	}
	if len(i.LV2) != lv2Length {
		return fmt.Errorf("len(LV2) should be %d", lv2Length)
	}
	if len(i.LV3) != lv3Length {
		return fmt.Errorf("len(LV3) should be %d", lv3Length)
	}
	if len(i.LV4) != lv4Length {
		return fmt.Errorf("len(LV4) should be %d", lv4Length)
	}
	return nil
}

// ParseID validates and parses given *mpm.Code as JPQR-ID.
func ParseID(c *mpm.Code) (*ID, error) {
	for _, v := range c.MerchantAccountInformation {
		var id ID
		if err := tlv.NewDecoder(strings.NewReader(v.Value), "lv", mpm.MaxSize, 2, 2, nil).Decode(&id); err != nil {
			return nil, err
		}
		if id.Prefix == idPrefix {
			if err := validateIDLength(&id); err != nil {
				return nil, err
			}
			return &id, nil
		}
	}
	return nil, errors.New("missing JPQR-ID")
}

// ParseIDFromString validates and parses given string as JPQR-ID.
func ParseIDFromString(v string) (*ID, error) {
	return ParseID(&mpm.Code{
		MerchantAccountInformation: []tlv.TLV{
			{Tag: "26", Length: "68", Value: v},
		},
	})
}
