const faker = require("faker"); // 랜덤 값 생성용 (선택적으로 사용)
// import * as message from 'message_pb';
const message = require("./message_pb")

module.exports = {

  randomMessage: (context, events, done) => {

    // 랜덤 메시지 생성
    const messageType = Math.floor(Math.random() * 2);
    const wrapper = new message.Wrapper();
    switch (messageType) {
      case 0: // ChangeDir
        const change = new message.ChangeDir();
        change.setAngle(Math.random() * 360);
        change.setIsmoved((true));
        wrapper.setChangedir(change);
        break;
      case 1: // DoDash
        const dash = new message.DoDash();
        dash.setDash(true);
        wrapper.setDodash(dash);
        break;
      case 2: // CreateBullet
        const bulletCreate = new message.CreateBullet();
        // bulletCreate.setAngle(Math.random() * 360);
        bulletCreate.setAngle(0);
        wrapper.setCreatebullet(bulletCreate);
        break;
    }

    // Protobuf 메시지 직렬화
    const binaryData = wrapper.serializeBinary();

    // WebSocket으로 전송
    context.ws.send(binaryData);

    // console.log("Sent Protobuf message:", message);
    done();
  },
};