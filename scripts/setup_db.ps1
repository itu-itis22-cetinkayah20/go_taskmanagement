# PostgreSQL Database Setup Script for Windows
# This script creates the databases needed for the application

Write-Host "Setting up PostgreSQL databases..." -ForegroundColor Green

# Create main database
try {
    createdb -h localhost -U postgres go_taskmanagement
    Write-Host "✓ Created database: go_taskmanagement" -ForegroundColor Green
} catch {
    Write-Host "ℹ Database go_taskmanagement already exists" -ForegroundColor Yellow
}

# Create test database
try {
    createdb -h localhost -U postgres go_taskmanagement_test
    Write-Host "✓ Created database: go_taskmanagement_test" -ForegroundColor Green
} catch {
    Write-Host "ℹ Database go_taskmanagement_test already exists" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Databases setup completed!" -ForegroundColor Green
Write-Host "Main database: go_taskmanagement" -ForegroundColor Cyan
Write-Host "Test database: go_taskmanagement_test" -ForegroundColor Cyan
Write-Host ""
Write-Host "You can now run the application with: go run main.go" -ForegroundColor Yellow
