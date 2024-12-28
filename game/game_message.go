package game

import (
	"jungle-royale/message"
	"jungle-royale/object"
	"log"

	"google.golang.org/protobuf/proto"
)

func (game *Game) OnMessage(data []byte, id string) {
	game.HandleMessage(id, data)
}

func (game *Game) HandleMessage(clientId string, data []byte) {
	var wrapper message.Wrapper
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		log.Printf("Failed to unmarshal message from client %s: %v", clientId, err)
		return
	}

	// dirChange message
	if dirChange := wrapper.GetChangeDir(); dirChange != nil {
		game.handleDirChange(clientId, dirChange)
	}

	if doDash := wrapper.GetDoDash(); doDash != nil {
		game.handleDoDash(clientId, doDash)
	}

	// bulletCreate message
	if createBullet := wrapper.GetCreateBullet(); createBullet != nil {
		game.handleBulletCreate(clientId, createBullet)
	}
}

func (game *Game) handleDirChange(clientId string, msg *message.ChangeDir) {
	if value, exists := game.state.MoverList.GetPlayers().Load(clientId); exists {
		player := value.(*object.Player)
		go player.DirChange(float64(msg.GetAngle()), msg.IsMoved)
	}
}

func (game *Game) handleBulletCreate(clientId string, msg *message.CreateBullet) {
	game.state.AddBullet(msg)
}

func (game *Game) handleDoDash(clientId string, msg *message.DoDash) {
	if value, exists := game.state.MoverList.GetPlayers().Load(clientId); exists {
		player := value.(*object.Player)
		go player.DoDash()
	}
}
