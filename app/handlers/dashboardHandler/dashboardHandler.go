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
		sysStat := make(map[string]interface{})
		for {
			memStats, err := utilities.Proc.Meminfo()
			if err != nil {
				panic(err)
			}
			cpuStat, err := utilities.Proc.CPUInfo()
			if err != nil {
				panic(err)
			}
			sysStat["freememory"] = memStats.MemFree
			sysStat["freeswap"] = memStats.SwapFree
			sysStat["totalmemory"] = memStats.MemTotal
			sysStat["totalswap"] = memStats.SwapTotal
			sysStat["cpucores"] = cpuStat[0].CPUCores

			data, err := json.Marshal(sysStat)
			if err != nil {
				panic(err)
			}
			err = c.WriteMessage(1, data)
			if err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))
}
