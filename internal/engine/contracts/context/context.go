package context

import "github.com/boilerplate/ebiten-template/internal/engine/app"

// ContextProvider gives engine components access to the shared application context.
type ContextProvider interface {
	// Use any to prevent life cycle imports

	// SetAppContext injects the shared application context.
	SetAppContext(appContext any)
	// AppContext returns the typed application context.
	AppContext() *app.AppContext
}
