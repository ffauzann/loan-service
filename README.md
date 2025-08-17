# Loan Service
This service is based on a personal boilerplate contain user usecase such as user authentication and management.

## Guide
### Prerequisite
- Programming language : `go@1.21` or later
- RDMS : `postgresql@14.11` or later
- Caching : `redis-server@7.2.1` or later
- Messaging : `kafka@3.7.1` or later
- MailHog (optional) : `mailhog@1.0.1` or later
- Secret management (optional) : `vault@1.16.2` or later
- Containerization (optional) : `docker@26.1.1` or later
- Orchestrator (optional) : `minikube@1.31.2` or later

### Service Configuration
#### With local config files :
1. Duplicate file `internal/app/auth.config.example.yaml` as `internal/app/auth.config.yaml`
2. Duplicate file `internal/app/loan.config.example.yaml` as `internal/app/loan.config.yaml`
3. Update both `auth.config.yaml` and `loan.config.yaml` files with your own credentials

#### With remote secret management :
1. Setup these environment variables :
    - `VAULT_ADDR` : `https://your-vault-server-host`
    - `VAULT_TOKEN` : `hvs.YOURVAULTTOKEN`
    - `VAULT_MOUNT_PATH` : `path/to/secret/engine/data/`
2. Create or use an existing **KV v2 engine** and make sure its the same engine as `VAULT_MOUNT_PATH` value
3. Create new secret in selected engine named `auth` and `user`
4. Convert both yaml config files into json format and then fill it with your own credentials

### Installation
#### Using `go run` :
1. Complete [Service Configuration](#service-configuration) above
2. Create new database schema under the same name as `database.sql.schema` value in `loan.config.yaml`
3. Run `go mod vendor`
4. Run `go run main.go`
5. Import Postman collection from `docs/`

#### Using Docker :
Assuming you already have your own docker enviroment setup, Follow these steps:
1. Complete [Service Configuration](#service-configuration) above. **NOTE: For local development, this service is not yet support using docker alongside remote secret manager.**
2. Update `docker-compose.yaml` as well if needed
3. Run `docker-compose up`
4. Import Postman collection from `docs/`

#### Using Kubernetes :
All manifests including helm configuration files has been tested in k8s cluster with 1 control-plane and 2 worker nodes.
Assuming you already have your own k8s cluster setup, you can duplicate all directories and files under `k8s/testing/` to ignored git directory under `k8s/local/`. Then configure manifests with your desired setup.

### API Standard Overview
This service is primarily using gRPC API to support microservice architecture communication between backend services and mobile apps integrations. Although, this service uses reverse proxy for each gRPC API to REST API server just in case of other integration is prefer to use REST API ecosystem. Section below is some specs you might find across all APIs :
- Path structure consist this prefix format `/<serviceName>/api/<apiVersion>/<g=grpc,r=rest>/<collection>*`
- Response body for when an error occurred consist these keys :
    - `code (int)` : Original gRPC code that is mapped by default to HTTP response code
    - `message (string)` : Error message
    - `details ([]string)` : If some error details or multiple errors are needed to be informed beside the error message
- Success response body for when the requested data is a list should have these keys :
    - `data` : consisting array of object of particular data and should return empty array instead of `NOT_FOUND` error if no data is available
    - `metadata` consisting :
        - `search (string)` : Search query parameter is given on request
        - `limit (int)` : Size of requested data
        - `page (int)` : Current page number
        - `order_by (string)`: Server-side sorting based on valid columns
        - `order_type (string)`: Ascending / descending sorting
        - `next (bool)` : Whether the list has more data without having the need for excessive data counting
        - `previous (bool)` : Whether the list is not on page 1

## License
[MIT] (https://github.com/ffauzann/loan-service/blob/main/LICENSE) 