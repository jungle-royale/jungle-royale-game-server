import { Scene } from 'phaser';
import * as message from '../../message_pb';

class Player {
    constructor(id, x, y, scene) {
        this.id = id; // 플레이어 ID
        this.x = x; // X 좌표
        this.y = y; // Y 좌표
        this.scene = scene; // Phaser Scene
        this.circle = null; // Phaser Graphics 객체
        this.last_dir = null;
        this.createCircle();
    }

    // 원 생성
    createCircle() {
        this.circle = this.scene.add.graphics();
        this.circle.fillStyle(this.scene.myPlayerId === this.id ? 0xff0000 : 0x00ff00, 1); // 내 플레이어는 빨간색, 나머지는 녹색
        this.circle.fillCircle(this.x, this.y, 0.5); // x, y 좌표와 반지름 10
    }

    // 원 이동
    updatePosition(x, y) {
        this.x = x;
        this.y = y;
        if (this.circle) {
            this.circle.clear(); // 이전 원 제거
            this.circle.fillStyle(this.scene.myPlayerId === this.id ? 0xff0000 : 0x00ff00, 1);
            this.circle.fillCircle(this.x, this.y, 0.5);
        }
    }

    // 원 제거
    destroy() {
        if (this.circle) {
            this.circle.destroy();
            this.circle = null;
        }
    }
}

class Bullet {
    constructor(bulletId, x, y, scene) {
        this.bulletId = bulletId; // 플레이어 ID
        this.x = x; // X 좌표
        this.y = y; // Y 좌표
        this.scene = scene; // Phaser Scene
        this.circle = null; // Phaser Graphics 객체
        this.createCircle();
    }

    // 원 생성
    createCircle() {
        this.circle = this.scene.add.graphics();
        this.circle.fillStyle(0x000000, 1);
        this.circle.fillCircle(this.x, this.y, 0.2); // x, y 좌표와 반지름 10
    }

    // 원 이동
    updatePosition(x, y) {
        this.x = x;
        this.y = y;
        if (this.circle) {
            this.circle.clear(); // 이전 원 제거
            this.circle.fillStyle(0x000000, 1);
            this.circle.fillCircle(this.x, this.y, 0.2);
        }
    }

    // 원 제거
    destroy() {
        if (this.circle) {
            this.circle.destroy();
            this.circle = null;
        }
    }
}

function sendDashMessage(socket) {
    const dash = new message.DoDash();
    dash.setDash(true);
    const wrapper = new message.Wrapper();
    wrapper.setDodash(dash);
    const binaryData = wrapper.serializeBinary();
    socket.send(binaryData);
}

function sendChangeMessage(socket, angle, isMoved) {

    // console.log("send change ", angle);

    if (socket.readyState !== WebSocket.OPEN) {
        console.error('WebSocket is not open.');
        return;
    }

    // Change 메시지 생성
    const change = new message.ChangeDir();
    change.setAngle(angle);
    change.setIsmoved(isMoved);

    // Wrapper 메시지 생성
    const wrapper = new message.Wrapper();
    wrapper.setChangedir(change);

    // 직렬화하여 바이너리 데이터로 변환
    const binaryData = wrapper.serializeBinary();

    // WebSocket으로 전송
    socket.send(binaryData);
}

function sendBulletCreateMessage(socket, playerId, x, y, angle) {

    if (socket.readyState !== WebSocket.OPEN) {
        console.error('WebSocket is not open.');
        return;
    }

    const bulletCreate = new message.CreateBullet();
    bulletCreate.setPlayerid(playerId);
    bulletCreate.setStartx(x);
    bulletCreate.setStarty(y);
    bulletCreate.setAngle(angle);

    // Wrapper 메시지 생성
    const wrapper = new message.Wrapper();
    wrapper.setCreatebullet(bulletCreate);

    // 직렬화하여 바이너리 데이터로 변환
    const binaryData = wrapper.serializeBinary();

    // WebSocket으로 전송
    socket.send(binaryData);
}

export class Game extends Scene {

    constructor() {
        super('Game');
        this.ID = "";
        this.socket = null; // WebSocket 연결
        this.players = {}; // 플레이어 데이터 저장
        this.bullets = {};
        this.last_x = 0;
        this.last_y = 0;
        this.dash = false;
    }

    create() {

        this.createGrid(1000, 1000, 20, 0x000000);

        let testRoomId = "test"

        // WebSocket 연결
        this.socket = new WebSocket(`ws://localhost:8000/room?roomId=${testRoomId}`)

        console.log(`ws://localhost:8000/ws?roomId=${testRoomId}`)

        // WebSocket 이벤트 처리
        this.socket.onopen = () => {
            console.log('Connected to WebSocket server');
        };

        this.socket.onmessage = async (event) => {
            try {
                let data;

                // WebSocket에서 받은 데이터 처리
                if (event.data instanceof Blob) {
                    // Blob 데이터를 ArrayBuffer로 변환
                    const arrayBuffer = await event.data.arrayBuffer();
                    data = new Uint8Array(arrayBuffer);
                } else if (event.data instanceof ArrayBuffer) {
                    // ArrayBuffer로 바로 처리
                    data = new Uint8Array(event.data);
                } else {
                    console.error("Unsupported WebSocket message format:", typeof event.data);
                    return;
                }

                // 디코드된 데이터 출력
                // console.log("Received raw data:", data);

                // Protobuf 메시지 디코드
                const wrapper = message.Wrapper.deserializeBinary(data);
                // console.log("Decoded Wrapper message:", wrapper.toObject());

                // Wrapper 내부 메시지 처리
                if (wrapper.hasGameinit()) {
                    const gameInit = wrapper.getGameinit();
                    this.handleGameInit(gameInit.toObject());
                } else if (wrapper.hasGamestate()) {
                    const gameState = wrapper.getGamestate();
                    this.handleGameState(gameState.toObject());
                } else if (wrapper.hasGamecount()) {
                    const gameCount = wrapper.getGamecount();
                    console.log(gameCount)
                } else {
                    console.error("Unknown message type received.");
                }
            } catch (error) {
                console.error('Failed to process WebSocket message:', error);
            }
        };


        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.socket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        // 키보드 입력 설정
        this.cursors = this.input.keyboard.createCursorKeys();

        // **화면 클릭 이벤트 등록**
        this.input.on('pointerdown', (pointer) => {
            const centerX = this.cameras.main.centerX;
            const centerY = this.cameras.main.centerY;
            const worldPoint = this.cameras.main.getWorldPoint(centerX, centerY);
            const dx = pointer.x - centerX;
            const dy = pointer.y - centerY;
            const radians = Math.atan2(dy, dx);
            const degrees = radians * (180 / Math.PI) + 90;
            sendBulletCreateMessage(
                this.socket,
                this.ID,
                worldPoint.x,
                worldPoint.y,
                degrees,
            );
        });

        this.cameras.main.setZoom(15)

        this.spaceKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.SPACE);
        this.keys = this.input.keyboard.addKeys({
            up: Phaser.Input.Keyboard.KeyCodes.W,
            down: Phaser.Input.Keyboard.KeyCodes.S,
            left: Phaser.Input.Keyboard.KeyCodes.A,
            right: Phaser.Input.Keyboard.KeyCodes.D
        });
    }

    update() {
        let x = 0;
        let y = 0;
        if (this.cursors.left.isDown || this.keys.left.isDown) {
            x -= 5.0;
        }
        if (this.cursors.right.isDown || this.keys.right.isDown) {
            x += 5.0;
        }
        if (this.cursors.up.isDown || this.keys.up.isDown) {
            y -= 5.0;
        }
        if (this.cursors.down.isDown || this.keys.down.isDown) {
            y += 5.0;
        }
        if (this.last_x != x || this.last_y != y) {
            this.last_x = x;
            this.last_y = y;

            if (x == 0 && y == 0)
                sendChangeMessage(this.socket, 0, false);
            else if (x == 0 && y == -5)
                sendChangeMessage(this.socket, 0, true);
            else if (x == 5 && y == -5)
                sendChangeMessage(this.socket, 45, true);
            else if (x == 5 && y == 0)
                sendChangeMessage(this.socket, 90, true);
            else if (x == 5 && y == 5)
                sendChangeMessage(this.socket, 135, true);
            else if (x == 0 && y == 5)
                sendChangeMessage(this.socket, 180, true);
            else if (x == -5 && y == 5)
                sendChangeMessage(this.socket, 225, true);
            else if (x == -5 && y == 0)
                sendChangeMessage(this.socket, 270, true);
            else if (x == -5 && y == -5)
                sendChangeMessage(this.socket, 315, true);
        }
        if (this.spaceKey.isDown && this.dash == false) {
            this.dash = true
            sendDashMessage(this.socket);
        }
        else if(this.dash == true) {
            this.dash = false;
        }
    }

    createGrid(width, height, gridSize, color) {
        const graphics = this.add.graphics();
        graphics.lineStyle(1, color, 1); // 선 스타일: 두께, 색상, 투명도

        // 세로선 그리기
        for (let x = 0; x <= width; x += gridSize) {
            graphics.beginPath();
            graphics.moveTo(x, 0);
            graphics.lineTo(x, height);
            graphics.strokePath();
        }

        // 가로선 그리기
        for (let y = 0; y <= height; y += gridSize) {
            graphics.beginPath();
            graphics.moveTo(0, y);
            graphics.lineTo(width, y);
            graphics.strokePath();
        }
    }

    handleGameInit(init) {
        console.log("ID ", init.id);
        if (init.id) {
            this.ID = init.id;
        }
    }

    handleGameState(state) {
        // console.log(state);
        if (state.playerstateList) {

            for (const player in this.playerstateList) {
                this.players[player].circle.clear();
            }

            // 새 playersList를 순회하면서 업데이트
            const playerIdsInState = new Set(); // state.playersList에서 확인된 player id를 저장할 Set
            state.playerstateList.forEach((player) => {
                playerIdsInState.add(player.id); // state에 있는 player id 저장
                if (!(player.id in this.players)) {
                    this.players[player.id] = new Player(player.id, player.x, player.y, this);
                }
                if (this.players[player.id]) {
                    this.players[player.id].updatePosition(player.x, player.y);
                }
                if (player.id == this.ID) {
                    this.cameras.main.centerOn(player.x, player.y);
                }
            });
    
            // this.players에 있지만 state.playersList에 없는 player 삭제
            for (const playerId in this.players) {
                if (!playerIdsInState.has(playerId)) {
                    delete this.players[playerId]; // this.players에서 제거
                }
            }
        }
    
        if (state.bulletstateList) {

            // 기존 bullets 처리
            for (const bullet in this.bullets) {
                this.bullets[bullet].circle.clear();
            }
    
            // 새 bulletstateList를 순회하면서 업데이트
            const bulletIdsInState = new Set(); // state.bulletstateList에서 확인된 bulletid를 저장할 Set
            state.bulletstateList.forEach((bullet) => {
                bulletIdsInState.add(bullet.bulletid); // state에 있는 bulletid 저장
                if (!(bullet.bulletid in this.bullets)) {
                    this.bullets[bullet.bulletid] = new Bullet(bullet.bulletid, bullet.x, bullet.y, this);
                }
                if (this.bullets[bullet.bulletid]) {
                    this.bullets[bullet.bulletid].updatePosition(bullet.x, bullet.y);
                }
            });
    
            // this.bullets에 있지만 state.bulletstateList에 없는 bullet 삭제
            for (const bulletId in this.bullets) {
                if (!bulletIdsInState.has(bulletId)) {
                    this.bullets[bulletId].circle.clear(); // 해당 bullet의 리소스 정리
                    delete this.bullets[bulletId]; // this.bullets에서 제거
                }
            }
        }
    }
    
}