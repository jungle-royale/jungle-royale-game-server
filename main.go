package main

import (
	"jungle-royale/game"
	"log"
	"runtime"
	"time"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// log.SetOutput(io.Discard)

	runtime.GOMAXPROCS(2)

	gameManager := game.NewGameManager()

	go func() {
		time.Sleep(1000 * time.Millisecond) // 3초
		gameManager.CreateRoom("test", 4, 200)
	}()

	gameManager.Listen()
}
