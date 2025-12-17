package context

import "github.com/leandroatallah/firefly/internal/engine/core"

type ContextProvider interface {
	// Use any to prevent life cycle imports
	SetAppContext(appContext any)
	AppContext() *core.AppContext
}
