package bootstrap

import (
	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
)


func CreateManager() *core.ContextManager {
	return localstorage.CreateManager()
}
 
