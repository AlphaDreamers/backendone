# Chat and Order Flow Documentation

This document explains how the frontend (Next.js) should interact with the backend services for order placement and chat functionality.

## Flow Overview

1. **Order Placement**: Buyer creates an order which generates a chat room
2. **Chat Initialization**: Both buyer and seller connect to the chat room via WebSocket
3. **Messaging**: Participants exchange messages in the chat room
4. **Notifications**: Seller receives real-time order notifications via SSE

## API Endpoints

### 1. Order Placement

**Endpoint**: `POST /orders`

**Request**:
```json
{
  "buyerId": "550e8400-e29b-41d4-a716-446655440000",
  "sellerId": "123e4567-e89b-12d3-a456-426614174000",
  "serviceId": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

**Response**:
```json
{
    "chat_room": "3ea676c5-9b66-40a9-be09-ac198b651407"
}
```

**Frontend Implementation**:
```javascript
async function placeOrder(orderData) {
  const response = await fetch('/orders', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(orderData),
  });
  return await response.json();
}

// Usage
const orderResponse = await placeOrder({
  buyerId: '550e8400-e29b-41d4-a716-446655440000',
  sellerId: '123e4567-e89b-12d3-a456-426614174000',
  serviceId: 'f47ac10b-58cc-4372-a567-0e02b2c3d479'
});
```

### 2. WebSocket Chat Connection

**Endpoint**: `ws://yourdomain.com/ws/chat`

**Initialization Message** (must be sent first after connection):
```json
{
  "userId": "123e4567-e89b-12d3-a456-426614174000",
  "chatRoomId": "3ea676c5-9b66-40a9-be09-ac198b651407"
}
```

**Frontend Implementation**:
```javascript
function setupChatConnection(userId, chatRoomId) {
  const socket = new WebSocket('ws://yourdomain.com/ws/chat');
  
  socket.onopen = () => {
    // Send initialization message
    socket.send(JSON.stringify({
      userId,
      chatRoomId
    }));
  };
  
  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    // Handle incoming message
    console.log('Received message:', message);
  };
  
  socket.onclose = () => {
    console.log('WebSocket connection closed');
  };
  
  return socket;
}

// Usage for seller
const sellerSocket = setupChatConnection(
  '123e4567-e89b-12d3-a456-426614174000',
  '3ea676c5-9b66-40a9-be09-ac198b651407'
);

// Usage for buyer
const buyerSocket = setupChatConnection(
  '550e8400-e29b-41d4-a716-446655440000',
  '3ea676c5-9b66-40a9-be09-ac198b651407'
);
```

### 3. Sending Messages

**Message Format**:
```json
{
  "from": "550e8400-e29b-41d4-a716-446655440000",
  "to": "123e4567-e89b-12d3-a456-426614174000",
  "chat_room_id": "3ea676c5-9b66-40a9-be09-ac198b651407",
  "body": "Hey there! Just wanted to check in.",
  "file": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/4gHYSUNDX1BST0ZJTEUAAQEAAAHIAAAAAAQwAABtbnRyUkdCIFhZWiAH4AABAAEAAAAAAABh...[truncated]...ppOaoqeoh7b76Z//9k="
}
```

**Frontend Implementation**:
```javascript
function sendMessage(socket, messageData) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(messageData));
  } else {
    console.error('WebSocket is not open');
  }
}

// Usage
sendMessage(buyerSocket, {
  from: '550e8400-e29b-41d4-a716-446655440000',
  to: '123e4567-e89b-12d3-a456-426614174000',
  chat_room_id: '3ea676c5-9b66-40a9-be09-ac198b651407',
  body: 'Hey there! Just wanted to check in.',
  file: 'data:image/jpeg;base64,...'
});
```

### 4. Order Notifications (SSE)

**Endpoint**: `GET /sse/seller/:sellerId`

**Frontend Implementation**:
```javascript
function setupSellerNotifications(sellerId) {
  const eventSource = new EventSource(`/sse/seller/${sellerId}`);
  
  eventSource.onmessage = (event) => {
    console.log('New notification:', event.data);
    // Display notification to seller
  };
  
  eventSource.onerror = (error) => {
    console.error('EventSource failed:', error);
  };
  
  return eventSource;
}

// Usage
const notifications = setupSellerNotifications('123e4567-e89b-12d3-a456-426614174000');
```

## Complete Flow Example

1. **Buyer places an order**:
   ```javascript
   const orderResponse = await placeOrder({
     buyerId: '550e8400-e29b-41d4-a716-446655440000',
     sellerId: '123e4567-e89b-12d3-a456-426614174000',
     serviceId: 'f47ac10b-58cc-4372-a567-0e02b2c3d479'
   });
   
   const chatRoomId = orderResponse.chat_room;
   ```

2. **Buyer connects to chat**:
   ```javascript
   const buyerSocket = setupChatConnection(
     '550e8400-e29b-41d4-a716-446655440000',
     chatRoomId
   );
   ```

3. **Seller connects to chat and sets up notifications**:
   ```javascript
   const sellerSocket = setupChatConnection(
     '123e4567-e89b-12d3-a456-426614174000',
     chatRoomId
   );
   
   const notifications = setupSellerNotifications('123e4567-e89b-12d3-a456-426614174000');
   ```

4. **Participants exchange messages**:
   ```javascript
   // Buyer sends message
   sendMessage(buyerSocket, {
     from: '550e8400-e29b-41d4-a716-446655440000',
     to: '123e4567-e89b-12d3-a456-426614174000',
     chat_room_id: chatRoomId,
     body: 'Hello, I have a question about my order'
   });
   
   // Seller sends message with image
   sendMessage(sellerSocket, {
     from: '123e4567-e89b-12d3-a456-426614174000',
     to: '550e8400-e29b-41d4-a716-446655440000',
     chat_room_id: chatRoomId,
     body: 'Here is the image you requested',
     file: 'data:image/jpeg;base64,...'
   });
   ```

## Error Handling

1. **WebSocket Errors**:
    - Check `readyState` before sending messages
    - Implement reconnection logic if connection drops

2. **File Uploads**:
    - Compress images before sending to reduce size
    - Handle upload failures gracefully

3. **Order Placement**:
    - Validate all UUIDs before sending
    - Handle 400/500 errors appropriately
