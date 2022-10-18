package app

import (
	"github.com/vpofe/which-fix-version/git"
	"github.com/vpofe/which-fix-version/tui"
)

type App struct {
	Model tui.Model
}

func NewApp(gc *git.GitConfig, searchMode tui.SearchMode) (app App) {
	app.Model = tui.InitialModel(gc, searchMode)
	return
}
