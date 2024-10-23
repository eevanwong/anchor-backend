# Anchor Backend

Go service which handles the backend of ANCHOR - Bike Locking Service

## Quickstart Guide

### Prerequisites

- **Go**: Make sure you have Go installed. You can download it from the official [Go website](https://golang.org/dl/).
- **Git**: Version control system. Download it [here](https://git-scm.com/).

### Getting Started

1. **Clone the repository**:

   ```sh
   git clone https://github.com/eevanwong/anchor-backend.git
   cd your_project_name
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

### Setup Database

This assumes docker is already setup on the system. Install the docker vscode extension.

1. **Setup Docker Container**

```sh
docker-compose up
```

2. **Run Migrations and Seed Database**

```sh
go run ./database/db_setup.go
```

3. **(Optional) Check the database**
   Via the Docker tab on vscode, attach the shell of the postgres dev_db to the CLI and run:

```sh
psql -U docker
```

Then run `\dt` to see the current schema and check changes.
