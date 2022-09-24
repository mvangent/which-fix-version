package app

import "github.com/vpofe/just-in-time/tui"

type App struct {
	Model tui.Model
}

func NewApp() (app App) {
	app.Model = tui.InitialModel()
	return
}
