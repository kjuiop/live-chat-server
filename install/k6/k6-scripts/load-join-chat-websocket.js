import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import ws from 'k6/ws';

export const options = {
    stages: [
        { duration: '10s', target: 2000 },
        { duration: '30s', target: 2000 },
        { duration: '30s', target: 1000 },
        { duration: '5s', target: 0 },
    ],
};

// 5000 명일 때 30개 1분 50명 (1사람당 송신)
export default function () {
    const url = `wss://domain/ws/room_id/jake_user_${__VU}/nickname_${__VU}?joinMessageCount=30`;

    const res = ws.connect(url, {}, (socket) => {
        // WebSocket 연결 성공 시 이벤트 설정
        socket.on('open', () => {
            console.log(`WebSocket connected ${__VU}`)

            // 5000 명 증가하는데
            // broadcast 도 보내는것 (몇개의 메시지
            if (__VU <= 30) {

                sleep(12)
                // 메시지 보내기 로직 (1명당 1분에 10개의 메시지 보내기)
                // 2초에 1개씩 1분에 30개 * 50 = 1500 개 송신하는거고
                // 남은 5000명이 저걸 받아줘야되는데

                for (let i = 0; i < 20; i++) {
                    const message = {
                        Method: "chat",
                        Params: {
                            UserId: `user-${__VU}`, // 고유한 UserId 사용
                            Message: `Message number ${i + 1} from user-${__VU}`,
                            Nickname: `User-${__VU}`,
                            PhotoUrl: "www.foo.com/a.jpg",
                            Type: "user defined string",
                            UserData: "notice"
                        }
                    };

                    socket.send(JSON.stringify(message));
                    // sleep(randomIntBetween(2, 6));
                    sleep(6);
                }
            }
        });

        socket.on('message', (message) => {
            // console.log(`WebSocket received ${__VU} message`)
        });

        socket.on('close', () => {
            console.log('WebSocket disconnected');
        });

        socket.on('error', (e) => {
            console.log('An unexpected error occured: ', e.error());
        });
    });

    check(res, {
        'WebSocket connection established': (r) => r && r.status === 101,
    });
}