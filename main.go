package main

import (
	"jungle-royale/game"
	"jungle-royale/socket"
)

func main() {
	var myGame = game.NewGame()
	_ = myGame

	socket.InitSocket()
}
