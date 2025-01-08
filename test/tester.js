// import { Wrapper, ChangeAngle, ChangeDir, DoDash, ChangeBulletState } from "./message_pb.js";
const message = require("./message_pb");


module.exports = {

  randomMessage: (context, events, done) => {

// 랜덤 메시지 생성
const messageType = Math.floor(Math.random() * 4);
const wrapper = new message.Wrapper();
switch (messageType) {
  case 0:  // ChangeAngle
    const ca = new message.ChangeAngle();
    ca.setAngle(Math.random() * 360);
    wrapper.setChangeangle(ca);
    break;
  case 1:  // ChangeDir
    const cd = new message.ChangeDir();
    cd.setAngle(Math.random() * 360);
    cd.setIsmoved(true);
    wrapper.setChangedir(cd);
    break;
  case 2:  // DoDash
    const dash = new message.DoDash();
    dash.setDash(true);
    wrapper.setDodash(dash);
    break;
  case 3:  // ChangeBulletState
    const cbs = new message.ChangeBulletState();
    cbs.setIsshooting(true)
    wrapper.setChangebulletstate(cbs);
    break;
}

// Protobuf 메시지 직렬화
const binaryData = wrapper.serializeBinary();

// WebSocket으로 전송
context.ws.send(binaryData);

// console.log("Sent Protobuf message:", message);
done();
  },

  upMessage: (context, events, done) => { 
    const change = new message.ChangeDir();
    const wrapper = new message.Wrapper();
    change.setAngle(0)
    change.setIsmoved(true)
    wrapper.setChangedir(change)
    const binaryData = wrapper.serializeBinary();
    context.ws.send(binaryData);
    done();
  },

  leftMessage: (context, events, done) => { 
    const change = new message.ChangeDir();
    const wrapper = new message.Wrapper();
    change.setAngle(270)
    change.setIsmoved(true)
    wrapper.setChangedir(change)
    const binaryData = wrapper.serializeBinary();
    context.ws.send(binaryData);
    done();
  },

  rightMessage: (context, events, done) => { 
    const change = new message.ChangeDir();
    const wrapper = new message.Wrapper();
    change.setAngle(90)
    change.setIsmoved(true)
    wrapper.setChangedir(change)
    const binaryData = wrapper.serializeBinary();
    context.ws.send(binaryData);
    done();
  },

  downMessage: (context, events, done) => { 
    const change = new message.ChangeDir();
    const wrapper = new message.Wrapper();
    change.setAngle(180)
    change.setIsmoved(true)
    wrapper.setChangedir(change)
    const binaryData = wrapper.serializeBinary();
    context.ws.send(binaryData);
    done();
  }
};