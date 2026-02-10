package requests

import (
	"fmt"
	"strings"
)

func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name is required")
	}
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("service name must not start with '-'")
	}
	for _, ch := range name {
		switch ch {
		case ';', ' ', '\n', '\r', '`', '\x00':
			return fmt.Errorf("service name contains invalid character")
		}
	}
	return nil
}
