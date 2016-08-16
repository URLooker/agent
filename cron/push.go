package cron

import (
	"log"
	"time"

	"github.com/urlooker/web/api"
	webg "github.com/urlooker/web/g"

	"github.com/urlooker/agent/backend"
	"github.com/urlooker/agent/g"
)

func Push() {
	hostname, err := g.Hostname()
	if err != nil {
		log.Println("get hostname err:", err)
		hostname = "null"
	}
	for {
		checkResults := make([]*webg.CheckResult, 0)
		itemResults := g.CheckResultQueue.PopBack(500)
		if len(itemResults) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for _, itemResult := range itemResults {
			checkResult := itemResult.(*webg.CheckResult)
			checkResults = append(checkResults, checkResult)
		}

		var resp string
		sendResultReq := api.SendResultReq{
			Hostname:     hostname,
			CheckResults: checkResults,
		}
		err := backend.CallRpc("Web.SendResult", sendResultReq, &resp)
		if err != nil {
			log.Println("error:", err)
		}

		if g.Config.Debug {
			log.Println("<=", resp)
		}
	}
	return
}
