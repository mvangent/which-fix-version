package app

import (
	"github.com/vpofe/which-fix-version/git"
	"github.com/vpofe/which-fix-version/tui"
)

type App struct {
	Model tui.Model
}

func NewApp(gc *git.GitConfig) (app App) {
	app.Model = tui.InitialModel(gc)
	return
}
