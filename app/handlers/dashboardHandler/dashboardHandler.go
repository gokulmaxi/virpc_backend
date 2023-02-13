package dashboardhandler

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Register(_route fiber.Router) {
	_route.Post("/memstat", func(c *fiber.Ctx) error {
		stats, err := utilities.Proc.Meminfo()
		if err != nil {
			panic(err)
		}
		fmt.Println(stats)
		data, err := json.Marshal(stats)
		if err != nil {
			panic(err)
		}
		return c.Send(data)
	})
	_route.Get("/mem", websocket.New(func(c *websocket.Conn) {
		for {
			stats, err := utilities.Proc.Meminfo()
			if err != nil {
				panic(err)
			}
			data, err := json.Marshal(stats)
			if err != nil {
				panic(err)
			}
			err = c.WriteMessage(1, data)
			if err != nil {
				break
			}
			fmt.Println(stats.MemFree)
			time.Sleep(2 * time.Second)
		}
	}))
}
