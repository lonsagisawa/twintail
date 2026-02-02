package requests

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestUpdateSettingsRequest_Validation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     UpdateSettingsRequest
		wantErr bool
	}{
		{
			name:    "valid english",
			req:     UpdateSettingsRequest{Lang: "en"},
			wantErr: false,
		},
		{
			name:    "valid japanese",
			req:     UpdateSettingsRequest{Lang: "ja"},
			wantErr: false,
		},
		{
			name:    "invalid language",
			req:     UpdateSettingsRequest{Lang: "fr"},
			wantErr: true,
		},
		{
			name:    "empty language",
			req:     UpdateSettingsRequest{Lang: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
