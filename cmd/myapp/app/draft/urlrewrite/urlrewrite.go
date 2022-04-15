// Package urlrewrite removes trailing slash from requests for an Ambient app.
package urlrewrite

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
}

// New returns a new urlrewrite plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName(context.Context) string {
	return "urlrewrite"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion(context.Context) string {
	return "1.0.0"
}

// GrantRequests returns a list of grants requested by the plugin.
func (p *Plugin) GrantRequests(context.Context) []ambient.GrantRequest {
	return []ambient.GrantRequest{
		{Grant: ambient.GrantRouterMiddlewareWrite, Description: "Access to redirect to the correct URL if the user request URL doesn't match."},
	}
}

// Middleware returns router middleware.
func (p *Plugin) Middleware(context.Context) []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.HandlePrefix,
	}
}

// HandlePrefix will handle URLs behind a proxy.
func (p *Plugin) HandlePrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlprefix := os.Getenv("AMB_URL_PREFIX")

		// If there is a prefix, then strip it out for all requests.
		if len(urlprefix) > 0 {
			r.URL.Path = path.Join("/", strings.TrimPrefix(r.URL.Path, urlprefix))
			p.Log.Debug("Rewrote URL: %v", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}
