package tlv

import (
	"strings"
	"testing"
)

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
				payload: "003239401ff0c21a4543a8ed5fbaa30ab02e",
				bufSize: 1,
			},
			wantErr: true,
		},
		{
			name: "read length error",
			args: args{
				payload: "003239401ff0c21a4543a8ed5fbaa30ab02e",
				bufSize: 3,
			},
			wantErr: true,
		},
		{
			name: "read value error",
			args: args{
				payload: "003239401ff0c21a4543a8ed5fbaa30ab02e",
				bufSize: 4,
			},
			wantErr: true,
		},
		{
			name: "pass",
			args: args{
				payload: "003239401ff0c21a4543a8ed5fbaa30ab02e",
				bufSize: 36,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var v struct{}
			err := NewDecoder(strings.NewReader(tt.args.payload), "emv", tt.args.bufSize, 2, 2, nil).Decode(&v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
