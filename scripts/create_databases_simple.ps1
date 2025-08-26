# PowerShell script to create PostgreSQL databases
Write-Host "Creating PostgreSQL databases..." -ForegroundColor Green

$env:PGPASSWORD = "1234"
$PSQL_PATH = "C:\Program Files\PostgreSQL\17\bin\psql.exe"

# Create main database
Write-Host "Creating database: go_taskmanagement" -ForegroundColor Yellow
& $PSQL_PATH -U postgres -h localhost -p 5432 -c "CREATE DATABASE go_taskmanagement;"

# Create test database  
Write-Host "Creating database: go_taskmanagement_test" -ForegroundColor Yellow
& $PSQL_PATH -U postgres -h localhost -p 5432 -c "CREATE DATABASE go_taskmanagement_test;"

Write-Host "Database creation completed!" -ForegroundColor Green

# Clean up
Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue
