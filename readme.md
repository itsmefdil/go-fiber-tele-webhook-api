### REST API BOT TELEGRAM WEBHOOK

### Description
This is a simple REST API that allows you to create, update, delete and send messages to a chat using a bot. The bot is created using the Telegram API.

### Features
- Basic Authentication
- Create a bot
- Get all bots
- Update a bot
- Delete a bot
- Send a message to a chat

### How to use

1. Clone the repository
2. Install Depedenices
```
go mod tidy
```
3. Copy the .env.example file to .env and fill in the required values
4. Run the application
```
go run main.go
```


### API Endpoints
1. Create a bot
params: 
Method : POST
URL : http://localhost:3000/bots
```
{
    "token": "your_token",
    "room_id": "your_chat_id",
    "thread_id": "your_thread_id",
}

```
2. Get all bots
Method : GET

```
curl -X GET "http://localhost:3000/bots"
```
3. Update a bot
Method : PUT
URL : http://localhost:3000/bots/:id
```
{
    "token": "your_token",
    "room_id": "your_chat_id",
    "thread_id": "your_thread_id",
}
```

4. Delete a bot
Method : DELETE
URL : http://localhost:3000/bots/:id

5. Send a message to a chat
Method : POST


### Technologies
- Go
- Fiber

### API TABLE   

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /bots | Get all bots |
| POST | /bots | Create a bot |
| PUT | /bots/:id | Update a bot |
| DELETE | /bots/:id | Delete a bot |
| POST | /webhook/:id/send | Send a message to a chat |