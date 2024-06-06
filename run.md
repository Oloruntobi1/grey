## Step-by-Step Guide to Run the Application
After cloing the repository please following the following steps to run the application locally.

### 1. Install Docker
First, ensure Docker is installed on your machine. You can download it from [Docker's official website](https://www.docker.com/get-started).

### 2. Start All Services
Run the following command to start all services using Docker:
```sh
make start-all-services
```

### 3 Create User
```sh
curl -X POST http://localhost:9292/api/create-user \
-H "Content-Type: application/json" \
-d '{
  "name": "userA",
  "email": "userA@example.com"
}'
```

### 4 Create Another User
```sh
curl -X POST http://localhost:9292/api/create-user \
-H "Content-Type: application/json" \
-d '{
  "name": "userB",
  "email": "userB@example.com"
}'
```

### 5 Create Wallet for User A
```sh
curl -X POST http://localhost:9292/api/create-wallet \
-H "Content-Type: application/json" \
-d '{
  "user_id": user-id-for-1,
  "initial_balance": 1000
}'
```

### 6 Create Wallet for User B
```sh
curl -X POST http://localhost:9292/api/create-wallet \
-H "Content-Type: application/json" \
-d '{
  "user_id": user-id-for-2,
  "initial_balance": 500
}'
```


## FUTURE WORK

### 7 Transfer from User A to User B
```sh
curl -X POST http://localhost:9292/api/transfer \
-H "Content-Type: application/json" \
-d '{
  "from_user_id": user-id-for-1,
  "to_user_id": user-id-for-2,
  "amount": 100
}'
```

### 8 Get User A's Transaction List
```sh
curl -X GET http://localhost:9292/api/users/user-id-for-1/transactions
```

### 9 Get User B's Transaction List
```sh
curl -X GET http://localhost:9292/api/users/user-id-for-2/transactions
```