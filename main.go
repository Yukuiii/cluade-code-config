package main

import (
	"embed"
	"io/fs"

	"claude-code-config/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of app structure
	app := app.NewApp()

	// Create a sub filesystem for the frontend/dist directory
	distFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "claude-code-config",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: distFS,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}