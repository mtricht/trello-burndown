package main

import (
	"runtime"

	"github.com/swordbeta/trello-burndown/src/util"
	"github.com/swordbeta/trello-burndown/src/watcher"
	"github.com/swordbeta/trello-burndown/src/webserver"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	util.InitConfig()
	go webserver.Start()
	watcher.Start()
}
