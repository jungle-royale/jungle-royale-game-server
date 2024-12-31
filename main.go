package main

import (
	"jungle-royale/game"
	"jungle-royale/network"
	"log"
	"runtime"
	"time"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// log.SetOutput(io.Discard)

	runtime.GOMAXPROCS(2)

	roomManager := network.NewRoomManager()

	var socket network.Socket = roomManager

	go func() {
		time.Sleep(1000 * time.Millisecond)                                                      // 3초
		var testGame network.Room = game.NewGame(&socket, 100, 200).SetReadyStatus().StartGame() // 플레이어 수, 게임 시간
		roomManager.RegisterRoom(network.RoomId("test"), &testGame)
	}()

	roomManager.Listen()
}
