package mpm_test

import (
	"reflect"
	"strings"
	"testing"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/tlv"
)

var emvSamplePayload = []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京8036003239401ff0c21a4543a8ed5fbaa30ab02e81360032c2fbf6dd646f4f36b617f10747c0b96163046F32")

type invalidFormat interface {
	InvalidFormat() bool
}

func TestDecoder_Decode(t *testing.T) {
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
				UnreservedTemplates: []tlv.TLV{
					{Tag: "80", Length: "36", Value: "003239401ff0c21a4543a8ed5fbaa30ab02e"},
					{Tag: "81", Length: "36", Value: "0032c2fbf6dd646f4f36b617f10747c0b961"},
				},
			},
		},
		{
			name: "pass: without MerchantInformation",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME8036003239401ff0c21a4543a8ed5fbaa30ab02e81360032c2fbf6dd646f4f36b617f10747c0b9616304EBCA"),
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
				UnreservedTemplates: []tlv.TLV{
					{Tag: "80", Length: "36", Value: "003239401ff0c21a4543a8ed5fbaa30ab02e"},
					{Tag: "81", Length: "36", Value: "0032c2fbf6dd646f4f36b617f10747c0b961"},
				},
			},
		},
		{
			name: "pass: without UnreservedTemplates",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京6304CDE7"),
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
			name: "err: GloballyUniqueIdentifier of UnreservedTemplate is empty",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京8000630484D1"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
		{
			name: "err: GloballyUniqueIdentifier of UnreservedTemplate is greater than 32",
			args: args{
				buf: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京8037003339401ff0c21a4543a8ed5fbaa30ab02ee6304D279"),
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
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
		tt := tt
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
	type args struct {
		in *mpm.Code
	}
	tests := []struct {
		name            string
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
					UnreservedTemplates: []tlv.TLV{
						{Tag: "80", Length: "36", Value: "003239401ff0c21a4543a8ed5fbaa30ab02e"},
						{Tag: "81", Length: "36", Value: "0032c2fbf6dd646f4f36b617f10747c0b961"},
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
					UnreservedTemplates: []tlv.TLV{
						{Tag: "80", Length: "36", Value: "003239401ff0c21a4543a8ed5fbaa30ab02e"},
						{Tag: "81", Length: "36", Value: "0032c2fbf6dd646f4f36b617f10747c0b961"},
					},
				},
			},
			want: []byte("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME8036003239401ff0c21a4543a8ed5fbaa30ab02e81360032c2fbf6dd646f4f36b617f10747c0b9616304EBCA"),
		},
		{
			name: "err: err: GloballyUniqueIdentifier of UnreservedTemplate is empty",
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
					UnreservedTemplates: []tlv.TLV{
						{Tag: "80", Length: "00", Value: ""},
					},
				},
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
		{
			name: "err: GloballyUniqueIdentifier of UnreservedTemplate is greater than 32",
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
					UnreservedTemplates: []tlv.TLV{
						{Tag: "80", Length: "37", Value: "003339401ff0c21a4543a8ed5fbaa30ab02e"},
					},
				},
			},
			wantErr: true,
			wantErrTypeFunc: func(err error) bool {
				e, ok := err.(invalidFormat)
				return ok && e.InvalidFormat()
			},
		},
		{
			name:    "err: cannot pass nil pointer",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf, err := mpm.Encode(tt.args.in)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErrTypeFunc != nil && !tt.wantErrTypeFunc(err) {
				t.Errorf("Encoder.Encode() unexpected error passed error = %v", err)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(buf, tt.want) {
					t.Errorf("Encoder.Encode() = %v, want %v", buf, tt.want)
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

func TestTlvDecode(t *testing.T) {
	type args struct {
		payload string
		bufSize int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "read tag error",
			args: args{
				payload: "0001a",
				bufSize: 1,
			},
			wantErr: true,
		},
		{
			name: "read length error",
			args: args{
				payload: "0001a",
				bufSize: 3,
			},
			wantErr: true,
		},
		{
			name: "read value error",
			args: args{
				payload: "0001a",
				bufSize: 4,
			},
			wantErr: true,
		},
		{
			name: "pass",
			args: args{
				payload: "0001a",
				bufSize: 5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var v struct{}
			err := tlv.NewDecoder(strings.NewReader(tt.args.payload), "emv", tt.args.bufSize, 2, 2, nil).Decode(&v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decoer.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
