#!/bin/bash

# PostgreSQL Database Setup Script
# This script creates the databases needed for the application

echo "Setting up PostgreSQL databases..."

# Create main database
createdb -h localhost -U postgres go_taskmanagement 2>/dev/null || echo "Database go_taskmanagement already exists"

# Create test database
createdb -h localhost -U postgres go_taskmanagement_test 2>/dev/null || echo "Database go_taskmanagement_test already exists"

echo "Databases created successfully!"
echo ""
echo "Main database: go_taskmanagement"
echo "Test database: go_taskmanagement_test"
echo ""
echo "You can now run the application with: go run main.go"
