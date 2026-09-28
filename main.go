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
	"sen-de-yaz/utilities"
	"strconv"
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
	"github.com/Elagoht/collage/pkg/collage"
)

//go:embed all:templates
var templatesFS embed.FS

//go:embed all:assets
var assetsFS embed.FS

//go:embed all:static
var staticFS embed.FS

var cacheDir = ".cache"

func main() {
	if file, err := utilities.LoadEnvFile("."); err != nil {
		log.Fatalf("sen-de-yaz: %v", err)
	} else if file != "" {
		log.Printf("sen-de-yaz: loaded environment from %s", file)
	}

	database, service, err := openUserService("app.sqlite")
	if err != nil {
		log.Fatalf("sen-de-yaz: initialize users: %v", err)
	}
	defer database.Close()

	buildFlag := flag.Bool("collage-build", false, "render the app to static files instead of serving it")
	outFlag := flag.String("out", "dist", "output directory for -collage-build")
	cleanFlag := flag.Bool("clean", false, "remove -out's existing contents before building")
	portFlag := flag.Int("port", envInt("PORT", 3000), "port to listen on (env PORT)")
	flag.Parse()

	devMode := os.Getenv("COLLAGE_DEV") == "1"

	app, err := newApp(devMode, *portFlag, service, stories.NewService(database))
	if err != nil {
		log.Fatalf("sen-de-yaz: %v", err)
	}

	if args := flag.Args(); len(args) > 0 {
		code, err := collage.DispatchCommands(context.Background(), app, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sen-de-yaz: %v\n", err)
		}
		os.Exit(code)
	}

	if *buildFlag {
		if err := staticBuild(app, *outFlag, *cleanFlag); err != nil {
			log.Fatalf("sen-de-yaz: static build: %v", err)
		}
		return
	}

	if err := app.ListenAndServe(); err != nil {
		log.Fatalf("sen-de-yaz: %v", err)
	}
}

func openUserService(path string) (*sql.DB, *users.UserService, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return database, users.NewService(database), nil
}

func newApp(devMode bool, port int, userService *users.UserService, storyService *stories.StoryService) (*collage.App, error) {
	pluginConfig, err := collage.LoadPluginConfig("plugins-config.json")
	if err != nil {
		return nil, fmt.Errorf("plugin configuration: %w", err)
	}

	csrfKey := os.Getenv("COLLAGE_CSRF_KEY")
	if csrfKey == "" {
		// Stable development key prevents form tokens from breaking on restart.
		// Set COLLAGE_CSRF_KEY to a random secret outside local development.
		csrfKey = "sen-de-yaz-development-csrf-key-change-me"
	}

	baseURL := envString("BASE_URL", "http://localhost:3000")
	plugins := []collage.Plugin{
		favicon.New(favicon.Options{
			FS:              assetsFS,
			Source:          "assets/icon.png",
			SVG:             "assets/icon.svg",
			Name:            "Sen de Yaz",
			ShortName:       "Sen de Yaz",
			ThemeColor:      "#6558f5",
			BackgroundColor: "#f7f8fc",
		}),
		secure.New(secure.Options{}),
		ratelimit.New(ratelimit.Options{}),
		// Register compression outside body-rewriting middleware so those
		// plugins see and can update the uncompressed HTML response.
		compress.New(compress.Options{}),
		honeypot.New(honeypot.Options{
			Key: []byte(envString("COLLAGE_HONEYPOT_KEY", "sen-de-yaz-development-honeypot-key-change-me")),
			// Enforced from process start; without it a path is only checked
			// after a page carrying its form has been served. Every POST form
			// in the app renders {{honeypot}}, so the prefixes are safe.
			Protect: []string{"/login", "/register", "/profile", "/stories", "/logout"},
		}),
		flash.New(flash.Options{Key: []byte(envString("COLLAGE_FLASH_KEY", "sen-de-yaz-development-flash-key-change-me"))}),
		meta.New(meta.Options{
			SiteName:        "Sen de Yaz",
			BaseURL:         baseURL,
			DefaultImage:    "/assets/icon.png",
			DefaultImageAlt: "Sen de Yaz logosu",
			Locales:         map[string]string{"en": "tr"},
		}),
		jsonld.New(),
		robots.New(robots.Options{}),
		sitemap.New(sitemap.Options{
			BaseURL: baseURL,
			Exclude: []string{"login", "register", "profile", "story-create"},
		}),
		minimizer.New(),
		optiimage.New(),
		accesslog.New(accesslog.Options{}),
		prometheus.New(prometheus.Options{Token: os.Getenv("METRICS_TOKEN")}),
		otel.New(otel.Options{}),
		htmlcheck.New(htmlcheck.Options{}),
	}
	if devMode {
		plugins = append(plugins, devtoolbar.New())
	}

	app, err := collage.New(&collage.Config{
		DevMode: devMode,
		Server: collage.ServerConfig{
			Host: envString("HOST", "localhost"),
			Port: port,
		},
		Template: collage.TemplateConfig{
			FS:        templatesFS,
			Root:      "templates",
			Extension: ".html",
		},
		Cache: collage.CacheConfig{
			Enabled:    true,
			Type:       "disk",
			Dir:        cacheDir,
			DefaultTTL: 5 * time.Minute,
		},
		PluginConfig: pluginConfig,
		Plugins:      plugins,
		Security: collage.SecurityConfig{
			CSRFKey: []byte(csrfKey),
		},
	})
	if err != nil {
		return nil, err
	}

	if err := register(app, userService, storyService); err != nil {
		return nil, err
	}

	assets, err := staticFiles(devMode)
	if err != nil {
		return nil, err
	}
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

func staticFiles(devMode bool) (fs.FS, error) {
	if devMode {
		if root, err := os.OpenRoot("static"); err == nil {
			return root.FS(), nil
		}
	}

	return fs.Sub(staticFS, "static")
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}

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
