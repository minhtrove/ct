#!/bin/bash
set -e

APP_NAME="ct-finance"
APP_DIR="/home/ec2-user/app"
SERVICE_NAME="ct-finance"

echo "🚀 Starting deployment..."

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
    mv ${APP_DIR}/${APP_NAME} ${APP_DIR}/${APP_NAME}.old
fi

# Move new binary
echo "📦 Installing new binary..."
mv ${APP_DIR}/${APP_NAME} ${APP_DIR}/${APP_NAME}

# Set permissions
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
