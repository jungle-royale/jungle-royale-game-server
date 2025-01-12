package serverlog

import (
	"jungle-royale/cons"
	"log"
)

type GameManagerLogger struct {
	GameLoggers []*GameLogger
}

func NewGameManagerLogger() *GameManagerLogger {
	return &GameManagerLogger{
		GameLoggers: make([]*GameLogger, cons.MaxGameNum),
	}
}

func (gml *GameManagerLogger) Log(logString string) {
	log.Printf("(GameManager) %s", logString)
}

func (gml *GameManagerLogger) AddGameLogger(gameIdx int, gameName string) *GameLogger {
	gml.GameLoggers[gameIdx].gameName = gameName
	return gml.GameLoggers[gameIdx]
}

type GameLogger struct {
	gameIdx  int
	gameName string
}

func NewGameLogger(gameIdx int, gameName string) *GameLogger {
	return &GameLogger{
		gameIdx:  gameIdx,
		gameName: gameName,
	}
}

func (gl *GameLogger) Log(logString string) {
	log.Printf("(Game[%d] \"%s\") %s", gl.gameIdx, gl.gameName, logString)
}

func (gl *GameLogger) GetGameName() string {
	return gl.gameName
}
