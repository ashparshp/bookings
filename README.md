## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/ashparshp/bookings.git
cd bookings
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
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

Update `database.yml` with your database credentials:

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

Make sure you have Soda CLI installed:

```bash
# Install Soda CLI
go install github.com/gobuffalo/pop/v6/soda@latest

# Run migrations for development
soda migrate

# Run migrations for production environment
soda migrate -e production
```

### 4. Configuration Options

The application supports the following command-line flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-dbhost` | Database host | localhost |
| `-dbport` | Database port | 5432 |
| `-dbname` | Database name | (required) |
| `-dbuser` | Database username | (required) |
| `-dbpassword` | Database password | (required) |
| `-dbssl` | SSL mode | disable |
| `-production` | Production mode | true |
| `-cache` | Template caching | true |
| `-mailhost` | SMTP server host | localhost |
| `-mailport` | SMTP server port | 1025 |
| `-mailusername` | SMTP username | "" |
| `-mailpassword` | SMTP password | "" |
| `-mailencryption` | Encryption (none/tls/ssl) | none |
| `-mailfrom` | Sender email address | noreply@bookings.com |
| `-mailfromname` | Sender name | "Bookings" |

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
go run cmd/web/*.go \
  -dbname=bookings \
  -dbuser=your_username \
  -dbpassword=your_password \
  -production=false \
  -cache=false
```

#### Production Mode

```bash
# Build the application
go build -o bookings cmd/web/*.go

# Run with production settings
./bookings \
  -dbname=your_dbname \
  -dbuser=your_dbuser \
  -dbpassword=your_password \
  -dbhost=your_dbhost \
  -dbport=5432 \
  -dbssl=require \
  -cache=true \
  -production=true \
  -mailhost=smtp.example.com \
  -mailport=587 \
  -mailusername=your_email@example.com \
  -mailpassword=your_email_password \
  -mailencryption=starttls \
  -mailfrom=noreply@bookings.com \
  -mailfromname="Bookings System"
```

### 6. Email Configuration

For development, you can use MailHog for local email testing:

```bash
# Install MailHog
go install github.com/mailhog/MailHog@latest

# Run MailHog
MailHog
```

Then access the web UI at http://localhost:8025

For production, configure your actual SMTP settings as command-line parameters.