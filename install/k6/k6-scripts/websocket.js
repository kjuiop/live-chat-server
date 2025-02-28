import { check } from 'k6';
import ws from 'k6/ws';

export const options = {
    vus: 10,  // 가상 사용자 수 (10명)
    duration: '30s',  // 테스트 지속 시간 (30초)
};

export default function () {
    const url = 'ws://host.docker.internal:8090/ws/chat/join/rooms/N2-01JG8HSMQQCV74V35251MDGDFH/user/jake';

    const res = ws.connect(url, {}, (socket) => {
        // WebSocket 연결 성공 시 이벤트 설정
        socket.on('open', () => {
            console.log('WebSocket connected');

            const message = {
                Method: 'chat',
                SendUserId: 'jungin-kim',
                Message: 'hello',
            };

            socket.send(JSON.stringify(message));
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
