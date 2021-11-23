package front

import (
	"log"
	"os"

	_ "github.com/creamlab/revcor/helpers"
	"github.com/evanw/esbuild/pkg/api"
)

var (
	developmentMode bool = false
	cmdBuildMode    bool = false
)

func init() {
	if os.Getenv("APP_ENV") == "DEV" {
		developmentMode = true
	}
	if os.Getenv("APP_ENV") == "BUILD_FRONT" {
		cmdBuildMode = true
	}
}

// API
func Build() {
	if !developmentMode && !cmdBuildMode {
		return
	}

	// build options for prod
	buildOptions := api.BuildOptions{
		EntryPoints:       []string{"front/js/main.js"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Engines: []api.Engine{
			{api.EngineChrome, "64"},
			{api.EngineFirefox, "53"},
			{api.EngineSafari, "11"},
			{api.EngineEdge, "79"},
		},
		Outdir: "public/scripts",
		Write:  true,
	}

	if developmentMode {
		buildOptions.MinifyWhitespace = false
		buildOptions.MinifyIdentifiers = false
		buildOptions.MinifySyntax = false
		buildOptions.Watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					for _, msg := range result.Errors {
						log.Printf("[error] js build: %v", msg.Text)
					}
				} else {
					if len(result.Warnings) > 0 {
						log.Printf("[info] js build success with %d warnings\n", len(result.Warnings))
						for _, msg := range result.Warnings {
							log.Printf("[warning] js build: %v\n", msg.Text)
						}
					} else {
						log.Println("[info] js build success")
					}
				}
			},
		}
	}

	build := api.Build(buildOptions)

	if len(build.Errors) > 0 {
		log.Printf("[error] js build: %v\n", build.Errors[0].Text)
	}
}
