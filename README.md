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
