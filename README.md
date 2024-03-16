# User Server Project

This project is a user server that provides registration, login, email verification, and recommendation functionalities. It is built with Golang, Gin framework, MySQL, and Redis.

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites
You need to have Docker and Docker Compose installed on your machine.

### Installation and Setup

1. Clone the repository.
2. Run `docker-compose up` to start the application and its dependencies.
3. The server should now be running on `localhost:8080`.


## API Endpoints
* `POST /register`: Register a new user.
* `POST /login`: Login with email and password.
* `GET /verify_email`: Verify the user's email with a verification code.
* `GET /recommendation`: Get a list of recommended items (requires authentication).