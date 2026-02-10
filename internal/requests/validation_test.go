package requests

import "testing"

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-service", false},
		{"valid with dots", "my.service.name", false},
		{"empty", "", true},
		{"starts with dash", "-malicious", true},
		{"contains semicolon", "name;rm -rf /", true},
		{"contains space", "my service", true},
		{"contains newline", "name\nmalicious", true},
		{"contains carriage return", "name\rmalicious", true},
		{"contains backtick", "name`id`", true},
		{"contains null byte", "name\x00malicious", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
