<h1 align="center">Welcome to Go Bank Transfer :bank:</h1>
<p>
  <a href="https://github.com/fabianoleittes/go-bank-transfer" target="_blank">
    <img src="https://github.com/fabianoleittes/go-bank-transfer/workflows/Test/badge.svg"
  </a>
  <img alt="Version" src="https://img.shields.io/badge/version-1.7.0-blue.svg?cacheSeconds=2592000" />
  <a href="https://goreportcard.com/badge/github.com/GSabadini/go-bank-transfer" target="_blank">
    <img alt="Build" src="https://goreportcard.com/badge/github.com/GSabadini/go-bank-transfer" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
  <a href="https://travis-ci.org/github/GSabadini/go-bank-transfer" target="_blank">
    <img alt="Build" src="https://travis-ci.org/GSabadini/go-bank-transfer.svg?branch=master" />
  </a>
  <a href="https://codecov.io/gh/GSabadini/go-bank-transfer">
    <img src="https://codecov.io/gh/GSabadini/go-bank-transfer/branch/master/graph/badge.svg" />
  </a>
</p>

- Go Bank Transfer is a simple API for some banking routines, such as creating accounts, listing accounts, listing balance for a specific account, transfers between accounts and listing transfers.

## Architecture
-  This is an attempt to implement a clean architecture, in case you don’t know it yet, here’s a reference https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

![Clean Architecture](cleanarch.png)

## Example Create Account Usecase

![Clean Architecture](create_account.png)

## Requirements/dependencies
- Docker
- Docker-compose

## Getting Started

- Environment variables

```sh
make init
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

## API Request

| Endpoint        | HTTP Method           | Description       |
| --------------- | :---------------------: | :-----------------: |
| `/v1/accounts` | `POST`                | `Create accounts` |
| `/v1/accounts` | `GET`                 | `List accounts`   |
| `/v1/accounts/{{account_id}}/balance`   | `GET`                |    `Find balance account` |
| `/v1/transfers`| `POST`                | `Create transfer` |
| `/v1/transfers`| `GET`                 | `List transfers`  |
| `/v1/health`| `GET`                 | `Health check`  |

#### Test endpoints API using Postman

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/8406204152f98cc33eac)

#### Test endpoints API using curl

- Creating new account

```bash
curl -i --request POST 'http://localhost:3001/v1/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Test",
    "cpf": "070.910.584-24",
    "balance": 100
}'
```

- Listing accounts

```bash
curl -i --request GET 'http://localhost:3001/v1/accounts'
```

- Fetching account balance

```bash
curl -i --request GET 'http://localhost:3001/v1/accounts/{{account_id}}/balance'
```

- Creating new transfer

```bash
curl -i --request POST 'http://localhost:3001/v1/transfers' \
--header 'Content-Type: application/json' \
--data-raw '{
	"account_origin_id": "{{account_id}}",
	"account_destination_id": "{{account_id}}",
	"amount": 100
}'
```

- Listing transfers

```bash
curl -i --request GET 'http://localhost:3001/v1/transfers'
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
