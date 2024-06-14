import http from 'k6/http';
import { sleep } from 'k6';

const url = "http://localhost:8080/api/v1/backends?region=europe"
const timeout = 0.5
export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'shared-iterations',
            vus: 1000,
            iterations: 500000,
            maxDuration: '5m',
        },
    },
};

export default function () {
    http.get(url);
    sleep(timeout);
}