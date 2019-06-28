package jpqr_test

import (
	"reflect"
	"testing"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/mpm/jpqr"
	"go.mercari.io/go-emv-code/tlv"
)

func TestDecode(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *mpm.Code
		wantErr bool
	}{
		{
			args: args{
				payload: []byte("0002016003xxx01021129300012D156000000000510A93FO3230Q31280012D1560000000103081234567826680019jp.or.paymentsjapan01130000000000001020400010306000001040600000153033925903xxx64180002JA0108メルペイ カフェ520441115802JP610710661436304DEE9"),
			},
			want: &mpm.Code{
				PayloadFormatIndicator:  "01",
				PointOfInitiationMethod: mpm.PointOfInitiationMethodStatic,
				MerchantAccountInformation: []tlv.TLV{
					{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
					{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
					{Tag: "26", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
				},
				MerchantCategoryCode: "4111",
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
			},
		},
		{
			name: "err: crc is invalid",
			args: args{
				payload: []byte("00020101021229300012d156000000000510a93fo3230q31280012d15600000001030812345678520441115802cn5914best transport6007beijing64200002zh0104最佳运输0202北京540523.7253031565502016233030412340603***0708a60086670902me91320016a0112233449988770708123456786304a13b"),
			},
			wantErr: true,
		},
		{
			name: "err: countryCode is not JP",
			args: args{
				payload: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A"),
			},
			wantErr: true,
		},
		{
			name: "err: transactionCurrency is not 392",
			args: args{
				payload: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802JP5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A011223344998877070812345678630420EE"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpqr.Decode(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		code *mpm.Code
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				code: &mpm.Code{
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: mpm.PointOfInitiationMethodStatic,
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
						{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
						{Tag: "26", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
					},
					MerchantCategoryCode: "4111",
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
				},
			},
			want: []byte("0002016003xxx01021129300012D156000000000510A93FO3230Q31280012D1560000000103081234567826680019jp.or.paymentsjapan01130000000000001020400010306000001040600000153033925903xxx64180002JA0108メルペイ カフェ520441115802JP610710661436304DEE9"),
		},
		{
			name: "err: countryCode is not JP",
			args: args{
				code: &mpm.Code{
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "27", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
					},
					MerchantCategoryCode:        "4111",
					TransactionCurrency:         "392",
					CountryCode:                 "CN",
					MerchantName:                "BEST TRANSPORT",
					MerchantCity:                "BEIJING",
					PostalCode:                  "",
					AdditionalDataFieldTemplate: "030412340603***0708A60086670902ME",
					MerchantInformation: mpm.NullMerchantInformation{
						LanguagePreference: "ZH",
						Name:               "最佳运输",
						City:               "北京",
						Valid:              true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "err: transactionCurrency is not 392",
			args: args{
				code: &mpm.Code{
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "28", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
					},
					MerchantCategoryCode:        "4111",
					TransactionCurrency:         "156",
					CountryCode:                 "JP",
					MerchantName:                "BEST TRANSPORT",
					MerchantCity:                "BEIJING",
					PostalCode:                  "",
					AdditionalDataFieldTemplate: "030412340603***0708A60086670902ME",
					MerchantInformation: mpm.NullMerchantInformation{
						LanguagePreference: "ZH",
						Name:               "最佳运输",
						City:               "北京",
						Valid:              true,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpqr.Encode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				gotDst, err := mpm.Decode(got)
				if err != nil {
					t.Errorf("unexpected error = %v", err)
				}

				wantDst, err := mpm.Decode(tt.want)
				if err != nil {
					t.Errorf("unexpected error = %v", err)
				}

				if !reflect.DeepEqual(gotDst, wantDst) {
					t.Errorf("Encode() = %v, want %v", gotDst, wantDst)
				}
			}
		})
	}
}
