package mpm_test

import (
	"fmt"
	"log"

	"github.com/mercari/go-emv-code/mpm"
)

func ExampleDecode() {
	c := mpm.Code{
		PayloadFormatIndicator:      "01",
		PointOfInitiationMethod:     mpm.PointOfInitiationMethodDynamic,
		MerchantCategoryCode:        "4111",
		TransactionCurrency:         "156",
		CountryCode:                 "CN",
		MerchantName:                "BEST TRANSPORT",
		MerchantCity:                "BEIJING",
		PostalCode:                  "",
		AdditionalDataFieldTemplate: "030412340603***0708A60086670902ME",
	}

	buf, err := mpm.Encode(&c)
	if err != nil {
		log.Fatal(err)
		return
	}

	dst, err := mpm.Decode(buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v\n", dst)

	// Output:
	// &{PayloadFormatIndicator:01 PointOfInitiationMethod:12 MerchantAccountInformation:[] MerchantCategoryCode:4111 TransactionCurrency:156 TransactionAmount:{String: Valid:false} TipOrConvenienceIndicator: ValueOfConvenienceFeeFixed:{String: Valid:false} ValueOfConvenienceFeePercentage:{String: Valid:false} CountryCode:CN MerchantName:BEST TRANSPORT MerchantCity:BEIJING PostalCode: AdditionalDataFieldTemplate:030412340603***0708A60086670902ME MerchantInformation:{LanguagePreference: Name: City: Valid:false}}
}
