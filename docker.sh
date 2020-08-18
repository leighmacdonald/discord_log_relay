#!/bin/bash
echo "Starting log relay..."
echo "$@"
exec /app/discord_log_relay server -c "${CHANNEL_ID}" -t "${TOKEN}"