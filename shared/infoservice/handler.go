package infoservice

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"wz2100.net/microlobby/shared/component"
	"wz2100.net/microlobby/shared/proto/infoservicepb/v1"
)

// Handler is the handler for wz2100.net/microlobby/shared/proto/infoservice
type Handler struct {
	comRegistry *component.Registry
	proxyURI    string
	apiVersion  string
	routes      []*infoservicepb.RoutesReply_Route
}

// NewHandler returns a new srv/user pb handler
func NewHandler(comRegistry *component.Registry, proxyURI, apiVersion string, routes []*infoservicepb.RoutesReply_Route) *Handler {
	return &Handler{comRegistry, proxyURI, apiVersion, routes}
}

// Health returns information about the health of this service.
func (h *Handler) Health(ctx context.Context, req *empty.Empty, rsp *infoservicepb.HealthReply) error {
	healthMap := h.comRegistry.Health(ctx)

	hasError := false

	rsp.Infos = make(map[string]*infoservicepb.HealthReply_HealthInfo)
	for name, info := range healthMap {
		if info.IsError {
			hasError = true
		}

		rsp.Infos[name] = &infoservicepb.HealthReply_HealthInfo{
			Message: info.Message,
			IsError: info.IsError,
		}
	}

	rsp.HasError = hasError

	return nil
}

// Routes returns the registered routes
func (h *Handler) Routes(ctx context.Context, req *empty.Empty, rsp *infoservicepb.RoutesReply) error {
	rsp.ProxyURI = h.proxyURI
	rsp.ApiVersion = h.apiVersion
	rsp.Routes = h.routes

	return nil
}
