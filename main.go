package main

import (
	"jungle-royale/game"
	"log"
	"runtime"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// log.SetOutput(io.Discard)

	runtime.GOMAXPROCS(2)

	gameManager := game.NewgameManager()

	gameManager.Listen()
}
