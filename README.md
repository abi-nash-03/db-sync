# db-sync

A simple CLI tool to sync your production MySQL database to your development server. Run it on your dev server to pull a fresh dump from production and restore it locally — keeping your dev environment in sync without manual effort.

---

## How It Works

```
Dev Server
──────────────────────────────────────
db-sync runs here
    │
    ├── connects to prod MySQL
    │   runs mysqldump → saves dump locally
    │        /tmp/dump_2026-03-21_14-30-00.sql
    │
    └── pipes dump into local dev MySQL
        cleans up dump file
```

---

## Prerequisites

Make sure these are installed on your **dev server** before proceeding:

| Tool | Purpose | Install |
|------|---------|---------|
| Go 1.21+ | Build the tool | https://go.dev/dl |
| mysqldump | Create the dump | `sudo apt install mysql-client` |
| mysql | Restore the dump | `sudo apt install mysql-client` |
| git | Clone the repo | `sudo apt install git` |

Verify everything is installed:

```bash
go version
mysqldump --version
mysql --version
git --version
```

---

## Installation

### Step 1 — Clone the repository

```bash
git clone https://github.com/yourusername/db-sync.git
cd db-sync
```

### Step 2 — Install dependencies

```bash
go mod tidy
```

### Step 3 — Build the binary

```bash
go build -o db-sync .
```

This creates a `db-sync` executable in the current directory.

To make it available system-wide:

```bash
sudo mv db-sync /usr/local/bin/
```

Verify the installation:

```bash
db-sync --version
```

---

## Configuration

Create a `config.yaml` file in the same directory where you run `db-sync`. A template is provided in the repo as `config.example.yaml`.

```bash
cp config.example.yaml config.yaml
```

Edit `config.yaml` with your credentials:

```yaml
source:
  host: "your-production-server-ip"
  port: 3306
  user: "readonly_user"
  password: "your-production-password"
  database: "your_production_db"

destination:
  host: "127.0.0.1"
  port: 3306
  user: "dev_user"
  password: "your-dev-password"
  database: "your_dev_db"
```

> ⚠️ **Never commit `config.yaml` to git.** It contains your database credentials.
> The `.gitignore` already excludes it. Always use `config.example.yaml` as the template.

---

## Setting Up the Production MySQL User

On your **production server**, create a restricted readonly user that only allows connections from your dev server:

```sql
CREATE USER 'readonly_user'@'your-dev-server-ip' IDENTIFIED BY 'strong-password';
GRANT SELECT, SHOW VIEW, EVENT, TRIGGER ON your_production_db.* TO 'readonly_user'@'your-dev-server-ip';
FLUSH PRIVILEGES;
```

> This ensures the tool can only **read** from production — it can never modify your live data.

---

## Usage

### Run a sync

```bash
db-sync --config config.yaml
```

Or using the shorthand:

```bash
db-sync -c config.yaml
```

**Expected output:**

```
✓ Dump created:     /tmp/dump_2026-03-21_14-30-00.sql
✓ Restore complete: your_dev_db is up to date
```

### Dry run — see what would happen without making any changes

```bash
db-sync --dry-run --config config.yaml
```

### Check the version

```bash
db-sync --version
```

### Help

```bash
db-sync --help
```

---

## Build with Version Injection

When building a release, inject the version from your git tag:

```bash
VERSION=$(git describe --tags --always)
go build -ldflags "-X main.version=${VERSION}" -o db-sync .
```

---

## Running Tests

Run all tests across the project:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test -v ./config/...
go test -v ./dumper/...
go test -v ./restore/...
```

---

## Project Structure

```
db-sync/
├── main.go               # entry point
├── go.mod                # module definition
├── go.sum                # dependency checksums
├── config.example.yaml   # template config (safe to commit)
├── config.yaml           # your actual config (never commit this)
├── cmd/
│   └── root.go           # CLI commands and flags
├── config/
│   ├── config.go         # config loading and validation
│   └── config_test.go
├── dumper/
│   ├── dumper.go         # mysqldump execution
│   └── dumper_test.go
└── restore/
    ├── restore.go        # mysql restore execution
    └── restore_test.go
```

---

## Troubleshooting

**`mysqldump not found in PATH`**
```bash
sudo apt install mysql-client
```

**`Access denied for user`**
Check that your production MySQL user has the correct grants and is allowed from your dev server IP. Re-run the `GRANT` statement above.

**`dump file not found`**
The dump step likely failed silently. Run with `-v` verbose logging and check `/tmp/` for partial dump files.

**`connection refused` on port 3306**
Your production server's firewall may be blocking port 3306 from your dev server IP. Check your firewall rules:
```bash
sudo ufw allow from your-dev-server-ip to any port 3306
```

---

## Contributing

Pull requests are welcome. For major changes, open an issue first to discuss what you'd like to change. Make sure all tests pass before submitting:

```bash
go fmt ./...
go vet ./...
go test ./...
```

---

## License

MIT