package app

import "github.com/vpofe/which-fix-version/tui"

type App struct {
	Model tui.Model
}

func NewApp() (app App) {
	app.Model = tui.InitialModel()
	return
}
