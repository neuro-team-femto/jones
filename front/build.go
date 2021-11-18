package front

import (
	"os"

	_ "github.com/creamlab/revcor/helpers"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

var (
	developmentMode bool = false
)

func init() {
	if os.Getenv("APP_ENV") == "DEV" {
		developmentMode = true
	}
}

// API

func Build() {
	// only build in certain conditions (= not when launching ducksoup in production)
	buildOptions := api.BuildOptions{
		EntryPoints:       []string{"front/src/main.js"},
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
		Outdir: "front/public/scripts",
		Write:  true,
	}
	if developmentMode {
		buildOptions.Watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					for _, msg := range result.Errors {
						log.Error().Msgf("[JS build] error: %v", msg.Text)
					}
				} else {
					if len(result.Warnings) > 0 {
						log.Info().Msgf("[JS build] success with %d warnings", len(result.Warnings))
						for _, msg := range result.Warnings {
							log.Info().Msgf("[JS build] warning: %v", msg.Text)
						}
					} else {
						log.Info().Msg("[JS build] success")
					}
				}
			},
		}
	}
	build := api.Build(buildOptions)

	if len(build.Errors) > 0 {
		log.Fatal().Msgf("JS build fatal error: %v", build.Errors[0].Text)
	}
}
