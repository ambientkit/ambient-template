// Package navigation provides a navigation page for an Ambient app.
package navigation

import (
	"context"
	"embed"

	"github.com/ambientkit/ambient"
)

//go:embed template/*.tmpl
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
}

// New returns a new navigation plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName(context.Context) string {
	return "navigation"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion(context.Context) string {
	return "1.0.0"
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests(context.Context) []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantRouterRouteWrite, Description: "Access to create routes for javascript."},
	}
}

// Routes sets routes for the plugin.
func (p *Plugin) Routes(context.Context) {
	p.Mux.Get("/dashboard/plugins/navigation", p.index)
}
