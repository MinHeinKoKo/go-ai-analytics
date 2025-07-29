#!/bin/bash

echo "🌱 AI Analytics Database Seeder"
echo "==============================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating default .env file..."
    cat > .env << EOF
PORT=8080
HOST=localhost
MONGODB_URI=mongodb://localhost:27017
MONGODB_DBNAME=ai-analytics
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
EOF
    echo "✅ Created .env file with default values"
fi

# Check if MongoDB is running
echo "🔍 Checking MongoDB connection..."
if ! mongosh --eval "db.runCommand('ping')" --quiet > /dev/null 2>&1; then
    echo "❌ MongoDB is not running or not accessible"
    echo "💡 Please start MongoDB first:"
    echo "   - Using Docker: make docker-run"
    echo "   - Or start your local MongoDB instance"
    exit 1
fi

echo "✅ MongoDB is running"

# Run the seeder
echo "🚀 Running database seeder..."
go run seed/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Database seeding completed successfully!"
    echo ""
    echo "📊 You can now:"
    echo "   1. Start the API server: make run"
    echo "   2. Start the frontend: cd web && npm run dev"
    echo "   3. Visit http://localhost:3000 to see the data"
else
    echo "❌ Database seeding failed"
    exit 1
fi
