import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    stages: [
        { duration: '2m', target: 100 }, // Ramp-up to 100 users over 2 minutes
        { duration: '5m', target: 100 }, // Stay at 100 users for 5 minutes
        { duration: '1m', target: 0 },   // Ramp-down to 0 users over 1 minute
    ],
};

// Generate random emails and usernames for the POST request
function getRandomUser() {
    const name = `User_${randomString(5)}`;
    const email = `${name.toLowerCase()}@example.com`;
    return { name: name, email: email, age: Math.floor(Math.random() * 40) + 20 };
}

export default function () {
    // GET request to /api/v1/order
    const getOrderRes = http.get('http://localhost:8080/api/v1/order');
    check(getOrderRes, {
        'GET /order status is 200': (r) => r.status === 200,
    });

    // POST request to /api/v1/user with random user data
    const user = getRandomUser();
    const payload = JSON.stringify(user);
    const params = { headers: { 'Content-Type': 'application/json' } };

    console.log(payload)

    const postUserRes = http.post('http://localhost:8080/api/v1/user', payload, params);
    check(postUserRes, {
        'POST /user status is 201': (r) => 
             //console.log('status is ', r.status)
            r.status === 201,
        
    });

    sleep(1);
}
