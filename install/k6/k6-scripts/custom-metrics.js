import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';
import ws from 'k6/ws';

const wsMsgsReceived = new Counter('sum_ws_msgs_received'); // 수신된 메시지 수
const wsMsgsSent = new Counter('sum_ws_msgs_sent'); // 전송된 메시지 수
const wsSessions = new Counter('sum_ws_sessions'); // 시작된 세션 수
const wsPing = new Trend('time_ws_ping'); // Ping-Pong 시간


export const options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: '5s', target: 0 },
    ],
};

export default function () {
    const url = 'wss://chat-kdev-kr.kollus.com/ws/N1_01JGZVP3BWY10GNRSX7R0C0QY3/jungin.kim@catenoid.net%2FAdmin?joinMessageCount=30';

    const res = ws.connect(url, {}, (socket) => {

        // WebSocket 연결 성공 시 이벤트 설정
        socket.on('open', () => {
            console.log(`user-${__VU} WebSocket connected`);

            // 메시지 전송
            const message = JSON.stringify({ type: 'ping' });
            socket.send(message);
            wsMsgsSent.add(1);

            // Ping-Pong 시간 측정
            const pingStart = Date.now();
            socket.on('message', (data) => {
                wsMsgsReceived.add(1); // 메시지 수신 카운트
                const parsedData = JSON.parse(data);
                if (parsedData.type === 'pong') {
                    wsPing.add(Date.now() - pingStart);
                }
            });
        });

        socket.on('message', (message) => {
            console.log(`user-${__VU} Received message: ${message}`);
        });

        socket.setTimeout(() => {
            socket.close();
        }, 30000);
        // 세션 종료 시간 측정
        socket.on('close', () => {
            console.log(`user-${__VU} WebSocket disconnected`);
        });

        socket.on('error', (e) => {
            console.error(`WebSocket error: ${e}`);
        });
    });

    check(res, {
        'WebSocket connection established': (r) => r && r.status === 101,
    });
}
