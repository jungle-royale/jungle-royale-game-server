package game

import (
	"fmt"
	"jungle-royale/cons"
	"jungle-royale/state"
	"net/http"
	"strconv"
)

type gameState struct {
	GameIdx    int
	GameName   string
	GameStatus string
	PlayerNum  int
	LastSec    int
}

func (gs *gameState) toString() string {
	return fmt.Sprintf("GameIdx: %d, GameName: %s, GameStatus: %s, PlayerNum: %d, LastSec: %d",
		gs.GameIdx, gs.GameName, gs.GameStatus, gs.PlayerNum, gs.LastSec)
}

func (gm *GameManager) SetServerManager() {

	// check all game state
	http.HandleFunc("/manage/check-all-game", func(w http.ResponseWriter, r *http.Request) {

		gm.gameManagerLogger.Log("monitoring: check all game state")

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, `{"status":405,"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		gameStateList := make([]*gameState, 0)
		for gameIdx := 0; gameIdx < cons.MaxGameNum; gameIdx++ {

			game := gm.gameRooms[gameIdx]

			var gameStatus string
			if game.playerNum != 0 {
				if game.state.GameState == state.Waiting || game.state.GameState == state.Counting {
					gameStatus = "Waiting"
				} else if game.state.GameState == state.Playing {
					gameStatus = "Playing"
				}
				gameStateList = append(gameStateList, &gameState{
					GameIdx:    gameIdx,
					GameName:   game.gameLogger.GetGameName(),
					GameStatus: gameStatus,
					PlayerNum:  game.state.Players.Length(),
					LastSec:    (game.state.LastGameTick*16)/1000 + cons.WAITING_MAP_CHUNK_NUM,
				})
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Playing Game List...")
		for _, gs := range gameStateList {
			res := gs.toString()
			fmt.Fprintln(w, res)
		}
	})

	// check game state
	http.HandleFunc("/manage/check-game", func(w http.ResponseWriter, r *http.Request) {

		gameIdxString := r.URL.Query().Get("gameIdx")
		if gameIdxString == "" {
			http.Error(w, "Missing gameIdx query parameter", http.StatusBadRequest)
			return
		}

		gm.gameManagerLogger.Log("monitoring: check game state " + gameIdxString)

		gameIdx, ok := strconv.Atoi(gameIdxString)
		if ok != nil {
			http.Error(w, "gameIdx parameter is invalid", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, `{"status":405,"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var gameState gameState
		game := gm.gameRooms[gameIdx]

		var gameStatus string
		if game.playerNum != 0 {
			if game.state.GameState == state.Waiting || game.state.GameState == state.Counting {
				gameStatus = "Waiting"
			} else if game.state.GameState == state.Playing {
				gameStatus = "Playing"
			}
			gameState.GameIdx = gameIdx
			gameState.GameName = game.gameLogger.GetGameName()
			gameState.GameStatus = gameStatus
			gameState.PlayerNum = game.state.Players.Length()
			gameState.LastSec = (game.state.LastGameTick*16)/1000 + cons.WAITING_MAP_CHUNK_NUM
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "Game %d state...\n", gameIdx)
			res := gameState.toString()
			fmt.Fprintln(w, res)
		}
	})

	// game force start
	http.HandleFunc("/manage/force-start", func(w http.ResponseWriter, r *http.Request) {

		gameIdxString := r.URL.Query().Get("gameIdx")
		if gameIdxString == "" {
			http.Error(w, "Missing gameIdx query parameter", http.StatusBadRequest)
			return
		}

		gameIdx, ok := strconv.Atoi(gameIdxString)
		if ok != nil {
			http.Error(w, "gameIdx parameter is invalid", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, `{"status":405,"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		game := gm.gameRooms[gameIdx]
		game.SetGameState(state.Counting)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Game %d force start\n", gameIdx)
	})
}
