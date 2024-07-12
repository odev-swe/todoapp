# TodoApp - Go

Although I understand that a todo app is quite a common project in the development world, I want to use it as a starting point for my learning journey. By incorporating the latest frameworks and technologies, I aim to enhance my technical skills.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Framework](#framework)
- [Features](#features)
- [License](#license)
- [Contributing](#contributing)
- [Contact](#contact)

## Installation

<!-- Important -->

### Important

Before you start, make sure you have the following prerequisites installed on your machine and go path is set up correctly otherwise you will face issues while running the application.:

Prerequisites:

- [Go](https://golang.org/dl/) (at least version 1.20)
- [Goose](https://github.com/pressly/goose) (for database migrations)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/) (optional)
- [Postman](https://www.postman.com/downloads/) (optional)
- [Swaggo](https://github.com/swaggo/swag)
- [dlv](https://github.com/go-delve/delve) (for debugging)

### Clone the repository

```bash
# Clone the repository
git clone https://github.com/odev-swe/todoapp.git

# Change the directory
cd todoapp

# Use docker-compose to start the application (simple and fast)
make dc-up

# If you want to terminate the application
make dc-down

# If you want to refer any other command, you can list all available commands or refer to the Makefile
make list
```

### Environment Variables

Refer to .env file to set up the environment variables.

## Usage

1. Swagger Documentation:
   Open your browser and go to http://localhost:3000/swagger/index.html to access the Swagger documentation. This provides a user-friendly interface to interact with the API.

## Framework

libraries and frameworks used in this project:

1. [Chi]() - A lightweight, idiomatic and composable router for building Go HTTP services.
2. [pgx]() - A PostgreSQL driver and toolkit for Go.
3. [goose]() - A database migration tool for Go.
4. [swaggo]() - Automatically generate RESTful API documentation with Swagger 2.0 for Go.
5. [testify]() - A toolkit with common assertions and mocks that plays nicely with the standard library.
6. [postgreSQL]() - A powerful, open-source object-relational database system.
7. [Docker]() - A platform for developing, shipping, and running applications in containers.
8. [Docker Compose]() - A tool for defining and running multi-container Docker applications.
9. [Make]() - A build automation tool that automatically builds executable programs and libraries from source code by reading files called Makefiles.
10. [zap]() - A fast, structured, leveled logging library for Go.
11. [jwt]() - A Go implementation of JSON Web Tokens (JWT).
12. [bcrypt]() - A Go package that provides functions for hashing passwords using the bcrypt adaptive hashing algorithm.

## Features

- [x] Authentication using JWT
- [x] Todos (CRUD)
- [x] Rate Limiting
- [ ] Caching
- [ ] Notifcation
- [ ] Testing
- [ ] CI/CD Pipeline with Github Actions (own runner)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

- **I'm a beginner in Go, so any help is appreciated.**

- **If you have any ideas, you are welcome to contribute as well and maybe can make this code a reference for those who want to learn.**

- **If you think my code has a bad format, feel free to help me improve it. We can learn together!**

If you would like to contribute to the project, please follow these steps:
Fork the repository
Create a new branch (git checkout -b feature-branch)
Make your changes
Commit your changes (git commit -m 'Add new feature')
Push to the branch (git push origin feature-branch)
Open a Pull Request

## Contact

Oscar - odevswe@gmail.com

Project Link: https://github.com/odev-swe/todoapp
