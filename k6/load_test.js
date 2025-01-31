import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 50 }, // Ramp up to 50 users
        { duration: '1m', target: 50 },  // Stay at 50 users
        { duration: '30s', target: 0 },  // Ramp down to 0
    ],
};

export default function () {
    const url = 'http://localhost:8080';
    // 生成するユーザーのユニークな識別子
    const uniqueId = Math.random().toString(36).substring(7);

    // Create User
    let res = http.post(`${url}/v1/users`, JSON.stringify({
        name: `Test User ${uniqueId}`,
        email: `test_${uniqueId}@example.com`,
    }), { headers: { 'Content-Type': 'application/json' } });

    const userCreated = check(res, {
        'create user status is 200': (r) => r.status === 200,
    });

    if (!userCreated) {
        // ユーザー作成に失敗した場合、後続のリクエストをスキップ
        return;
    }

    // Extract user ID
    let userId = res.json().id;

    // Create Product
    res = http.post(`${url}/v1/products`, JSON.stringify({
        name: `Test Product ${uniqueId}`,
        price: 99.99,
    }), { headers: { 'Content-Type': 'application/json' } });

    const productCreated = check(res, {
        'create product status is 200': (r) => r.status === 200,
    });

    if (!productCreated) {
        // 商品作成に失敗した場合、後続のリクエストをスキップ
        return;
    }

    // Extract product ID
    let productId = res.json().id;

    // Create Order
    res = http.post(`${url}/v1/orders`, JSON.stringify({
        user_id: userId,
        product_id: productId,
        quantity: 2,
    }), { headers: { 'Content-Type': 'application/json' } });

    const orderCreated = check(res, {
        'create order status is 200': (r) => r.status === 200,
    });

    if (!orderCreated) {
        // 注文作成に失敗した場合、後続のリクエストをスキップ
        return;
    }

    // Extract order ID
    let orderId = res.json().id;

    // Get Order
    res = http.get(`${url}/v1/orders/${orderId}`);
    check(res, {
        'get order status is 200': (r) => r.status === 200,
    });

    // Add Inventory
    res = http.post(`${url}/v1/inventory`, JSON.stringify({
        product_id: productId,
        quantity: 10,
    }), { headers: { 'Content-Type': 'application/json' } });

    check(res, {
        'add inventory status is 200': (r) => r.status === 200,
    });

    // Get Inventory
    res = http.get(`${url}/v1/inventory/${productId}`);
    check(res, {
        'get inventory status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
