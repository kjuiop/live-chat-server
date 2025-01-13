import { check, sleep } from 'k6';
import ws from 'k6/ws';

export const options = {
    stages: [
        { duration: '30s', target: 5000 },
        { duration: '5m', target: 5000 },
        { duration: '3s', target: 0 },
    ],
};

export default function () {
    const url = 'wss://chat-kdev-kr.kollus.com/ws/N1_01JH7PKPT841ECCAY3G2DVAJCF/jungin.kim@catenoid.net/Admin?joinMessageCount=30';

    const res = ws.connect(url, {}, (socket) => {
        // WebSocket 연결 성공 시 이벤트 설정
        socket.on('open', () => {
            console.log('WebSocket connected');
        });

        socket.on('message', (message) => {
            console.log(`Received message: ${message}`);
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