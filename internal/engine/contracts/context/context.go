package context

import "github.com/boilerplate/ebiten-template/internal/engine/app"

type ContextProvider interface {
	// Use any to prevent life cycle imports
	SetAppContext(appContext any)
	AppContext() *app.AppContext
}
