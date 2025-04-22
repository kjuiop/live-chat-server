import ws from 'k6/ws';
import { check } from 'k6';
import { Trend } from 'k6/metrics';

// Metrics 정의
const connectTime = new Trend('websocket_connect_time');
const messageLatency = new Trend('websocket_message_latency');

export const options = {
    scenarios: {
        chat_test: {
            executor: 'per-vu-iterations', // 각 사용자가 고정된 횟수로 실행
            vus: 1000, // 최대 가상 사용자 수
            iterations: 1, // 사용자당 1회 실행
            maxDuration: '5m', // 전체 테스트 최대 시간
        },
    },
};

export default function () {
    const url = `ws://host.docker.internal:8090/ws/chat/join/rooms/N1-01JGRGTBYTM5K3ZWCGF4PMNKEB/user/user_${__VU}`;
    const params = { tags: { user: `user_${__VU}` } };

    const res = ws.connect(url, params, function (socket) {
        const connectStart = Date.now();

        socket.on('open', () => {
            connectTime.add(Date.now() - connectStart);
            console.log(`Connected: user_${__VU}`);

            if (__VU <= 50) {
                // 50명은 메시지를 지속적으로 전송
                const interval = setInterval(() => {
                    socket.send(
                        JSON.stringify({
                            Method: 'chat',
                            SendUserId: `user_${__VU}`,
                            Message: `Message from user_${__VU} at ${Date.now()}`,
                        })
                    );
                }, 1000); // 1초 간격 메시지 전송

                // 30초 후 메시지 전송 중단 및 연결 유지
                socket.setTimeout(() => {
                    clearInterval(interval);
                }, 30000);
            } else {
                // 950명은 1분 대기 후 강제 종료
                socket.setTimeout(() => {
                    console.log(`Forcefully disconnecting user_${__VU}`);
                    socket.close(); // 강제 종료
                }, 60000); // 1분 대기 후 종료
            }
        });

        socket.on('message', (message) => {
            const receivedTime = Date.now();
            const payload = JSON.parse(message);

            if (payload.type === 'message') {
                messageLatency.add(receivedTime - payload.timestamp);
                console.log(`Message received by user_${__VU}: ${message}`);
            }
        });

        socket.on('close', () => {
            console.log(`Disconnected: user_${__VU}`);
        });

        socket.on('error', (e) => {
            console.error(`Error for user_${__VU}: ${e.error()}`);
        });
    });

    check(res, {
        'status is 101': (r) => r && r.status === 101, // WebSocket 연결 성공 여부
    });
}
