const protobuf = require("protobufjs");
const faker = require("faker"); // 랜덤 값 생성용 (선택적으로 사용)
let root = null;

// Protobuf 정의 파일 로드
protobuf.load("message.proto", (err, loadedRoot) => {
  if (err) throw err;
  root = loadedRoot;
});

module.exports = {
  randomMessage: (context, events, done) => {
    if (!root) {
      console.error("Protobuf not loaded yet");
      return done();
    }

    // 메시지 타입 정의
    const Wrapper = root.lookupType("Wrapper");
    const ChangeDir = root.lookupType("ChangeDir");
    const DoDash = root.lookupType("DoDash");
    const CreateBullet = root.lookupType("CreateBullet");

    // 랜덤 메시지 생성
    const messageType = Math.floor(Math.random() * 3);
    let message;
    switch (messageType) {
      case 0: // ChangeDir
        message = Wrapper.create({
          MessageType: {
            changeDir: ChangeDir.create({
              angle: Math.random() * 360,
              isMoved: Math.random() > 0.5,
            }),
          },
        });
        break;
      case 1: // DoDash
        message = Wrapper.create({
          MessageType: {
            doDash: DoDash.create({
              dash: Math.random() > 0.5,
            }),
          },
        });
        break;
      case 2: // CreateBullet
        message = Wrapper.create({
          MessageType: {
            createBullet: CreateBullet.create({
              playerId: faker.datatype.uuid(), // 랜덤 Player ID
              startX: Math.random() * 1000,   // 랜덤 X 좌표
              startY: Math.random() * 1000,   // 랜덤 Y 좌표
              angle: Math.random() * 360,     // 랜덤 각도
            }),
          },
        });
        break;
    }

    // Protobuf 메시지 직렬화
    const buffer = Wrapper.encode(message).finish();

    // WebSocket으로 전송
    context.ws.send(buffer);

    // console.log("Sent Protobuf message:", message);
    done();
  },
};
