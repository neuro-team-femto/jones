package front

import (
	"log"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/neuro-team-femto/jones/config"
	_ "github.com/neuro-team-femto/jones/helpers"
)

// API
func Build() {
	// mode configuration, false by default
	var developmentMode, cmdBuildMode bool
	if config.Mode == "DEV" {
		developmentMode = true
	}
	if config.Mode == "BUILD_FRONT" {
		cmdBuildMode = true
	}
	if !developmentMode && !cmdBuildMode {
		return
	}

	// build options for prod
	buildOptions := api.BuildOptions{
		EntryPoints:       []string{"front/js/main.js"},
		Bundle:            true,
		MinifyWhitespace:  !developmentMode,
		MinifyIdentifiers: !developmentMode,
		MinifySyntax:      !developmentMode,
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
		ctx, err := api.Context(buildOptions)
		if err != nil {
			log.Fatal(err)
		} else if watchErr := ctx.Watch(api.WatchOptions{}); watchErr != nil {
			log.Fatal(watchErr)
		}
	}

	logsPlugin := api.Plugin{
		Name: "Logger",
		Setup: func(build api.PluginBuild) {
			build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
				if len(result.Errors) > 0 {
					for _, msg := range result.Errors {
						log.Printf("[error] js build: %v", msg.Text)
					}
				} else {
					if len(result.Warnings) > 0 {
						log.Printf("[front] js build success with %d warnings\n", len(result.Warnings))
						for _, msg := range result.Warnings {
							log.Printf("[warning] js build: %v\n", msg.Text)
						}
					} else {
						log.Println("[front] js build success")
					}
				}
				return api.OnEndResult{}, nil
			})
		},
	}
	buildOptions.Plugins = []api.Plugin{logsPlugin}
	build := api.Build(buildOptions)

	if len(build.Errors) > 0 {
		log.Printf("[front][error] js build: %v\n", build.Errors[0].Text)
	}
}
