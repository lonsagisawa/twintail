package handlers

import (
	"twintail/internal/services"
)

type Container struct {
	Service  *ServiceHandler
	Endpoint *EndpointHandler
	Settings *SettingsHandler
}

func NewContainer(tailscale *services.TailscaleService) *Container {
	return &Container{
		Service:  NewServiceHandler(tailscale),
		Endpoint: NewEndpointHandler(tailscale),
		Settings: NewSettingsHandler(),
	}
}
