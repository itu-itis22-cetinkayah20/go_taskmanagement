# PowerShell script to create PostgreSQL databases
# This script creates the required databases for the Go Task Management application

Write-Host "🗄️ Creating PostgreSQL databases for Go Task Management..." -ForegroundColor Green

# Set environment variables for PostgreSQL
$env:PGPASSWORD = "1234"
$PSQL_PATH = "C:\Program Files\PostgreSQL\17\bin\psql.exe"

# Check if psql exists
if (-not (Test-Path $PSQL_PATH)) {
    Write-Host "❌ psql.exe not found at $PSQL_PATH" -ForegroundColor Red
    Write-Host "Please ensure PostgreSQL is installed and update the PSQL_PATH variable" -ForegroundColor Yellow
    exit 1
}

Write-Host "📋 Creating databases..." -ForegroundColor Cyan

try {
    # Create main database
    Write-Host "Creating database: go_taskmanagement" -ForegroundColor Yellow
    & $PSQL_PATH -U postgres -h localhost -p 5432 -c "CREATE DATABASE go_taskmanagement;" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database 'go_taskmanagement' created successfully" -ForegroundColor Green
    } else {
        Write-Host "⚠️ Database 'go_taskmanagement' may already exist" -ForegroundColor Yellow
    }

    # Create test database
    Write-Host "Creating database: go_taskmanagement_test" -ForegroundColor Yellow
    & $PSQL_PATH -U postgres -h localhost -p 5432 -c "CREATE DATABASE go_taskmanagement_test;" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database 'go_taskmanagement_test' created successfully" -ForegroundColor Green
    } else {
        Write-Host "⚠️ Database 'go_taskmanagement_test' may already exist" -ForegroundColor Yellow
    }

    Write-Host "🎯 Database creation completed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "📊 Database Information:" -ForegroundColor Cyan
    Write-Host "  - Main Database: go_taskmanagement" -ForegroundColor White
    Write-Host "  - Test Database: go_taskmanagement_test" -ForegroundColor White
    Write-Host "  - Host: localhost" -ForegroundColor White
    Write-Host "  - Port: 5432" -ForegroundColor White
    Write-Host "  - Username: postgres" -ForegroundColor White
    Write-Host "  - Password: 1234" -ForegroundColor White
    Write-Host ""
    Write-Host "🚀 You can now run your application with database support!" -ForegroundColor Green

} catch {
    Write-Host "❌ Error creating databases: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
} finally {
    # Clean up environment variable
    Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue
}
