# RealPay

    This is a simple payment server application built in Go. It provides an API for managing the lifecycle of payment transactions, including initialization, confirmation, rejection, and status retrieval.

## Features

- **Transaction Initialization**: Clients can create new payment transactions by providing details like the amount, currency, payer phone number, etc. The server generates a unique transaction ID and returns the transaction details along with a payment URL.
- **Payment Confirmation**: The server offers an endpoint to confirm a payment transaction. This would typically involve verifying the transaction status, checking the account balance, and updating the transaction status accordingly.
- **Payment Rejection**: Clients can use the server to reject a payment transaction, updating the status to "Failed" and sending a webhook notification.
- **Payment Status Retrieval**: Clients can query the server to get the current status of a payment transaction.
- **Utility Functions**: The server provides helper functions for generating unique IDs, sending JSON responses, and notifying external systems about payment events.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Gorilla Mux library (for routing)

### Installation

1. Clone the repository:
```
git clone https://github.com/ndg23/realpay.git
```

2. Change to the project directory:
```
cd realpay
```

3. Install the dependencies:
```
go get -d ./...
```

### Running the Server

1. Build the application:
```
go build -o realpay cmd/main.go
```

2. Run the server:
```
./realpay
```

The server will start listening on `http://localhost:8080`.

## API Endpoints

- `POST /payments`: Initialize a new payment transaction
- `PUT /payments/{id}/confirm`: Confirm a payment transaction
- `PUT /payments/{id}/reject`: Reject a payment transaction
- `GET /payments/{id}`: Retrieve the status of a payment transaction

## Future Improvements

- Implement database integration for storing and retrieving payment transactions
- Add support for more payment methods (e.g., credit card, mobile wallet)
- Improve error handling and logging
- Implement authentication and authorization mechanisms
- Add unit tests and integration tests

## Contributing

If you'd like to contribute to this project, please follow the standard GitHub workflow:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes and commit them
4. Push your changes to your fork
5. Submit a pull request

## License

This project is licensed under the [MIT License](LICENSE).