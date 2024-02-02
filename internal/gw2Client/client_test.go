package gw2Client

import (
	"reflect"
	"testing"
)

func TestClient_ListAllitens(t *testing.T) {
	tests := []struct {
		name    string
		want    []int
		wantErr bool
	}{
		{
			name:    "asd",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{}
			got, err := c.ListAllIds()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAllIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAllIds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
