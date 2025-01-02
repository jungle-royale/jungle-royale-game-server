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

	gameManager := game.NewGameManager()

	// go func() {
	// 	time.Sleep(1000 * time.Millisecond) // 3초
	// 	gameManager.CreateGame("test", 10, 200)
	// }()

	gameManager.Listen()
}
