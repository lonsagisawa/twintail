package handlers

import (
	"twintail/internal/services"
)

type FullTailscaleService interface {
	TailscaleService
	EndpointService
}

type Container struct {
	Service  *ServiceHandler
	Endpoint *EndpointHandler
	Settings *SettingsHandler
}

func NewContainer(tailscale FullTailscaleService) *Container {
	return &Container{
		Service:  NewServiceHandler(tailscale),
		Endpoint: NewEndpointHandler(tailscale),
		Settings: NewSettingsHandler(),
	}
}

func NewContainerWithTailscale(tailscale *services.TailscaleService) *Container {
	return NewContainer(tailscale)
}
