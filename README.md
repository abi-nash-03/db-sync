# db-sync

A CLI tool to sync your production MySQL database to your development server. Run it on your dev server — it dumps production and restores it locally on a schedule.

## Prerequisites

- Go 1.21+
- `mysql-client` (`mysqldump` and `mysql` must be in PATH)

## Installation

```bash
git clone https://github.com/yourusername/db-sync.git
cd db-sync
go mod tidy
go build -o db-sync .
```

## Configuration

```bash
cp config.example.yaml config.yaml
```

Edit `config.yaml`:

```yaml
source:
  host: "your-production-ip"
  port: 3306
  user: "readonly_user"
  password: "your-password"
  database: "production_db"

destination:
  host: "127.0.0.1"
  port: 3306
  user: "dev_user"
  password: "dev-password"
  database: "dev_db"

schedule: "0 2 * * 0"   # every Sunday at 2am — leave empty to run once

notify:
  slack_webhook: ""      # optional
```

## Production MySQL Setup

Run this on your production server to create a restricted user:

```sql
CREATE USER 'readonly_user'@'your-dev-server-ip' IDENTIFIED BY 'password';
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON production_db.* TO 'readonly_user'@'your-dev-server-ip';
FLUSH PRIVILEGES;
```

## Usage

```bash
# run once
db-sync -c config.yaml

# run on schedule (blocks until Ctrl+C)
db-sync -c config.yaml

# dry run — see what would happen without making changes
db-sync -c config.yaml --dry-run

# debug logs
db-sync -c config.yaml --debug
```

## Running Tests

```bash
go test ./...
```