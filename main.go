package main

import (
	"fmt"
	"jungle-royale/game"
	"os"
	"runtime"
	"strconv"
)

// intelliJ 연결 확인 - start, end

// 실행인자 처리
// 프로덕션과 테스트 실행 분리

// prod이면
// 아무것도 없으면 dev
// dev일 땐 minplayer랑 시간을 받을 수 잇음
// prod는 코드에 고정해놓기 - 설정값들!!

func main() {

	var cpu int = 2
	var debug bool = false
	var player = 10
	var time = 200

	configurationFromArgs(&cpu, &debug, &player, &time)

	// log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// log.SetOutput(io.Discard)

	runtime.GOMAXPROCS(cpu)

	gameManager := game.NewGameManager(debug)

	// dev 환경 실행
	if debug {
		go func() {
			gameManager.SetNewGame("test", player, time)
		}()
	}

	gameManager.Listen() // block
}

func configurationFromArgs(
	cpu *int,
	debug *bool,
	player *int,
	time *int,
) {
	// 실행 인자에서 첫 번째 항목 이후만 가져옴
	args := os.Args[1:]

	// 플래그 처리
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "debug", "d":
			*debug = true
		case "cpu":
			if i+1 < len(args) {
				value, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("Error: cpu 뒤에는 숫자가 와야 합니다.")
				} else {
					*cpu = value
					i++ // 숫자를 처리했으므로 i를 한 칸 더 증가
				}
			} else {
				fmt.Println("Error: cpu 뒤에 값이 없습니다.")
			}
		case "player":
			if i+1 < len(args) {
				value, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("Error: player 뒤에는 숫자가 와야 합니다.")
				} else {
					*player = value
					i++ // 숫자를 처리했으므로 i를 한 칸 더 증가
				}
			} else {
				fmt.Println("Error: player 뒤에 값이 없습니다.")
			}

		case "time":
			// -t 뒤에 숫자가 있는지 확인
			if i+1 < len(args) {
				value, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("Error: time 뒤에는 숫자가 와야 합니다.")
				} else {
					*time = value
					i++ // 숫자를 처리했으므로 i를 한 칸 더 증가
				}
			} else {
				fmt.Println("Error: time 뒤에 값이 없습니다.")
			}

		default:
			fmt.Printf("Unknown flag: %s\n", args[i])
		}
	}
}
