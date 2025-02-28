import ws from 'k6/ws';
import { check } from 'k6';
import { Trend, Counter } from 'k6/metrics';

// 커스텀 메트릭 정의
const wsConnecting = new Trend('ws_connecting'); // 연결 시간
const wsMsgsReceived = new Counter('ws_msgs_received'); // 수신된 메시지 수
const wsMsgsSent = new Counter('ws_msgs_sent'); // 전송된 메시지 수
const wsPing = new Trend('ws_ping'); // Ping-Pong 시간
const wsSessionDuration = new Trend('ws_session_duration'); // 세션 지속 시간
const wsSessions = new Counter('ws_sessions'); // 시작된 세션 수

export const options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    const url = 'ws://example.com/socket';

    const response = ws.connect(url, {}, (socket) => {
        const connectStart = Date.now();
        wsSessions.add(1); // 세션 시작 카운트

        socket.on('open', () => {
            // 연결 시간 기록
            wsConnecting.add(Date.now() - connectStart);

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

        // 에러 핸들링
        socket.on('error', (e) => {
            console.error(`WebSocket error: ${e.error()}`);
        });

        // 세션 종료 시간 측정
        const sessionStart = Date.now();
        socket.on('close', () => {
            wsSessionDuration.add(Date.now() - sessionStart);
        });

        // 5초 후 연결 종료
        socket.setTimeout(() => {
            socket.close();
        }, 5000);
    });

    // 연결 성공 여부 확인
    check(response, {
        'status is 101': (r) => r && r.status === 101,
    });
}
