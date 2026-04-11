# db-sync

A CLI tool to sync your production MySQL database to your development server. Run it on your dev server — it dumps production and restores it locally on a schedule.

## Prerequisites

- Go 1.21+ (only needed if building from source)
- `mysql-client` (`mysqldump` and `mysql` must be in PATH)

---

## Installation

**Option A — Download binary (recommended)**
```bash
wget https://github.com/yourusername/db-sync/releases/download/v1.0.0/db-sync-linux-amd64
chmod +x db-sync-linux-amd64
sudo mv db-sync-linux-amd64 /usr/local/bin/db-sync
```

**Option B — Build from source**
```bash
git clone https://github.com/yourusername/db-sync.git
cd db-sync
go mod tidy
go build -o db-sync .
sudo mv db-sync /usr/local/bin/
```

---

## Configuration

```bash
sudo mkdir -p /etc/db-sync
sudo cp config.example.yaml /etc/db-sync/config.yaml
sudo nano /etc/db-sync/config.yaml
```

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

> ⚠️ Never commit `config.yaml` — it contains credentials. Commit `config.example.yaml` instead.

---

## Production MySQL Setup

Run on your **production server** to create a restricted readonly user:

```sql
CREATE USER 'readonly_user'@'your-dev-server-ip' IDENTIFIED BY 'password';
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON production_db.* TO 'readonly_user'@'your-dev-server-ip';
FLUSH PRIVILEGES;
```

---

## Usage

```bash
# run once
db-sync -c /etc/db-sync/config.yaml

# dry run — see what would happen without making changes
db-sync -c /etc/db-sync/config.yaml --dry-run

# debug logs
db-sync -c /etc/db-sync/config.yaml --debug
```

---

## Running as a Service (Recommended)

Set up `db-sync` as a systemd service so it runs automatically on schedule and survives server reboots.

**Create the service file:**
```bash
sudo nano /etc/systemd/system/db-sync.service
```

```ini
[Unit]
Description=db-sync — MySQL production to dev sync
After=network.target mysql.service

[Service]
Type=simple
ExecStart=/usr/local/bin/db-sync -c /etc/db-sync/config.yaml
Restart=on-failure
RestartSec=10
User=ubuntu
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Enable and start:**
```bash
sudo systemctl daemon-reload
sudo systemctl start db-sync
sudo systemctl enable db-sync    # auto-start on reboot
```

**Useful commands:**
```bash
sudo systemctl status db-sync          # check status
sudo journalctl -u db-sync -f          # live logs
sudo journalctl -u db-sync -n 50       # last 50 log lines
sudo systemctl restart db-sync         # restart
sudo systemctl stop db-sync            # stop
```

---

## Running Tests

```bash
go test ./...
```
