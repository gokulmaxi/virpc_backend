package utilities

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store = InitSessison()

func InitSessison() *session.Store {
	_store := session.New()
	return _store
}
