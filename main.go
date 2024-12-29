package main

import (
	"jungle-royale/game"
	"jungle-royale/network"
	"time"
)

func main() {

	roomManager := network.NewRoomManager()

	var socket network.Socket = roomManager

	go func() {
		time.Sleep(1000 * time.Millisecond) // 3초
		var testGame network.Room = game.NewGame(&socket, 2).SetReadyStatus().StartGame()
		roomManager.RegisterRoom(network.RoomId("test"), &testGame)
	}()

	roomManager.Listen()
}
