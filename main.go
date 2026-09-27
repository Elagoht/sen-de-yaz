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
	"strconv"
	"time"

	favicon "github.com/Elagoht/collage-favicon"
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
	database, service, err := openUserService("data/db/app.sqlite")
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
		Plugins: []collage.Plugin{
			favicon.New(favicon.Options{
				FS:              assetsFS,
				Source:          "assets/icon.png",
				SVG:             "assets/icon.svg",
				Name:            "Sen de Yaz",
				ShortName:       "Sen de Yaz",
				ThemeColor:      "#6558f5",
				BackgroundColor: "#f7f8fc",
			}),
		},
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
