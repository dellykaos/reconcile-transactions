# Reconciliation Service

This service is used to identify unmatched and discrepant transactions between internal data (system transactions) and external data (bank statements).

The tech stack used for this services are:

- Golang version 1.23.0
- PostgreSQL 14.8

**Table of Content:**

- [How To Run](#how-to-run)
  - [Setup Environment](#1-setup-environment)
  - [Setup Database](#2-setup-database)
  - [Run API](#3-run-api)
  - [Run Reconcile Job](#4-run-reconcile-job)
  - [Create Docker Container](#5-create-docker-container-for-deployment)
- [Documentation](#documentation)
  - [Get List Reconcile Job](#get-reconciliation-list)
  - [Get Reconcile Job By ID](#get-reconciliation-job-request-by-id)
  - [Create Reconcile Job](#create-reconciliation-job-request)
  - [Process Reconcile Job](#reconciliation-job-process)

## How to run

### 1. Setup environment

Copy file `env.sample` to `.env`

```shell
cp env.sample .env
```

Fill all required configuration needed.

```text
ENV=development

DB_HOST=localhost # your database host
DB_PORT=5432 # your database port
DB_USER= # your database user
DB_PASS= # your database password
DB_NAME= # your database name

SERVER_PORT=8080 # server API port

USE_LOCAL_STORAGE=true # use local storage as file storage
LOCAL_STORAGE_DIR=/temp_storage # dir location to store uploaded csv files, currently would use path $CWD/$LOCAL_STORAGE_DIR

# This GCS Config would not be used if USE_LOCAL_STORAGE value is true
GCS_KEY_JSON= # GCS Secret Key JSON to access Bucket
GCS_BUCKET= # GCS Bucket Name
GCS_PROJECT_ID= # GCS Project ID to store the file
```

### 2. Setup database

Please make sure you have PostgreSQL 14.8 installed and `golang-migrate` CLI for PostgreSQL is installed.

To Install `golang-migrate` CLI:

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
```

(Optional)
If you want to use dockerized PostgreSQL, you can setup it by using these command.

```shell
docker pull postgres:14.8-alpine3.18 # pull docker image
docker run --name postgres14 -d -p 5432:5432 -e POSTGRES_PASSWORD=your_password -e POSTGRES_USER=your_user postgres:14.8-alpine3.18 # run docker image postgres14
```

If the database is already started, run this to setup the database:

```shell
make setup-db
```

To run the DB Migration individually:

```shell
make migrate-up
```

If you need to rollback the database, you can use this command:

```shell
make migrate-down
```

### 3. Run API

To run API Server

```shell
make run
# or
go run cmd/api/main.go
```

### 4. Run Reconcile Job

To run reconcile processer job

```shell
make run-job
#or
go run cmd/reconcile-job/main.go
```

### 5. Create Docker Container for Deployment

To create docker container for deployment API

```shell
make docker.api
```

To create docker container for deployment Reconcile Job

```shell
make docker.reconcile-job
```

To try running your docker container deployment API

```shell
docker run --rm --name recon-api -p 8080:8080 --env-file .env --platform=linux/amd64 delly/amartha-recon-api:demo-amd64
```

To try running your docker container deployment Reconcile Job

```shell
docker run --rm --name recon-job --env-file .env --platform=linux/amd64 delly/amartha-recon-job:demo-amd64
```

Since file storage implementation used in this project is using local storage, so to ensure that create new request Reconcile and Process Reconcile Job works, you need to implement contract `FileStorageRepository` with storage that can be accessed by the docker container, for example store it in Google Cloud Storage, or something like that.

## Documentation

### Get Reconciliation List

![get reconciliation list image](https://www.planttext.com/api/plantuml/png/XSz13e9030NGVK_nBq0aBaaqI50shemUOCI9rjGPCbF2zHqqId3Zkg__jsLK4xH_21ghEDZMkvQ5ZR9ts7DKxCGFHCegzeyvHHkGhTzYqt61Pdl48imM8dt68wsh0jSKQaJmShZxSwIwGZOB2bRxu7xPb9JmsFw5opp7m7gRD2GTIbHQTqdVhfu0)

Path: `/reconciliations?limit=10&offset=0`<br/>
Method: `GET`<br/>
Query Params:

- limit (integer)
  - min: 1
  - max: 100
- offset (integer)
  - min: 0

Response:

Success: Status Code 200 (OK)

```json
{
    "data": [
        {
            "id": 1,
            "status": "SUCCESS",
            "discrepancy_threshold": 0,
            "system_transaction_csv_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307103607000_1pFvighg/Recon test - system_trx (3).csv",
            "bank_transaction_csv_paths": [
                {
                    "bank_name": "BCA",
                    "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307104925000_kLqrbPOp/Recon test - bca_trx (1).csv"
                },
                {
                    "bank_name": "BRI",
                    "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307105323000_pHNa2RG2/Recon test - bri_trx (1).csv"
                }
            ],
            "start_date": "2024-10-01T00:00:00Z",
            "end_date": "2024-11-28T00:00:00Z"
        }
    ],
    "meta": {
        "limit": 10,
        "offset": 0,
        "total": 1
    }
}
```

### Get Reconciliation Job Request by ID

![get reconciliation job request by id](https://www.planttext.com/api/plantuml/png/ZP7BRW8n34Nt-Gh_069KiM9HKSHZZ-YwcZ86Z8AD74SiAluTPqPJ65LgkuWLA_SUFp9BLglbSuGr6cnm9xoZIBMHiASfHuDLb6i8HXRnJzLxGeNHQwTvkz0KriijZ7LWIUElatn-K7CBlQvu5lCf79pVYi4LVYlei9Z3QC1KjApyKrZ7PpUBmLuoDm7WKST1fSbloAIQ18m9duozgU1AxZkod80IN90RueE__OPygIguaXrxuyFL5XeY6s7yBwzhiHiMH05LFHBlHS_jPai9RxsSCFFe7nlk)

Path: `/reconciliations/:id`<br/>
Method: `GET`<br/>
Params:

- id (integer)

Response:

Success:
Status code 200 (OK)

```json
{
    "data": {
        "id": 1,
        "status": "SUCCESS",
        "system_transaction_csv_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307103607000_1pFvighg/Recon test - system_trx (3).csv",
        "bank_transaction_csv_paths": [
            {
                "bank_name": "BCA",
                "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307104925000_kLqrbPOp/Recon test - bca_trx (1).csv"
            },
            {
                "bank_name": "BRI",
                "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307105323000_pHNa2RG2/Recon test - bri_trx (1).csv"
            }
        ],
        "discrepancy_threshold": 0,
        "error_information": "",
        "result": {
            "total_transaction_processed": 14,
            "total_transaction_matched": 13,
            "total_transaction_unmatched": 1,
            "total_discrepancy_amount": 2321979252,
            "missing_transactions": [
                {
                    "id": "ABC-136",
                    "amount": 2321231231,
                    "type": "CREDIT",
                    "time": "2024-11-25T02:44:21+07:00"
                }
            ],
            "missing_bank_transactions": {
                "BCA": [
                    {
                        "id": "BCA-132",
                        "amount": 123,
                        "type": "CREDIT",
                        "time": "2024-11-23T00:00:00Z"
                    },
                    {
                        "id": "BCA-133",
                        "amount": 42131,
                        "type": "DEBIT",
                        "time": "2024-11-25T00:00:00Z"
                    }
                ],
                "BRI": [
                    {
                        "id": "BRI-131",
                        "amount": 241231,
                        "type": "CREDIT",
                        "time": "2024-11-18T00:00:00Z"
                    },
                    {
                        "id": "BRI-132",
                        "amount": 222222,
                        "type": "CREDIT",
                        "time": "2024-11-18T00:00:00Z"
                    },
                    {
                        "id": "BRI-133",
                        "amount": 242314,
                        "type": "CREDIT",
                        "time": "2024-11-28T00:00:00Z"
                    }
                ]
            }
        },
        "start_date": "2024-10-01T00:00:00Z",
        "end_date": "2024-11-28T00:00:00Z",
        "created_at": "2024-11-23T20:58:27.119625+07:00",
        "updated_at": "2024-11-23T20:58:27.119625+07:00"
    }
}
```

Not Found:
Status Code 404 (Not Found)

```json
{
    "message": "reconciliation job not found"
}
```

### Create Reconciliation Job Request

![create reconciliation job request](https://www.planttext.com/api/plantuml/png/RP1B3i8m34JtFeKlKF5PTe5ALNKBed20q1em2WaaBhq-ATz4OcF9dZTZouKNvQI_QDnGQqsejvwyOAtj020iclugEqyEiyLBMruvn_MgsUB4ZNtBcfMmDHu--iZMhAaHwzIHSlJgJdW84uZ6c4MHYRSgtvRd0ZpRlOUgJFWyQD8xyqEGkoWaeEFLNsm-dU70SafvACXquH_m0000)

Path: `/reconciliations`<br/>
Method: `POST`<br/>
Form Data:

- start_date (date)
- end_date (date)
- system_transaction_file (file)
- discrepancy_threshold (float) - in percentage
  - Min: 0
- bank_names (string) - can be multiple
- bank_transaction_files (file) - can be multiple

Sample CSV file can be found under directory `test/data`

cURL example:

```shell
curl --location 'localhost:8080/reconciliations' \
--form 'start_date="2024-10-01"' \
--form 'end_date="2024-11-28"' \
--form 'system_transaction_file=@"/path/to/system/file.csv"' \
--form 'bank_names="BCA"' \
--form 'bank_transaction_files=@"/path/to/bca/file.csv"' \
--form 'bank_names="BRI"' \
--form 'bank_transaction_files=@"/path/to/bri/file.csv"' \
--form 'discrepancy_threshold="0.1"'
```

Response:

Success:
Status Code 201 (Created)

```json
{
    "data": {
        "id": 1,
        "status": "PENDING",
        "system_transaction_csv_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307103607000_1pFvighg/Recon test - system_trx (3).csv",
        "bank_transaction_csv_paths": [
            {
                "bank_name": "BCA",
                "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307104925000_kLqrbPOp/Recon test - bca_trx (1).csv"
            },
            {
                "bank_name": "BRI",
                "file_path": "/Users/delly/latihan/paystone/amartha/temp_storage/1732370307105323000_pHNa2RG2/Recon test - bri_trx (1).csv"
            }
        ],
        "discrepancy_threshold": 0,
        "error_information": "",
        "result": null,
        "start_date": "2024-10-01T00:00:00Z",
        "end_date": "2024-11-28T00:00:00Z",
        "created_at": "2024-11-23T20:58:27.119625+07:00",
        "updated_at": "2024-11-23T20:58:27.119625+07:00"
    }
}
```

Invalid Params: Status Code 400 (Bad Request)

```json
{
    "message": "bank names and bank transaction files length must be same"
}
```

### Reconciliation Job Process

![reconciliation job process](https://www.planttext.com/api/plantuml/png/ZL513i8m3Blt5Vx0Fh03cc3Ym94VT5q6GwMPsbJmVB9nO1i8k5HAxDY9MoMnKVBLuqYEW-izuS0DzfvlnaWlMdz2fjukSXXRnGRrjiI91DPxn173XPjawYqAHUViKd79CQofrWi2ppl0qkLDYEwz6FA9PbFeE8TMPptpe4K4MNT-4HHPwwvbXyYEKlenCrwSXzRAp1sQfkG4OQJi9X5TeBEQNJk9VClZATOkR6awvQyOb6agVVKlpGC0)
