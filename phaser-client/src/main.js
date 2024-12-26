import { Game } from './scenes/Game';

const config = {
    type: Phaser.AUTO,
    width: 1000,
    height: 1000,
    physics: {
        default: 'arcade', // 기본 물리 엔진 설정
        // arcade: {
        //     gravity: { y: 300 }, // 중력 설정
        //     debug: false        // 디버그 모드 비활성화
        // }
    },
    parent: 'game-container',
    backgroundColor: '0xAAAAAA',
    scale: {
        mode: Phaser.Scale.FIT,
        autoCenter: Phaser.Scale.CENTER_BOTH
    },
    scene: [
        Game,
    ]
};

export default new Phaser.Game(config);