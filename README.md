# Anchor Backend

Go service which handles the backend of ANCHOR - Bike Locking Service

PN532 + ESP32 Service Repository [here](https://github.com/eevanwong/anchor-arduino)

## Quickstart Guide

### Prerequisites

- **Go**: Make sure you have Go installed. You can download it from the official [Go website](https://golang.org/dl/).
- **Git**: Version control system. Download it [here](https://git-scm.com/).

### Getting Started

1. **Clone the repository**:

   ```sh
   git clone https://github.com/eevanwong/anchor-backend.git
   cd anchor-backend
   ```

2. **Install dependencies**:

   ```sh
   go mod tidy
   ```

3. **Run the server**:

   ```sh
   go run ./main.go
   ```

4. **Check the server**:
   Open your browser and navigate to `http://localhost:8080`. You should see a message indicating that the server is running.

### Run Service

This assumes docker is already setup on the system. Install the docker vscode extension.

```sh
docker-compose build --no-cache && docker-compose up
```

### Testing

**Connect to the database**

```sh
docker exec -it anchor-backend-dev-db-1 bash
psql -U docker
```
- Then run `\dt` to see the current schema and check changes.

**Sample Endpoint**

```sh
curl -X POST http://localhost:8080/api/lock \
     -H "Content-Type: application/json" \
     -d '{"rack_id":1, "user_name":"John Doe", "user_email":"johndoe@gmail.com", "user_phone":"1234561234"}'
```


