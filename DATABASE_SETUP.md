# PostgreSQL Database Setup for Go Task Management

## Prerequisites
1. PostgreSQL must be installed on your system
2. PostgreSQL service should be running

## Setup Steps

### 1. Using psql (Command Line)
```bash
# Connect to PostgreSQL as superuser
psql -U postgres

# Create databases
CREATE DATABASE go_taskmanagement;
CREATE DATABASE go_taskmanagement_test;

# Exit psql
\q
```

### 2. Using pgAdmin (GUI)
1. Open pgAdmin
2. Connect to your PostgreSQL server
3. Right-click on "Databases" → Create → Database
4. Create `go_taskmanagement`
5. Create `go_taskmanagement_test`

### 3. Environment Configuration
Make sure your `.env` file has the correct PostgreSQL credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_taskmanagement
DB_SSLMODE=disable

TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=go_taskmanagement_test
TEST_DB_SSLMODE=disable
```

### 4. Test Connection
```bash
# Test main database
psql -h localhost -U postgres -d go_taskmanagement -c "SELECT 1;"

# Test test database
psql -h localhost -U postgres -d go_taskmanagement_test -c "SELECT 1;"
```

## Running the Application

### Start the Server
```bash
go run main.go
```

### Run Tests
```bash
# All tests
go test ./test/contract -v

# Specific auth flow test
go test ./test/contract -run TestRealAuthenticationFlow -v
```

## Troubleshooting

### Authentication Failed
If you see "password authentication failed", update your `.env` file with the correct PostgreSQL password.

### Database Not Found
Make sure the databases are created using the commands above.

### Connection Refused
- Check if PostgreSQL service is running
- Verify the port (default: 5432)
- Check if PostgreSQL is listening on localhost
