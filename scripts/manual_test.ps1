# Manual API Tests Script for Task Management
# Usage: 
# 1) Start the server: go run main.go
# 2) In a new PowerShell: .\scripts\manual_test.ps1

$base = 'http://localhost:8080'

Write-Host '--- Register new user ---'
$response = Invoke-RestMethod -Method Post -Uri "$base/register" `
    -Body (@{ username='alice'; password='secret' } | ConvertTo-Json) `
    -ContentType 'application/json'
$response | ConvertTo-Json -Depth 5 | Write-Host

Write-Host '--- Login user ---'
$login = Invoke-RestMethod -Method Post -Uri "$base/login" `
    -Body (@{ username='alice'; password='secret' } | ConvertTo-Json) `
    -ContentType 'application/json'
$token = $login.token
Write-Host "Token: $token"

Write-Host '--- Create a new task ---'
$task = Invoke-RestMethod -Method Post -Uri "$base/tasks" `
    -Headers @{ Authorization = "Bearer $token" } `
    -Body (@{ title='Test Görev'; details='Detay' } | ConvertTo-Json) `
    -ContentType 'application/json'
$task | ConvertTo-Json | Write-Host

Write-Host '--- List my tasks ---'
$tasks = Invoke-RestMethod -Method Get -Uri "$base/tasks" `
    -Headers @{ Authorization = "Bearer $token" }
$tasks | ConvertTo-Json -Depth 5 | Write-Host

Write-Host '--- List public tasks ---'
$public = Invoke-RestMethod -Method Get -Uri "$base/tasks/public"
$public | ConvertTo-Json -Depth 5 | Write-Host

Write-Host 'Manual API testing completed.'
