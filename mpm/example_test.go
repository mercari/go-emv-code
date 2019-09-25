package mpm_test

import (
	"fmt"
	"log"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/tlv"
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
		UnreservedTemplates: []tlv.TLV{
			{Tag: "80", Length: "36", Value: "003239401ff0c21a4543a8ed5fbaa30ab02e"},
		},
	}

	buf, err := mpm.Encode(&c)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := mpm.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", dst)

	// Output:
	// &{PayloadFormatIndicator:01 PointOfInitiationMethod:12 MerchantAccountInformation:[] MerchantCategoryCode:4111 TransactionCurrency:156 TransactionAmount:{String: Valid:false} TipOrConvenienceIndicator: ValueOfConvenienceFeeFixed:{String: Valid:false} ValueOfConvenienceFeePercentage:{String: Valid:false} CountryCode:CN MerchantName:BEST TRANSPORT MerchantCity:BEIJING PostalCode: AdditionalDataFieldTemplate:030412340603***0708A60086670902ME MerchantInformation:{LanguagePreference: Name: City: Valid:false} UnreservedTemplates:[{Tag:80 Length:36 Value:003239401ff0c21a4543a8ed5fbaa30ab02e}]}
}
