package main

import (
	"runtime"

	"github.com/swordbeta/trello-burndown/src/backend"
	"github.com/swordbeta/trello-burndown/src/frontend"
	"github.com/swordbeta/trello-burndown/src/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	util.InitConfig()
	go frontend.Start()
	backend.Start()
}
