# Go Bank Transfer :bank: :money_with_wings:
- Go Bank Transfer is a simple API for some banking routines, such as creating accounts, listing accounts, listing balance for a specific account, transfers between accounts and listing transfers.

## Architecture
-  This is an attempt to implement a clean architecture, in case you don't already know, here is a reference https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

![Clean Architecture](cleanarch.png)

## Requirements/dependencies
- Golang (not obligatory)
- Docker
- Docker-compose

## Getting Started

- Environment variables

```sh
create .env through .env.example
```

- Starting API

```sh
make up
```

- Run tests in container (it is necessary to have the application started)

```sh
make test-container
```

- Run tests local (it is necessary to have golang installed)

```sh
make test
```

- View logs

```sh
make logs
```

- Enter in container

```sh
make enter-container
```

#### Test endpoints API using Postman

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/38286dacebab6fa2975f)

#### Test endpoints API using curl

- Creating new account

```bash
curl -i --request POST 'http://localhost:3001/api/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Test",
    "cpf": "070.910.584-64",
    "balance": 1.00
}'
```

- Listing accounts

```bash
curl -i --request GET 'http://localhost:3001/api/accounts'
```

- Fetching account balance

```bash
curl -i --request GET 'http://localhost:3001/api/accounts/{{account_id}}/balance'
```

- Creating new transfer

```bash
curl -i --request POST 'http://localhost:3001/api/transfers' \
--header 'Content-Type: application/json' \
--data-raw '{
	"account_destination_id": "{account_id}",
	"account_origin_id": "{account_id}",
	"amount": 1.00
}'
```

- Listing transfers

```bash
curl -i --request GET 'http://localhost:3001/api/transfers'
```

## Git workflow
- Gitflow

## Code status
- Development

## Author
- Gabriel Sabadini Facina - [GSabadini](https://github.com/GSabadini)

## License
Copyright © 2020 [GSabadini](https://github.com/GSabadini).<br />
This project is [MIT](LICENSE) licensed.
