import { check, sleep } from 'k6';
import ws from 'k6/ws';

export const options = {
    stages: [
        { duration: '5s', target: 500 },
        { duration: '1s', target: 0 },
    ],
};

export default function () {
    const url = 'wss://domain/ws/N2_01JGDPNW4DYYZCSZ87P46DBEF9/jungin.kim@catenoid.net%2FAdmin?joinMessageCount=30';

    const res = ws.connect(url, {}, (socket) => {
        // WebSocket 연결 성공 시 이벤트 설정
        socket.on('open', () => {
            console.log('WebSocket connected');

            // 5000 명 증가하는데
            // broadcast 도 보내는것 (몇개의 메시지

            // 10명
            // if (__VU <= 10) {
            //     // 메시지 보내기 로직 (1명당 1분에 10개의 메시지 보내기)
            //     for (let i = 0; i < 10; i++) {
            //         const message = {
            //             Method: "chat",
            //             Params: {
            //                 UserId: `user-${__VU}`, // 고유한 UserId 사용
            //                 Message: `Message number ${i + 1} from user-${__VU}`,
            //                 Nickname: `User-${__VU}`,
            //                 PhotoUrl: "www.foo.com/a.jpg",
            //                 Type: "user defined string",
            //                 UserData: "notice"
            //             }
            //         };
            //         socket.send(JSON.stringify(message));
            //         sleep(6); // 6초 간격으로 메시지를 보냄 (1분에 10개)
            //     }
            // }
        });

        socket.on('message', (message) => {
            console.log(`Received message: ${message}`);
        });

        socket.on('close', () => {
            console.log('WebSocket disconnected');
        });

        socket.on('error', (e) => {
            console.error(`WebSocket error: ${e}`);
        });
    });

    check(res, {
        'WebSocket connection established': (r) => r && r.status === 101,
    });
}