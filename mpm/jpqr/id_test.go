package jpqr_test

import (
	"reflect"
	"testing"

	"go.mercari.io/go-emv-code/mpm"
	"go.mercari.io/go-emv-code/mpm/jpqr"
	"go.mercari.io/go-emv-code/tlv"
)

func TestParseID(t *testing.T) {
	type args struct {
		c *mpm.Code
	}
	tests := []struct {
		name    string
		args    args
		want    *jpqr.ID
		wantErr bool
	}{
		{
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "26", Length: "68", Value: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001"},
						{Tag: "29", Length: "30", Value: "0012D156000000000510A93FO3230Q"},
						{Tag: "31", Length: "28", Value: "0012D15600000001030812345678"},
					},
				},
			},
			want: &jpqr.ID{"jp.or.paymentsjapan", "0000000000001", "0001", "000001", "000001"},
		},
		{
			name: "fail: malformed payload",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "26", Length: "68", Value: "foobarbaz"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid prefix",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "27", Length: "68", Value: "0019jp.co.paymentsjapan011300000000000010204000103060000010406000001"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV1)",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "28", Length: "67", Value: "0019jp.or.paymentsjapan01120000000000010204000103060000010406000001"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV2)",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "29", Length: "67", Value: "0019jp.or.paymentsjapan01130000000000001020300003060000010406000001"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV3)",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "30", Length: "67", Value: "0019jp.or.paymentsjapan01130000000000001020400010305000000406000001"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV4)",
			args: args{
				c: &mpm.Code{
					MerchantAccountInformation: []tlv.TLV{
						{Tag: "31", Length: "67", Value: "0019jp.or.paymentsjapan01130000000000001020400010306000001040500000"},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpqr.ParseID(tt.args.c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseID() got = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseIDFromString(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    *jpqr.ID
		wantErr bool
	}{
		{
			args: args{
				src: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001",
			},
			want: &jpqr.ID{Prefix: "jp.or.paymentsjapan", LV1: "0000000000001", LV2: "0001", LV3: "000001", LV4: "000001"},
		},
		{
			name: "fail: malformed payload",
			args: args{
				src: "foobarbaz",
			},
			wantErr: true,
		},
		{
			name: "fail: empty src is err",
			args: args{
				src: "",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid prefix",
			args: args{
				src: "0019jp.co.paymentsjapan011300000000000010204000103060000010406000001",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV1)",
			args: args{
				src: "0019jp.or.paymentsjapan01130000000000001020300003060000010406000001",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV2)",
			args: args{
				src: "0019jp.or.paymentsjapan01130000000000001020300003060000010406000001",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV3)",
			args: args{
				src: "0019jp.or.paymentsjapan01130000000000001020400010305000000406000001",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid length (LV4)",
			args: args{
				src: "0019jp.or.paymentsjapan01130000000000001020400010306000001040500000",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpqr.ParseIDFromString(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseIDFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIDFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_String(t *testing.T) {
	type fields struct {
		Prefix string
		LV1    string
		LV2    string
		LV3    string
		LV4    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Prefix: "jp.or.paymentsjapan",
				LV1:    "0000000000001",
				LV2:    "0001",
				LV3:    "000001",
				LV4:    "000001",
			},
			want: "0019jp.or.paymentsjapan011300000000000010204000103060000010406000001",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			i := &jpqr.ID{
				Prefix: tt.fields.Prefix,
				LV1:    tt.fields.LV1,
				LV2:    tt.fields.LV2,
				LV3:    tt.fields.LV3,
				LV4:    tt.fields.LV4,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("ID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
