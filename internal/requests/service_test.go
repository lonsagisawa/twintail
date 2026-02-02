package requests

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestStoreServiceRequest_Validation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     StoreServiceRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "https",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "missing service name",
			req: StoreServiceRequest{
				Protocol:    "https",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "ftp",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "non-numeric port",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "https",
				ExposePort:  "abc",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "missing destination",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "https",
				ExposePort:  "443",
			},
			wantErr: true,
		},
		{
			name: "tcp protocol",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "tcp",
				ExposePort:  "5432",
				Destination: "localhost:5432",
			},
			wantErr: false,
		},
		{
			name: "tcp+tls protocol",
			req: StoreServiceRequest{
				ServiceName: "my-service",
				Protocol:    "tcp+tls",
				ExposePort:  "5432",
				Destination: "localhost:5432",
			},
			wantErr: false,
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

func TestStoreServiceRequest_Default(t *testing.T) {
	var req StoreServiceRequest
	defaults := req.Default()

	if defaults.Protocol != "https" {
		t.Errorf("expected Protocol 'https', got '%s'", defaults.Protocol)
	}
	if defaults.ExposePort != "443" {
		t.Errorf("expected ExposePort '443', got '%s'", defaults.ExposePort)
	}
}

func TestStoreServiceRequest_ToParams(t *testing.T) {
	req := StoreServiceRequest{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}

	params := req.ToParams()

	if params.ServiceName != req.ServiceName {
		t.Errorf("expected ServiceName '%s', got '%s'", req.ServiceName, params.ServiceName)
	}
	if params.Protocol != req.Protocol {
		t.Errorf("expected Protocol '%s', got '%s'", req.Protocol, params.Protocol)
	}
	if params.ExposePort != req.ExposePort {
		t.Errorf("expected ExposePort '%s', got '%s'", req.ExposePort, params.ExposePort)
	}
	if params.Destination != req.Destination {
		t.Errorf("expected Destination '%s', got '%s'", req.Destination, params.Destination)
	}
}
