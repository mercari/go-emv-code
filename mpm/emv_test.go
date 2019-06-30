package mpm_test

import (
	"io"
	"reflect"
	"testing"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/tlv"
)

var emvSamplePayload = []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A")

func TestDecoder_Decode(t *testing.T) {
	type invalidFormat interface {
		InvalidFormat() bool
	}

	type invalidCRC interface {
		InvalidCRC() bool
	}

	type args struct {
		buf []byte
	}
	tests := []struct {
		name            string
		args            args
		want            *mpm.Code
		wantErr         bool
		wantErrTypeFunc func(error) bool
	}{
		{
			name: "pass",
			args: args{
				buf: emvSamplePayload,
			},
			want: &mpm.Code{
				PayloadFormatIndicator:  "01",
				PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
				MerchantAccountInformation: []tlv.TLV{
					{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
					{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
				},
				MerchantCategoryCode: "4111",
				TransactionCurrency:  "156",
				TransactionAmount: mpm.NullString{
					String: "23.72",
					Valid:  true,
				},
				TipOrConvenienceIndicator:   mpm.TipOrConvenienceIndicatorPrompt,
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
		{
			name: "pass: without MerchantInformation",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304FF8B"),
			},
			want: &mpm.Code{
				PayloadFormatIndicator:  "01",
				PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
				MerchantAccountInformation: []tlv.TLV{
					{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
					{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
				},
				MerchantCategoryCode: "4111",
				TransactionCurrency:  "156",
				TransactionAmount: mpm.NullString{
					String: "23.72",
					Valid:  true,
				},
				TipOrConvenienceIndicator:   mpm.TipOrConvenienceIndicatorPrompt,
				CountryCode:                 "CN",
				MerchantName:                "BEST TRANSPORT",
				MerchantCity:                "BEIJING",
				PostalCode:                  "",
				AdditionalDataFieldTemplate: "030412340603***0708A60086670902ME",
			},
		},
		{
			name: "err: invalid header",
			args: args{
				buf: []byte("01021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
		{
			name: "err: missing CRC",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A011223344998877070812345678"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
		{
			name: "err: CRC is not valid",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A01122334499887707081234567863040000"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidCRC)
				return ok && e.InvalidCRC()
			},
		},
		{
			name: "err: too short payload",
			args: args{
				buf: []byte("foo"),
			},
			wantErr: true,
		},
		{
			name: "err: malformed payload",
			args: args{
				buf: []byte("00020101021126680019jp.or.paymentsjapan0113aaaaaaaaaaaaa0204bbbb0306cccccc0406dddddd5204000053033925802JP5925123456789012345678901234560151234567890123456110123456786499あいうえおかきくけこ\bあいうえおかきくけこあいうえおかきくけこあいうえおかきくけこあいうえおかきくけこあいうえおかきくけこ\bあいうえおかきくけこあいうえおかきくけこあいうえおかきくけこあいうえおかきくけ6304E0C6"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst, err := mpm.Decode(tt.args.buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErrTypeFunc != nil && !tt.wantErrTypeFunc(err) {
				t.Errorf("Decoder.Decode() unexpected error passed error = %v", err)
			}

			if tt.want != nil && !reflect.DeepEqual(dst, tt.want) {
				t.Errorf("Decoder.Decode() = %v, want %v", dst, tt.want)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := mpm.Decode(emvSamplePayload); err != nil {
			b.Error(err)
		}
	}
}

func TestEncoder_Encode(t *testing.T) {
	type fields struct {
		r *io.Writer
	}
	type args struct {
		in *mpm.Code
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            []byte
		wantErr         bool
		wantErrTypeFunc func(error) bool
	}{
		{
			name: "pass",
			args: args{
				in: &mpm.Code{
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
						{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
					},
					MerchantCategoryCode: "4111",
					TransactionCurrency:  "156",
					TransactionAmount: mpm.NullString{
						String: "23.72",
						Valid:  true,
					},
					TipOrConvenienceIndicator:   mpm.TipOrConvenienceIndicatorPrompt,
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
			want: emvSamplePayload,
		},
		{
			name: "pass: without MerchantInformation",
			args: args{
				in: &mpm.Code{
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: mpm.PointOfInitiationMethodDynamic,
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
						{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
					},
					MerchantCategoryCode: "4111",
					TransactionCurrency:  "156",
					TransactionAmount: mpm.NullString{
						String: "23.72",
						Valid:  true,
					},
					TipOrConvenienceIndicator:   mpm.TipOrConvenienceIndicatorPrompt,
					CountryCode:                 "CN",
					MerchantName:                "BEST TRANSPORT",
					MerchantCity:                "BEIJING",
					PostalCode:                  "",
					AdditionalDataFieldTemplate: "030412340603***0708A60086670902ME",
				},
			},
			want: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304FF8B"),
		},
		{
			name:    "err: cannot pass nil pointer",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := mpm.Encode(tt.args.in)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErrTypeFunc != nil && !tt.wantErrTypeFunc(err) {
				t.Errorf("Encoder.Encode() unexpected error passed error = %v", err)
			}

			if tt.want != nil {
				inDst, err := mpm.Decode(buf)
				if err != nil {
					t.Errorf("unexpected error = %v", err)
				}

				wantDst, err := mpm.Decode(tt.want)
				if err != nil {
					t.Errorf("unexpected error = %v", err)
				}

				if !reflect.DeepEqual(inDst, wantDst) {
					t.Errorf("Encoder.Encode() = %v, want %v", inDst, wantDst)
				}
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	code := &mpm.Code{
		PayloadFormatIndicator:      "01",
		PointOfInitiationMethod:     mpm.PointOfInitiationMethodDynamic,
		MerchantCategoryCode:        "4111",
		TransactionCurrency:         "156",
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
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := mpm.Encode(code); err != nil {
			b.Error(err)
		}
	}
}
