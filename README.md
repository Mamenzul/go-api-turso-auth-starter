# Project Name

## Technologies Used

- Go chi API framework
- Turso DB
- go-pkgz/auth for direct authentication

## Description

This project is a Go API built using the chi framework. It utilizes Turso DB as the database and go-pkgz/auth for direct authentication. The API provides various endpoints for interacting with the database and authenticating users.

## Installation

To install and set up the project, follow these steps:

1. Clone the repository: `git clone https://github.com/Mamenzul/go-api-starter.git`
2. Change to the project directory: `cd go-api-starter`
3. Install dependencies: `go mod download`
4. Set up the database: [Turso quickstart](https://docs.turso.tech/quickstart)
5. Set up the database:
   - Create a `.env` file in the project directory.
   - Add the following line to the `.env` file:
     ```
     DATABASE_URL=[url]?authToken=[token]
     ```
   - Save the `.env` file.
6. Build and run the API: `make run`

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). Please see the [LICENSE](./LICENSE) file for more details.
