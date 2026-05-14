#!/usr/bin/env bash
# Deployment script for Ubuntu 22.04 VPS
# Run as root once: bash scripts/deploy.sh
set -euo pipefail

APP_DIR=/opt/vcita-agent
APP_USER=vcita-agent
SERVICE_NAME=vcita-agent

# Create dedicated system user (no login shell)
id "$APP_USER" &>/dev/null || useradd -r -s /bin/false -d "$APP_DIR" "$APP_USER"

# Create directory structure with restricted permissions
mkdir -p "$APP_DIR"/{data,logs}
chown -R "$APP_USER":"$APP_USER" "$APP_DIR"
chmod 750 "$APP_DIR"
chmod 700 "$APP_DIR/data" "$APP_DIR/logs"

# Install systemd service
cat > /etc/systemd/system/"$SERVICE_NAME".service << SERVICE
[Unit]
Description=vcita AI Patient Assistant
After=network.target

[Service]
Type=simple
User=$APP_USER
WorkingDirectory=$APP_DIR
EnvironmentFile=$APP_DIR/.env
ExecStart=$APP_DIR/vcita-agent
Restart=on-failure
RestartSec=5s
# Harden the service
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=$APP_DIR/data $APP_DIR/logs

[Install]
WantedBy=multi-user.target
SERVICE

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
echo "Deploy script complete. Copy binary and .env to $APP_DIR, then: systemctl start $SERVICE_NAME"
