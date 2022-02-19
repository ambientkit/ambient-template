// Package mvp provides a template for building a plugin for Ambient apps.
package mvp

import (
	"github.com/ambientkit/ambient"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	// PluginBase is a struct that provides empty functions to satisfy the
	// `Plugin` interface in ambient.go. This allows plugins to work with newer
	// versions of the interface as long as functions defined below still match
	// the interface.
	*ambient.PluginBase
}

// New returns a new mvp plugin. This should be modified to include any values
// or objects that need to be passed in before it's enabled. Any example would
// be a password hash or a flag for debug mode.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name. This name should match the package name.
// PluginName should be globally unique. Only lowercase letters, numbers,
// and underscores are permitted. Must start with with a letter.
func (p *Plugin) PluginName() string {
	return "mvp"
}

// PluginVersion returns the plugin version. This version must follow
// https://semver.org/.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}
