import message_pb2
import asyncio
import websockets
import requests
import random
import time
import httpx

websocket_array = []
# http_url = "http://localhost:8000"
http_url = "http://game-api.eternalsnowman.com:8080"
# ws_url = "ws://localhost:8000"
ws_url = "ws://game-api.eternalsnowman.com:8080"
min_player = 100
room_num = 1

async def create_room(room_num):
    create_room_url = http_url + "/api/create-game"
    data = {
        "roomId": "room" + str(room_num),
        "MinPlayers": min_player,
	    "MaxPlayTime": 60
    }
    async with httpx.AsyncClient() as client:
        response = await client.post(create_room_url, json=data)
        print("Created room:", response)
    

# roomId, clientId: int
async def connect_websocket(roomId, clientId):
    websocket_url = ws_url + f"/room?roomId=room{roomId}&clientId=r{roomId}c{clientId}&username={'이름'}"
    # websocket_url = ws_url + f"/room?roomId=0xoOrFVW3VYDuTUBEo_aaBpmQzJ9gwIgwiidILpzpt8&clientId=r{roomId}c{clientId}&username={'이름'}"
    # websocket_url = ws_url + f"/room?roomId=test&clientId=r{roomId}c{clientId}"
    websocket = await websockets.connect(websocket_url)

    asyncio.create_task(receive_data(websocket))

    websocket_array.append(websocket)
    print(f"Connected WebSocket {roomId}-{clientId}")
    
async def receive_data(websocket):
    """서버에서 데이터를 지속적으로 수신하여 버림."""
    try:
        async for message in websocket:
            pass  # 데이터 버림
    except Exception as e:
        print(f"Error receiving data: {e}")
    
    
async def do_somthing(websocket_idx):
    try:
        websocket = websocket_array[websocket_idx]
        wrapper = message_pb2.Wrapper()
        message_type = random.randint(0, 3)
        if message_type == 0:
            ca = message_pb2.ChangeAngle()
            ca.angle = random.uniform(0, 360)
            wrapper.changeAngle.CopyFrom(ca)

        elif message_type == 1:
            cd = message_pb2.ChangeDir()
            cd.angle = random.uniform(0, 360)
            cd.isMoved = True
            wrapper.changeDir.CopyFrom(cd)

        elif message_type == 2:
            dash = message_pb2.DoDash()
            dash.dash = True
            wrapper.doDash.CopyFrom(dash)

        elif message_type == 3:
            cbs = message_pb2.ChangeBulletState()
            cbs.isShooting = True
            wrapper.changeBulletState.CopyFrom(cbs)

        binary_data = wrapper.SerializeToString()
        await websocket.send(binary_data)
    except Exception as e:
        print(f"Error sending message to WebSocket {websocket_idx}: {e}")
    

async def close_websockets():
    close_tasks = [ws.close() for ws in websocket_array]
    await asyncio.gather(*close_tasks)
    print("all websocket connect close")
    print()
    

# 메인 함수
async def main():
    try:
        rn = room_num
        create_room_tasks = [create_room(i) for i in range(rn)]
        await asyncio.gather(*create_room_tasks)
        print(f"create {rn} rooms")
        print()
        
        mp = min_player
        connect_websocket_tasks = [connect_websocket(i, j) for i in range(rn) for j in range(mp)]
        await asyncio.gather(*connect_websocket_tasks)
        print(f"{rn * mp} users connect")
        print()
        
        semaphore = asyncio.Semaphore(1000)  # 동시에 최대 100개의 작업 실행
        
        for i in range(10000):
            
            async def limited_do_somthing(j):
                async with semaphore:
                    await do_somthing(j)
            
            do_somthing_tasks = [limited_do_somthing(j) for j in range(rn * mp)]
            await asyncio.gather(*do_somthing_tasks)
            print(f"do somthing {i} times")
            await asyncio.sleep(0.1)


    finally:
        await close_websockets()

# asyncio 실행
asyncio.run(main())