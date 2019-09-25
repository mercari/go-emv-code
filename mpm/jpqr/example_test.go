package jpqr_test

import (
	"fmt"
	"log"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/mpm/jpqr"
	"go.mercari.io/go-emv-code/tlv"
)

func ExampleDecode() {
	c := mpm.Code{
		PayloadFormatIndicator:  "01",
		PointOfInitiationMethod: mpm.PointOfInitiationMethodStatic,
		MerchantAccountInformation: []tlv.TLV{
			{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
			{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
			{Tag: "26", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
		},
		MerchantCategoryCode: "5812",
		TransactionCurrency:  "392",
		CountryCode:          "JP",
		MerchantName:         "xxx",
		MerchantCity:         "xxx",
		PostalCode:           "1066143",
		MerchantInformation: mpm.NullMerchantInformation{
			LanguagePreference: "JA",
			Name:               "メルペイ カフェ",
			Valid:              true,
		},
	}

	buf, err := jpqr.Encode(&c)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := jpqr.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", dst)

	// Output:
	// &{PayloadFormatIndicator:01 PointOfInitiationMethod:11 MerchantAccountInformation:[{Tag:29 Length:30 Value:0012D156000000000510A93FO3230Q} {Tag:31 Length:28 Value:0012D15600000001030812345678} {Tag:26 Length:68 Value:0019jp.or.paymentsjapan011300000000000010204000103060000010406000001}] MerchantCategoryCode:5812 TransactionCurrency:392 TransactionAmount:{String: Valid:false} TipOrConvenienceIndicator: ValueOfConvenienceFeeFixed:{String: Valid:false} ValueOfConvenienceFeePercentage:{String: Valid:false} CountryCode:JP MerchantName:xxx MerchantCity:xxx PostalCode:1066143 AdditionalDataFieldTemplate: MerchantInformation:{LanguagePreference:JA Name:メルペイ カフェ City: Valid:true} UnreservedTemplates:[]}
}
