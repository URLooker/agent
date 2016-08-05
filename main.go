package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urlooker/agent/backend"
	"github.com/urlooker/agent/cron"
	"github.com/urlooker/agent/g"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	backend.InitClients(g.Config.Web.Addrs)

	g.Init()
}

func main() {
	go cron.Push()
	cron.StartCheck()
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
