# eoplatform

Event organizer platform for MSIB mini porject.

## Features

- Authentication
- Role-based authorization
- Account menagement
- Bank account management
- CRUD for EO services
- Customer order with payment gateway integration
- Customer feedback with sentiment analysis

## Requirements

- Go v1.19.1
- Goose v3.7.0
- Docker 20.10.20
- Docker compose v2.12.0
- MySQL v8.0
- GCP credentials for gcloud
- Midtrans server key
- Terraform v1.3.3 (optional)

## Usage

1. Download and install required dependencies
2. Refer to [Google Cloud documentation](https://cloud.google.com/natural-language/docs/setup) to setup Natural Language API
3. Refer to [Midtrans documentation](https://api-docs.midtrans.com/) to setup environment and retrieve server key
4. Fill all variables in `.env` file (you also need to fill `Makefile` and `docker-compose.yaml` if you want to use them)
5. Create a new database and run migration using `make migrateup`
6. Run the app!

## Directories

| Folder         | Description                                 |
| -------------- | ------------------------------------------- |
| eoplatform     | Root folder                                 |
| ├── config     | Application configurations                  |
| ├── db         | Database connection                         |
| ├── helper     | Custom helper functions                     |
| ├── migration  | SQL files for migration                     |
| ├── model      | Database models                             |
| ├── repository | Database access interfaces                  |
| ├── request    | HTTP request objects                        |
| ├── response   | HTTP response objects                       |
| ├── server     | Server objects--including handlers & routes |
| ├── terraform  | Infrastructure configurations               |
| └── test       | Test cases                                  |

## Terraform Vars Example

```
project = "project_id"
region  = "asia-southeast2"
zone    = "asia-southeast2-a"

env_vars = [
  {
    name : "APP_ENV",
    value : "production",
  },
  {
    name : "HTTP_PORT",
    value : "8000",
  },
  {
    name : "AUTH_SECRET",
    value : "ssssshhhhhhhhhh",
  },
  {
    name : "AUTH_COST",
    value : "10",
  },
  {
    name : "AUTH_EXP_HOURS",
    value : "1",
  },
  {
    name : "DB_DRIVER",
    value : "mysql",
  },
  {
    name : "DB_USER",
    value : "root",
  },
  {
    name : "DB_PASS",
    value : "root",
  },
  {
    name : "DB_NAME",
    value : "eoplatform",
  },
  {
    name : "DB_HOST",
    value : "localhost",
  },
  {
    name : "DB_PORT",
    value : "3306",
  },
  {
    name : "SMTP_HOST",
    value : "smtp.gmail.com",
  },
  {
    name : "SMTP_PORT",
    value : "587",
  },
  {
    name : "EMAIL_ADDRESS",
    value : "email_or_username",
  },
  {
    name : "EMAIL_PASSWORD",
    value : "password",
  },
  {
    name : "MIDTRANS_BASE_URL",
    value : "https://api.sandbox.midtrans.com",
  },
  {
    name : "MIDTRANS_SERVER_KEY",
    value : "server_key",
  },
]
```
