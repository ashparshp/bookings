## Bookings and Reservations

A modern hotel booking management system built with Go, featuring a web interface for guests to make reservations and an admin panel for hotel staff to manage bookings.

🔗 **GitHub Repository**: [https://github.com/ashparshp/bookings](https://github.com/ashparshp/bookings)

- Built in Go version 1.15
- Uses the [chi router](github.com/go-chi/chi)
- Uses [alex edwards scs session management](github.com/alexedwards/scs)
- Uses [nosurf](github.com/justinas/nosurf)

## Features

### Guest Features
- **Room Availability Search**: Check room availability for specific dates
- **Online Reservations**: Make reservations for available rooms
- **Contact Form**: Send inquiries to hotel staff
- **Responsive Design**: Mobile-friendly interface with Bootstrap
- **Email Confirmations**: Automated booking confirmations

### Admin Features
- **Dashboard**: Overview of hotel operations and statistics
- **Reservation Management**: View, edit, and process reservations
- **User Management**: Handle user accounts and authentication
- **Real-time Notifications**: Session-based alerts and messages
- **CSRF Protection**: Secure forms with token validation

### Technical Features
- **Session Management**: Secure user sessions with SCS
- **Database Integration**: PostgreSQL with Buffalo Pop migrations
- **Email Service**: Background email processing with channels
- **Template Caching**: Optimized Go HTML template rendering
- **Form Validation**: Comprehensive input validation
- **Repository Pattern**: Clean database abstraction layer
- **Testing Suite**: Unit tests with mock repositories

## 🛠 Technology Stack

- **Backend**: Go 1.19+
- **Database**: PostgreSQL 12+
- **Session Store**: SCS (Alexedwards Session)
- **Templating**: Go HTML Templates
- **Email**: Go Simple Mail with goroutines
- **Router**: Chi Router with middleware
- **Frontend**: HTML, CSS, JavaScript, Bootstrap
- **Migrations**: Buffalo Pop (Soda)
- **Deployment**: Render.com ready

## 📁 Project Structure

```
bookings/
├── cmd/web/                    # Application entry point
│   ├── main.go                # Main application file with CLI flags
│   └── send-mail.go           # Background email service
├── internal/
│   ├── config/                # Application configuration
│   ├── driver/                # Database driver and connection
│   ├── forms/                 # Form validation logic
│   ├── handlers/              # HTTP handlers and business logic
│   │   ├── handlers.go        # Main handlers (HomePage, AboutPage, etc.)
│   │   └── setup_test.go      # Test configuration
│   ├── helpers/               # Utility functions
│   ├── models/                # Data models and structs
│   ├── render/                # Template rendering engine
│   │   └── setup_test.go      # Render test setup
│   └── repository/            # Database operations
│       └── dbrepo/            # PostgreSQL and test repositories
├── templates/                 # HTML templates
├── static/                    # Static assets
│   ├── admin/                 # Admin panel assets
│   ├── css/                   # Stylesheets
│   ├── images/                # Images and icons
│   └── js/                    # JavaScript files
├── email-templates/           # Email HTML templates
│   └── basic.html             # Basic email template
├── migrations/                # Database migrations (Soda/Pop)
│   ├── *_create_user_table.*
│   ├── *_create_rooms_table.*
│   └── *_create_restrictions_table.*
├── pkg/                       # Public packages
├── database.yml               # Soda database configuration

## 🔧 Prerequisites

- **Go** 1.19 or later
- **PostgreSQL** 12 or later
- **Soda CLI** (Buffalo Pop) for migrations
- **SMTP Server** (for email functionality)

### Installing Soda CLI

```bash
# Install Soda (Buffalo Pop) for database migrations
go install github.com/gobuffalo/pop/v6/soda@latest

# Verify installation
soda --version
```

## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/ashparshp/bookings.git
cd bookings
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

#### Create PostgreSQL Database

```bash
# Create database
createdb bookings

# For production environment
createdb bookings_production
```

#### Configure Database

Update [`database.yml`](database.yml) with your database credentials:

```yaml
development:
  dialect: postgres
  database: bookings
  user: your_username
  password: your_password
  host: localhost
  port: 5432
  pool: 5

production:
  dialect: postgres
  database: bookings_production
  user: your_production_user
  password: your_production_password
  host: your_production_host
  port: 5432
  pool: 25
```

#### Run Database Migrations

```bash
# Run migrations for development
soda migrate

# Run migrations for production environment
soda migrate -e production

# Check migration status
soda migrate status

# Rollback last migration (if needed)
soda migrate down
```

### 4. Environment Configuration

The application accepts the following command-line flags:

```bash
-dbhost       # Database host (default: localhost)
-dbport       # Database port (default: 5432)
-dbname       # Database name (required)
-dbuser       # Database username (required)
-dbpassword   # Database password (required)
-dbssl        # SSL mode: disable, prefer, require (default: disable)
-production   # Run in production mode (default: true)
-cache        # Use template caching (default: true)
```

### 5. Build and Run

#### Development Mode

Using the provided script:
```bash
# Make script executable
chmod +x run.sh

# Run development server
./run.sh
```

Or manually:
```bash
go run cmd/web/*.go -dbname=bookings -dbuser=your_username -dbpassword=your_password -production=false -cache=false
```

#### Production Mode

```bash
# Build the application
go build -o bookings cmd/web/*.go

# Run with production settings
./bookings -dbname=bookings_production -dbuser=your_username -dbpassword=your_password -production=true -cache=true
```

## 📧 Email Configuration

The application uses SMTP for sending emails. Configure your SMTP settings in [`cmd/web/send-mail.go`](cmd/web/send-mail.go):

```go
server.Host = "smtp.gmail.com"
server.Port = 587
server.Username = "your_email@gmail.com"
server.Password = "your_app_password"
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Test specific package
go test ./internal/handlers/
```

### Test Configuration

The application includes comprehensive test setup in:
- [`internal/handlers/setup_test.go`](internal/handlers/setup_test.go)
- [`internal/render/setup_test.go`](internal/render/setup_test.go)

Tests use mock repositories and in-memory sessions for isolated testing.

## 🚀 Deployment

### Render.com Deployment

The application includes a [`render.yaml`](render.yaml) file for easy deployment:

```yaml
services:
  - type: web
    name: bookings-app
    env: go
    buildCommand: go build -o bookings cmd/web/*.go
    startCommand: ./bookings -dbname=bookings -dbuser=ashparsh -cache=false -production=true
    envVars:
      - key: PORT
        value: 8080
```

### Manual Deployment

1. **Build for target platform**:
```bash
# For Linux (most cloud providers)
GOOS=linux GOARCH=amd64 go build -o bookings cmd/web/*.go
```

2. **Set up production database**:
```bash
soda migrate -e production
```

3. **Run with production settings**:
```bash
./bookings -dbname=your_prod_db -dbuser=your_prod_user -dbpassword=your_prod_password -production=true
```

## 🛣 API Endpoints

### Public Routes
- `GET /` - Home page
- `GET /about` - About page  
- `GET /contact` - Contact page
- `POST /contact` - Submit contact form
- `GET /generals-quarters` - Generals Quarters room page
- `GET /majors-suite` - Majors Suite room page
- `GET /make-reservation` - Reservation form
- `POST /make-reservation` - Submit reservation
- `GET /reservation-summary` - Booking confirmation

### Authentication Routes
- `GET /user/login` - Login page
- `POST /user/login` - Login submission
- `GET /user/logout` - Logout

### Admin Routes (Protected)
- `GET /admin/dashboard` - Admin dashboard
- `GET /admin/reservations-new` - New reservations
- `GET /admin/reservations-all` - All reservations
- `GET /admin/reservations/{src}/{id}/show` - Show reservation details
- `POST /admin/reservations/{src}/{id}` - Update reservation
- `POST /admin/process-reservation/{src}/{id}/{action}` - Process reservation

### AJAX Endpoints
- `POST /search-availability-json` - Check room availability (JSON)
- `POST /choose-room/{id}` - Select specific room

## 💾 Database Schema

### Main Tables (via Soda migrations)

- **users**: Admin users and authentication
- **rooms**: Hotel room information and pricing
- **reservations**: Booking details and guest information
- **restrictions**: Room availability restrictions
- **room_restrictions**: Junction table for room-specific restrictions

### Key Models

```go
type Reservation struct {
    ID        int       `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
    RoomID    int       `json:"room_id"`
    Room      Room      `json:"room"`
}
```

## 🔧 Database Management

### Soda Commands

```bash
# Create a new migration
soda generate fizz create_new_table

# Check migration status
soda migrate status

# Run migrations
soda migrate
soda migrate -e production

# Rollback migrations
soda migrate down
soda migrate down -e production

# Reset database (careful!)
soda reset
```

### Migration Files

Migration files are located in [`migrations/`](migrations/) directory:
- `*.up.fizz` - Forward migrations
- `*.down.fizz` - Rollback migrations

## 🔐 Security Features

- **CSRF Protection**: Implemented on all forms using nosurf middleware
- **Session Security**: Secure cookie configuration with SameSite protection
- **Input Validation**: Comprehensive form validation in [`internal/forms/`](internal/forms/)
- **SQL Injection Prevention**: Parameterized queries through repository pattern
- **XSS Protection**: Template escaping and validation

## 🎯 Key Features Implementation

### Repository Pattern
```go
type DatabaseRepo interface {
    AllUsers() bool
    InsertReservation(res models.Reservation) (int, error)
}
```

### Session Management
```go
session = scs.New()
session.Lifetime = 24 * time.Hour
session.Cookie.Persist = true
session.Cookie.Secure = app.InProduction
```

### Background Email Processing
```go
mailChan := make(chan models.MailData)
// Background goroutine processes emails
go listenForMail()
```
---