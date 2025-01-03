import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
    stages: [
        {
            duration: '2m',
            target: 500,
        },
        {
            duration: '3m',
            target: 500,
        },
    ],
};
export default function () {
    http.get('https://chat-kdev-kr.kollus.com/ws/N1/system/health-check/local');
    sleep(1);
}
