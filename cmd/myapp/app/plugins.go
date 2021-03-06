package app

import (
	"log"
	"os"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient-template/cmd/myapp/app/draft/hello"
	"github.com/ambientkit/ambient-template/cmd/myapp/app/draft/navigation"
	"github.com/ambientkit/plugin/generic/author"
	"github.com/ambientkit/plugin/generic/bearblog"
	"github.com/ambientkit/plugin/generic/bearcss"
	"github.com/ambientkit/plugin/generic/charset"
	"github.com/ambientkit/plugin/generic/debugpprof"
	"github.com/ambientkit/plugin/generic/description"
	"github.com/ambientkit/plugin/generic/disqus"
	"github.com/ambientkit/plugin/generic/envinfo"
	"github.com/ambientkit/plugin/generic/googleanalytics"
	"github.com/ambientkit/plugin/generic/pluginmanager"
	"github.com/ambientkit/plugin/generic/prism"
	"github.com/ambientkit/plugin/generic/robots"
	"github.com/ambientkit/plugin/generic/rssfeed"
	"github.com/ambientkit/plugin/generic/sitemap"
	"github.com/ambientkit/plugin/generic/stackedit"
	"github.com/ambientkit/plugin/generic/styles"
	"github.com/ambientkit/plugin/generic/viewport"
	"github.com/ambientkit/plugin/middleware/cors"
	"github.com/ambientkit/plugin/middleware/etagcache"
	"github.com/ambientkit/plugin/middleware/gzipresponse"
	"github.com/ambientkit/plugin/middleware/healthcheck"
	"github.com/ambientkit/plugin/middleware/logrequest"
	"github.com/ambientkit/plugin/middleware/notrailingslash"
	"github.com/ambientkit/plugin/middleware/redirecttourl"
	"github.com/ambientkit/plugin/middleware/uptimerobotok"
	"github.com/ambientkit/plugin/router/awayrouter"
	"github.com/ambientkit/plugin/sessionmanager/scssession"
	"github.com/ambientkit/plugin/templateengine/htmlengine"
)

var (
	// StorageSitePath is the location of the site file.
	StorageSitePath = "storage/site.bin"
	// StorageSessionPath is the location of the session file.
	StorageSessionPath = "storage/session.bin"
)

// Plugins defines the plugins.
func Plugins() *ambient.PluginLoader {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	passwordHash := os.Getenv("AMB_PASSWORD_HASH")
	if len(passwordHash) == 0 {
		log.Fatalf("app: environment variable is missing: %v\n", "AMB_PASSWORD_HASH")
	}

	// Define the session manager so it can be used as a core plugin and
	// middleware.
	sessionManager := scssession.New(secretKey)

	return &ambient.PluginLoader{
		// Core plugins are implicitly trusted.
		Router:         awayrouter.New(nil),
		TemplateEngine: htmlengine.New(),
		SessionManager: sessionManager,
		// Trusted plugins are those that are typically needed to boot so they
		// will be enabled and given full access.
		TrustedPlugins: map[string]bool{
			"pluginmanager": true, // Page to manage plugins.
			"bearblog":      true, // Bear Blog login page.
			"bearcss":       true, // Bear Blog styling.
		},
		Plugins: []ambient.Plugin{
			// Standard library plugins.
			charset.New(),              // Charset to the HTML head.
			bearblog.New(passwordHash), // Bear Blog functionality.
			bearcss.New(),              // Bear Blog styling.
			debugpprof.New(),           // Go pprof debug endpoints.
			viewport.New(),             // Viewport in the HTML head.
			author.New(),               // Author in the HTML head.
			description.New(),          // Description the HTML head.
			pluginmanager.New(),        // Page to manage plugins.
			prism.New(),                // Prism CSS for codeblocks.
			stackedit.New(),            // Stackedit for editing markdown.
			googleanalytics.New(),      // Google Analytics.
			disqus.New(),               // Disqus for comments for blog posts.
			robots.New(),               // Robots file.
			sitemap.New(),              // Sitemap generator.
			rssfeed.New(),              // RSS feed generator.
			styles.New(),               // Style editing page.
			envinfo.New(),              // Show environment variables on the server.

			// App plugins.
			hello.New(),
			navigation.New(),
		},
		Middleware: []ambient.MiddlewarePlugin{
			logrequest.New(),      // Log every request as INFO.
			sessionManager,        // Session manager middleware.
			gzipresponse.New(),    // Compress all HTTP responses.
			etagcache.New(),       // Cache using Etag.
			redirecttourl.New(),   // Redirect to production URL.
			cors.New(),            // Enable CORS on /api/.
			healthcheck.New(),     // Provide 200 on /api/healthcheck.
			uptimerobotok.New(),   // Provide 200 on HEAD /.
			notrailingslash.New(), // Redirect all requests with a trailing slash.
		},
	}
}
