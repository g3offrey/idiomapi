#!/bin/bash

# API Example Script - Demonstrates the TODO API
# Make sure the server is running on http://localhost:8080

set -e

API_URL="${API_URL:-http://localhost:8080}"

echo "🚀 TODO API Example Usage"
echo "========================="
echo ""

# Health check
echo "1️⃣  Health Check:"
echo "GET $API_URL/health"
curl -s "$API_URL/health" | jq . || echo "Server not running. Start with: make run"
echo ""
echo ""

# Create a todo
echo "2️⃣  Create a TODO:"
echo "POST $API_URL/api/v1/todos"
TODO_ID=$(curl -s -X POST "$API_URL/api/v1/todos" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Study Go best practices and idioms",
    "completed": false
  }' | jq -r '.id')
echo "Created TODO with ID: $TODO_ID"
echo ""

# Get the todo
echo "3️⃣  Get TODO by ID:"
echo "GET $API_URL/api/v1/todos/$TODO_ID"
curl -s "$API_URL/api/v1/todos/$TODO_ID" | jq .
echo ""
echo ""

# Create more todos
echo "4️⃣  Create more TODOs:"
curl -s -X POST "$API_URL/api/v1/todos" \
  -H "Content-Type: application/json" \
  -d '{"title": "Build REST API", "completed": true}' > /dev/null
curl -s -X POST "$API_URL/api/v1/todos" \
  -H "Content-Type: application/json" \
  -d '{"title": "Write tests", "completed": true}' > /dev/null
curl -s -X POST "$API_URL/api/v1/todos" \
  -H "Content-Type: application/json" \
  -d '{"title": "Deploy to production", "completed": false}' > /dev/null
echo "Created 3 more TODOs"
echo ""

# List all todos
echo "5️⃣  List all TODOs (paginated):"
echo "GET $API_URL/api/v1/todos?page=1&page_size=10"
curl -s "$API_URL/api/v1/todos?page=1&page_size=10" | jq .
echo ""
echo ""

# Filter completed todos
echo "6️⃣  Filter completed TODOs:"
echo "GET $API_URL/api/v1/todos?completed=true"
curl -s "$API_URL/api/v1/todos?completed=true" | jq .
echo ""
echo ""

# Update todo
echo "7️⃣  Update TODO:"
echo "PUT $API_URL/api/v1/todos/$TODO_ID"
curl -s -X PUT "$API_URL/api/v1/todos/$TODO_ID" \
  -H "Content-Type: application/json" \
  -d '{"completed": true}' | jq .
echo ""
echo ""

# Delete todo
echo "8️⃣  Delete TODO:"
echo "DELETE $API_URL/api/v1/todos/$TODO_ID"
curl -s -X DELETE "$API_URL/api/v1/todos/$TODO_ID" -w "\nStatus: %{http_code}\n"
echo ""
echo ""

echo "✅ API Demo Complete!"
