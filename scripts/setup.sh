#!/bin/bash

# CertiTrack Setup Script
# This script sets up the development environment for CertiTrack

set -e

echo "🚀 Setting up CertiTrack development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_requirements() {
    print_status "Checking requirements..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ first."
        exit 1
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js 18+ first."
        exit 1
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm first."
        exit 1
    fi
    
    print_success "All requirements are satisfied!"
}

# Create environment file
setup_environment() {
    print_status "Setting up environment variables..."
    
    if [ ! -f .env ]; then
        cp .env.example .env
        print_success "Created .env file from .env.example"
    else
        print_warning ".env file already exists, skipping..."
    fi
}

# Setup backend
setup_backend() {
    print_status "Setting up backend..."
    
    cd backend
    
    # Download Go dependencies
    print_status "Downloading Go dependencies..."
    go mod tidy
    
    # Create necessary directories
    mkdir -p storage/documents
    mkdir -p storage/temp
    mkdir -p logs
    
    print_success "Backend setup completed!"
    cd ..
}

# Setup frontend
setup_frontend() {
    print_status "Setting up frontend..."
    
    cd frontend
    
    # Install npm dependencies
    print_status "Installing npm dependencies..."
    npm install
    
    print_success "Frontend setup completed!"
    cd ..
}

# Setup database
setup_database() {
    print_status "Setting up database..."
    
    # Start Docker services
    print_status "Starting Docker services..."
    docker-compose up -d postgres redis mailhog
    
    # Wait for services to be ready
    print_status "Waiting for services to start..."
    sleep 15
    
    # Check if services are running
    if docker-compose ps | grep -q "Up"; then
        print_success "Docker services are running!"
    else
        print_error "Failed to start Docker services"
        exit 1
    fi
}

# Create initial migration files
create_migrations() {
    print_status "Creating database migration files..."
    
    mkdir -p backend/migrations
    
    # Create initial migration
    cat > backend/migrations/001_initial_schema.up.sql << 'EOF'
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE
);

-- People table
CREATE TABLE people (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id VARCHAR(50) UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    department VARCHAR(100),
    position VARCHAR(100),
    hire_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Equipment table
CREATE TABLE equipment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_number VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(100),
    location VARCHAR(200),
    purchase_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Certification types table
CREATE TABLE certification_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    default_validity_period INTEGER,
    requires_renewal BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_people_employee_id ON people(employee_id);
CREATE INDEX idx_equipment_asset_number ON equipment(asset_number);
EOF

    cat > backend/migrations/001_initial_schema.down.sql << 'EOF'
DROP TABLE IF EXISTS certification_types;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
EOF

    print_success "Migration files created!"
}

# Main setup function
main() {
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    CertiTrack Setup                         ║"
    echo "║              Certification Management System                 ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""
    
    check_requirements
    setup_environment
    setup_backend
    setup_frontend
    setup_database
    create_migrations
    
    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    Setup Complete!                          ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""
    print_success "CertiTrack development environment is ready!"
    echo ""
    echo "Next steps:"
    echo "1. Run 'make dev' to start the development servers"
    echo "2. Backend API: http://localhost:8080"
    echo "3. Frontend: http://localhost:3000"
    echo "4. Mailhog (email testing): http://localhost:8025"
    echo ""
    echo "For more information, check the README.md file."
}

# Run main function
main "$@"