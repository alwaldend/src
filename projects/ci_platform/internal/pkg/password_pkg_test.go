package pkg

import (
	"testing"
)

func TestEncodedHashCompare(t *testing.T) {
	type args struct {
		password    string
		encodedHash string
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantErr   bool
	}{
		{
			"static comparison test",
			args{
				"password",
				"$argon2id$v=19$m=65536,t=3,p=2$UCT6dxcsFxs+GrRXXIfAWQ$R2MDRMgNlVwLnnIivxZMNS9dA7PkPx5ir5+sn2xZmKg",
			},
			true,
			false,
		},
		{
			"dynamic comparison test",
			args{
				"password",
				EncodedHashGenerate("password"),
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, err := EncodedHashCompare(tt.args.password, tt.args.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodedHashCompare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("EncodedHashCompare() = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}
