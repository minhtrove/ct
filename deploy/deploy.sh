#!/bin/bash
set -e

APP_NAME="ct-finance"
APP_DIR="$(cd "$(dirname "$0")" && pwd)"
SERVICE_NAME="ct-finance"

echo "🚀 Starting deployment..."
echo "📁 APP_DIR: ${APP_DIR}"
echo "📁 Working directory: $(pwd)"
echo "📁 Files present:"
ls -lah "${APP_DIR}/"

# Stop the service if it's running
if sudo systemctl is-active --quiet ${SERVICE_NAME}; then
    echo "⏸️  Stopping ${SERVICE_NAME} service..."
    sudo systemctl stop ${SERVICE_NAME}
fi

# Backup old binary
if [ -f "${APP_DIR}/${APP_NAME}.old" ]; then
    rm -f ${APP_DIR}/${APP_NAME}.old
fi

if [ -f "${APP_DIR}/${APP_NAME}" ]; then
    echo "📦 Backing up old binary..."
    mv ${APP_DIR}/${APP_NAME} ${APP_DIR}/${APP_NAME}.old
fi

# New binary is already in place from tar extraction
# Just set permissions
echo "📦 Setting permissions on new binary..."
chmod +x ${APP_DIR}/${APP_NAME}

# Install/update systemd service
echo "⚙️  Setting up systemd service..."
sudo cp ${APP_DIR}/app.service /etc/systemd/system/${SERVICE_NAME}.service
sudo systemctl daemon-reload
sudo systemctl enable ${SERVICE_NAME}

# Start the service
echo "▶️  Starting ${SERVICE_NAME} service..."
sudo systemctl start ${SERVICE_NAME}

# Check status
sleep 2
if sudo systemctl is-active --quiet ${SERVICE_NAME}; then
    echo "✅ Deployment successful! Service is running."
    sudo systemctl status ${SERVICE_NAME} --no-pager
else
    echo "❌ Service failed to start. Rolling back..."
    if [ -f "${APP_DIR}/${APP_NAME}.old" ]; then
        mv ${APP_DIR}/${APP_NAME}.old ${APP_DIR}/${APP_NAME}
        sudo systemctl start ${SERVICE_NAME}
    fi
    exit 1
fi

# Show recent logs
echo ""
echo "📋 Recent logs:"
sudo journalctl -u ${SERVICE_NAME} -n 20 --no-pager

echo ""
echo "✨ Deployment complete!"
