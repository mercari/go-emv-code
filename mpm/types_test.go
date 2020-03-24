package mpm_test

import (
	"testing"

	"go.mercari.io/go-emv-code/mpm"
)

func TestTipOrConvenienceIndicator_Tokenize(t *testing.T) {
	tests := []struct {
		name    string
		give    *mpm.TipOrConvenienceIndicator
		want    string
		wantErr bool
	}{
		{
			name:    "give nil",
			give:    nil,
			want:    "",
			wantErr: false,
		},
		{
			name: "give mpm.TipOrConvenienceIndicatorPrompt",
			give: func() *mpm.TipOrConvenienceIndicator {
				p := mpm.TipOrConvenienceIndicatorPrompt
				return &p
			}(),
			want:    "01",
			wantErr: false,
		},
		{
			name: "give mpm.TipOrConvenienceIndicatorFixed",
			give: func() *mpm.TipOrConvenienceIndicator {
				p := mpm.TipOrConvenienceIndicatorFixed
				return &p
			}(),
			want:    "02",
			wantErr: false,
		},
		{
			name: "give mpm.TipOrConvenienceIndicatorPercentage",
			give: func() *mpm.TipOrConvenienceIndicator {
				p := mpm.TipOrConvenienceIndicatorPercentage
				return &p
			}(),
			want:    "03",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			dst, err := tt.give.Tokenize()

			if (err != nil) != tt.wantErr {
				t.Errorf("TipOrConvenienceIndicator.Tokenize error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != dst {
				t.Errorf("TipOrConvenienceIndicator.Tokenize = %v, want %v", dst, tt.want)
			}
		})
	}
}

func TestTipOrConvenienceIndicator_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    []rune
		wantErr bool
	}{
		{
			name:    "give nil",
			give:    nil,
			wantErr: true,
		},
		{
			name:    "give empty",
			give:    []rune{},
			wantErr: true,
		},
		{
			name:    "give unexpected string",
			give:    []rune("wrong_value"),
			wantErr: true,
		},
		{
			name:    "give mpm.TipOrConvenienceIndicatorPrompt",
			give:    []rune(mpm.TipOrConvenienceIndicatorPrompt),
			wantErr: false,
		},
		{
			name:    "give mpm.TipOrConvenienceIndicatorFixed",
			give:    []rune(mpm.TipOrConvenienceIndicatorFixed),
			wantErr: false,
		},
		{
			name:    "give mpm.TipOrConvenienceIndicatorPercentage",
			give:    []rune(mpm.TipOrConvenienceIndicatorPercentage),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var ind mpm.TipOrConvenienceIndicator
			err := ind.Scan(tt.give)

			if (err != nil) != tt.wantErr {
				t.Errorf("TipOrConvenienceIndicator.Scan error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
