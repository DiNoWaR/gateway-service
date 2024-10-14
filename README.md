## Payment Gateway Service
### Overview
This project implements a payment gateway system consisting of three main services:

- Gateway Service: The core service responsible for handling deposit and withdrawal requests, integrating with multiple payment gateways, and processing transactions asynchronously.
- SOAP Gateway (Mock): A mock service simulating a payment provider using the SOAP protocol.
- REST Gateway (Mock): A mock service simulating a payment provider using the REST protocol.
The project is fully containerized using Docker and Docker Compose. PostgreSQL is used as the database to store transaction information.

Features
Gateway Service: Handles deposit and withdrawal requests, processes them asynchronously, and notifies clients via callbacks.
SOAP & REST Gateway Mocks: Simulate external payment gateways that process transactions (deposits and withdrawals).
Asynchronous Processing: Transactions are processed in the background, and clients are notified through a callback once the transaction is complete.
PostgreSQL Database: Stores transaction details, including status, amount, and timestamps.
Dockerized Environment: The entire system is containerized using Docker, and all services, including the database, can be started with Docker Compose.
Architecture
The project consists of the following services:

#### Gateway Service:
Receives deposit and withdrawal requests from clients.
Validates requests and forwards them to the appropriate payment gateway (either REST or SOAP).
Processes transactions asynchronously and stores transaction details in PostgreSQL.
Notifies clients via a callback once the transaction is complete.

#### SOAP Gateway (Mock):
- Simulates a SOAP-based payment gateway.
Handles deposit and withdrawal requests, returning mocked transaction results.
Responds to requests from the Gateway Service and processes transactions.

#### REST Gateway (Mock):
- Simulates a REST-based payment gateway.
Handles deposit and withdrawal requests, returning mocked transaction results.
Responds to requests from the Gateway Service and processes transactions.

#### Postgres Database:
Stores all transaction data, including transaction IDs, reference IDs, account details, status, and timestamps.
Provides a reliable data store for querying and managing transactions.


### Prerequisites
Api spec file is located in the folder **api**. 
To launch the entire service run the command from **dev** folder
```
docker compose up
```
After launched you can make all user requests

### Example Requests
#### Deposit Request
Send a deposit request to the Gateway Service:
```
    curl -X POST http://localhost:9090/deposit \
     -H "Content-Type: application/json" \
     -d '{
           "amount": 100.50,
           "currency": "USD",
           "account_id": "ACC123",
           "gateway_id": "rest_gateway"
         }'
```

#### Withdrawal Request
Send a withdrawal request to the Gateway Service:

```
curl -X POST http://localhost:9090/withdraw \
     -H "Content-Type: application/json" \
     -d '{
           "amount": 50.00,
           "currency": "EUR",
           "account_id": "ACC123",
           "gateway_id": "soap_gateway"
         }'
```

#### Get Transaction Request
```
curl -X POST http://localhost:9090/transaction \
     -H "Content-Type: application/json" \
     -d '{"reference_id": "5da37158-d41d-4280-bcef-2e88b12214e6"}'
```

#### Get All User Transactions Request
```
curl -X POST http://localhost:9090/transactions \
     -H "Content-Type: application/json" \
     -d '{"account_id": "ACC123"}'
```