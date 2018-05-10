package main

import (
	"github.com/swordbeta/trello-burndown/src/util"
	"github.com/swordbeta/trello-burndown/src/watcher"
	"github.com/swordbeta/trello-burndown/src/webserver"
)

func main() {
	util.InitConfig()
	go webserver.Start()
	watcher.Start()
}
