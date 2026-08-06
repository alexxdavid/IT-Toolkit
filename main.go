package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"ittoolkit/internal/platform"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if !platform.WebView2Installed() {
		platform.ShowMessageBox(
			"Solutions IT Toolkit",
			"Microsoft Edge WebView2 Runtime is required to run Solutions IT Toolkit but was not found.\n\n"+
				"The app will now open the WebView2 installer download page.\n"+
				"Install it and launch Solutions IT Toolkit again.",
		)
		_ = platform.OpenURL("https://go.microsoft.com/fwlink/p/?LinkId=2124703")
		return
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Solutions IT Toolkit",
		Width:     1280,
		Height:    820,
		MinWidth:  960,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 17, G: 24, B: 39, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
