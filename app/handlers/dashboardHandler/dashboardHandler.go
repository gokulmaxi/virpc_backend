package dashboardhandler

import (
	"app/utilities"
	"encoding/json"
	"fmt"
	"math"
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
	_route.Get("/sysinfo", websocket.New(func(c *websocket.Conn) {
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
			StorageStat, err := utilities.GetStorage()
			if err != nil {
				panic(err)
			}
			// NOTE meminfo returns memoru in kb
			sysStat["freememory"] = float64(*memStats.MemAvailable) / math.Pow(1000, 2)
			sysStat["freeswap"] = float64(*memStats.SwapFree) / math.Pow(1000, 2)
			sysStat["totalmemory"] = float64(*memStats.MemTotal) / math.Pow(1000, 2)
			sysStat["totalswap"] = float64(*memStats.SwapTotal) / math.Pow(1000, 2)
			sysStat["cpucores"] = cpuStat[0].CPUCores
			sysStat["totalstorage"] = StorageStat.Totalsize
			sysStat["freestorage"] = StorageStat.Available
			data, err := json.Marshal(sysStat)
			if err != nil {
				panic(err)
			}
			err = c.WriteMessage(1, data)
			if err != nil {
				break
			}
			time.Sleep(35 * time.Second)
		}
	}))
}
