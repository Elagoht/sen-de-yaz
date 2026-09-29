package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sen-de-yaz/data/db"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utils"
	"strings"
	"time"

	accesslog "github.com/Elagoht/collage-accesslog"
	compress "github.com/Elagoht/collage-compress"
	devtoolbar "github.com/Elagoht/collage-devtoolbar"
	favicon "github.com/Elagoht/collage-favicon"
	flash "github.com/Elagoht/collage-flash"
	honeypot "github.com/Elagoht/collage-honeypot"
	htmlcheck "github.com/Elagoht/collage-htmlcheck"
	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	minimizer "github.com/Elagoht/collage-minimizer"
	optiimage "github.com/Elagoht/collage-opti-image"
	otel "github.com/Elagoht/collage-otel"
	prometheus "github.com/Elagoht/collage-prometheus"
	ratelimit "github.com/Elagoht/collage-ratelimit"
	robots "github.com/Elagoht/collage-robots"
	secure "github.com/Elagoht/collage-secure"
	sitemap "github.com/Elagoht/collage-sitemap"
	validate "github.com/Elagoht/collage-validate"
	"github.com/Elagoht/collage/pkg/collage"
)

//go:embed all:templates
var templatesFS embed.FS

//go:embed all:assets
var assetsFS embed.FS

//go:embed all:static
var staticFS embed.FS

// Disk cache directory
var cacheDir = ".cache"

// Runs the app; exits only after run's deferred cleanup is done
func main() {
	os.Exit(run())
}

// Loads env, initializes services and serves or builds the app, returns the exit code
func run() int {
	// Load ENV files. Do not accept not having one.
	if file, err := utils.LoadEnvFile("."); err != nil {
		log.Printf("sen-de-yaz: %v", err)
		return 1
	} else if file != "" {
		log.Printf("sen-de-yaz: loaded environment from %s", file)
	}

	// Initialize database and services
	database, usersService, storiesService, err := initializeDomains("app.sqlite")
	if err != nil {
		log.Printf("sen-de-yaz: initialize database: %v", err)
		return 1
	} // Don't forget to close it on shutdown
	defer database.Close()

	// Startup flags setup
	buildFlag := flag.Bool("collage-build", false, "render the app to static files instead of serving it")
	outFlag := flag.String("out", "dist", "output directory for -collage-build")
	cleanFlag := flag.Bool("clean", false, "remove -out's existing contents before building")
	portFlag := flag.Int("port", utils.EnvInt("PORT"), "port to listen on (env PORT)")
	flag.Parse()

	// Create an App as configured
	devMode := os.Getenv("COLLAGE_DEV") == "1"
	app, err := newApp(devMode, *portFlag, usersService, storiesService)
	if err != nil {
		log.Printf("sen-de-yaz: %v", err)
		return 1
	}

	// Dispatch flag commands
	if args := flag.Args(); len(args) > 0 {
		code, err := collage.DispatchCommands(context.Background(), app, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sen-de-yaz: %v\n", err)
		}
		return code
	}

	// Build App if needed
	if *buildFlag {
		if err := staticBuild(app, *outFlag, *cleanFlag); err != nil {
			log.Printf("sen-de-yaz: static build: %v", err)
			return 1
		}
		return 0
	}

	// If not closed yet, run the app
	if err := app.ListenAndServe(); err != nil {
		log.Printf("sen-de-yaz: %v", err)
		return 1
	}
	return 0
}

// Initialize database and services
func initializeDomains(path string) (
	*sql.DB,
	*users.UserService,
	*stories.StoryService,
	error,
) {
	database, err := db.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	// Session cookies are Secure when the site is served over HTTPS
	secureCookies := strings.HasPrefix(utils.EnvString("BASE_URL"), "https://")
	return database, users.NewService(database, secureCookies), stories.NewService(database), nil
}

// Create Collage App with project specific needs
func newApp(
	devMode bool,
	port int,
	userService *users.UserService,
	storyService *stories.StoryService,
) (*collage.App, error) {
	// Load plugin configs from typed json file.
	pluginConfig, err := collage.LoadPluginConfig("plugins-config.json")
	if err != nil {
		return nil, fmt.Errorf("plugin configuration: %w", err)
	}

	// CSRF key is required for security
	csrfKey := utils.EnvString("COLLAGE_CSRF_KEY")
	// Especially for SEO data
	baseURL := utils.EnvString("BASE_URL")
	// Configures plugins
	plugins := []collage.Plugin{
		// Creates required files from source png or svg files
		favicon.New(favicon.Options{
			FS:              assetsFS,
			Source:          "assets/icon.png",
			SVG:             "assets/icon.svg",
			Name:            "Sen de Yaz",
			ShortName:       "Sen de Yaz",
			ThemeColor:      "#6558f5",
			BackgroundColor: "#f7f8fc",
		}),
		// Defines security headers
		// Inline scripts run only with {{cspNonce}}; dev mode only reports
		secure.New(secure.Options{
			CSP: "default-src 'self'; script-src 'self' 'nonce-{nonce}'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'self'",
		}),
		// Creates rate limits
		ratelimit.New(ratelimit.Options{}),
		// Compresses the responses with brotli or gzip depending on Accept-Encoding
		compress.New(compress.Options{}),
		// Bot protection with honeypots placed with {{honeypot}}
		honeypot.New(honeypot.Options{
			Key: []byte(utils.EnvString("COLLAGE_HONEYPOT_KEY")),
		}),
		// Short lived cookies for HTTP response messages
		flash.New(flash.Options{Key: []byte(utils.EnvString("COLLAGE_FLASH_KEY"))}),
		// Registers form validation utility functions for templates
		validate.New(validate.Options{}),
		// Metadata options for pages
		meta.New(meta.Options{
			SiteName:        "Sen de Yaz",
			BaseURL:         baseURL,
			DefaultImage:    "/assets/icon.png",
			DefaultImageAlt: "Sen de Yaz logosu",
		}),
		// Structured Json-LD meta data for pages, a WebSite node on every page
		jsonld.NewWith(jsonld.Config{SiteName: "Sen de Yaz", SiteURL: baseURL}),
		// Robots.txt creation
		robots.New(robots.Options{}),
		// Sitemap creation. Panel pages need a session, a crawler only ever
		// gets the login redirect there, so only the public pages are listed.
		sitemap.New(sitemap.Options{
			BaseURL: baseURL,
			Exclude: []string{"home", "profile", "stories", "story-create", "story-detail"},
		}),
		// Minimize html, css, js, json responses
		minimizer.New(),
		// Optimize images by creating resized versions if width and height is set
		optiimage.New(),
		// Access log
		accesslog.New(accesslog.Options{}),
		// Prometheus
		prometheus.New(prometheus.Options{Token: os.Getenv("METRICS_TOKEN")}),
		// Open Telemetry
		otel.New(otel.Options{}),
		// Semantic HTML checks
		htmlcheck.New(htmlcheck.Options{}),
	}
	// Load devtoolbar only on dev mode
	if devMode {
		plugins = append(plugins, devtoolbar.New())
	}

	// Finally Create the Collage Application
	app, err := collage.New(&collage.Config{
		// Basic configurations
		DevMode: devMode,
		Server: collage.ServerConfig{
			Host: utils.EnvString("HOST"),
			Port: port,
		},
		// Register templates
		Template: collage.TemplateConfig{
			FS:        templatesFS,
			Root:      "templates",
			Extension: ".html",
		},
		// Default cache rules
		Cache: collage.CacheConfig{
			Enabled:    true,
			Type:       "disk",
			Dir:        cacheDir,
			DefaultTTL: 5 * time.Minute,
		},
		// Registers plugins and their configs
		PluginConfig: pluginConfig,
		Plugins:      plugins,
		// All unsupported page paths returns "TR" here
		Locale: collage.LocaleConfig{
			Default: "tr",
		},
		// Security configuration
		Security: collage.SecurityConfig{
			CSRFKey: []byte(csrfKey),
		},
	})
	if err != nil {
		return nil, err
	}

	// Register app's needs and routes
	if err := register(app, userService, storyService); err != nil {
		return nil, err
	}

	// Mount secured Static files to app
	assets, err := staticFiles(devMode)
	if err != nil {
		return nil, err
	}
	// Create the required directories
	if err := app.Mount("/static/", assets); err != nil {
		return nil, fmt.Errorf("mount static files: %w", err)
	}
	if err := os.MkdirAll("uploads/profile", 0o755); err != nil {
		return nil, fmt.Errorf("create uploads directory: %w", err)
	}
	if err := app.Mount("/uploads/", os.DirFS("uploads")); err != nil {
		return nil, fmt.Errorf("mount uploads: %w", err)
	}

	return app, nil
}

// Create a secure storage path that doesn't allow links to follow outer places
func staticFiles(devMode bool) (fs.FS, error) {
	if devMode {
		if root, err := os.OpenRoot("static"); err == nil {
			return root.FS(), nil
		}
	}

	return fs.Sub(staticFS, "static")
}

// Build the configured app
func staticBuild(app *collage.App, outDir string, clean bool) error {
	builder, err := collage.NewBuilder(app, collage.BuildOptions{
		OutDir: outDir,
		Clean:  clean,
	})
	if err != nil {
		return err
	}

	report, buildErr := builder.Build(context.Background())
	collage.PrintBuildReport(os.Stdout, report, buildErr)

	return buildErr
}
