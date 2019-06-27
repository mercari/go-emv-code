/*
Package jpqr implements encoding and decoding of JPQR as defined in JPQR MPM Guideline.
*/
package jpqr

import (
	"fmt"

	"github.com/mercari/go-emv-code/mpm"
)

// Decode decodes payload and validates as JPQR.
func Decode(payload []byte) (*mpm.Code, error) {
	c, err := mpm.Decode(payload, []mpm.ValidatorFunc{
		validateID,
		validateCountryCodeIsJP,
		validateTransactionCurrency,
		validatePostalCode,
		validateMerchantInformation,
	}...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Encode encodes to EMV Payment Code payload.
func Encode(c *mpm.Code) ([]byte, error) {
	return mpm.Encode(c, []mpm.ValidatorFunc{
		validateID,
		validateCountryCodeIsJP,
		validateTransactionCurrency,
		validatePostalCode,
		validateMerchantInformation,
	}...)
}

func validateID(c *mpm.Code) error {
	if _, err := ParseID(c); err != nil {
		return mpm.NewInvalidFormat(fmt.Sprintf("jpqr: %s", err))
	}
	return nil
}

const countryCode = "JP"

func validateCountryCodeIsJP(c *mpm.Code) error {
	if c.CountryCode == countryCode {
		return nil
	}
	return mpm.NewInvalidFormat(fmt.Sprintf("jpqr: CountryCode should be %s", countryCode))
}

const transactionCurrency = "392"

func validateTransactionCurrency(c *mpm.Code) error {
	if c.TransactionCurrency == transactionCurrency {
		return nil
	}
	return mpm.NewInvalidFormat(fmt.Sprintf("jpqr: TransactionCurrency should be %s", transactionCurrency))
}

func validatePostalCode(c *mpm.Code) error {
	if c.PostalCode == "" {
		return mpm.NewInvalidFormat("jpqr: PostalCode should be represented")
	}
	return nil
}

const languagePreference = "JA"

func validateMerchantInformation(c *mpm.Code) error {
	if !c.MerchantInformation.Valid {
		return mpm.NewInvalidFormat("jpqr: MerchantInformation should be represented")
	}
	if c.MerchantInformation.LanguagePreference != languagePreference {
		return mpm.NewInvalidFormat(fmt.Sprintf("jpqr: MerchantInformation.LanguagePreference should be %s", languagePreference))
	}
	if c.MerchantInformation.City != "" {
		return mpm.NewInvalidFormat("jpqr: MerchantInformation.City is not necessary")
	}
	return nil
}
