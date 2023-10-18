#!/bin/bash
set -e

# Extract and install MinIO Client (mc)
gzip -dc /opt/bin/mc.gz > /usr/local/bin/mc
chmod +x /usr/local/bin/mc

# Start MinIO server
/usr/bin/docker-entrypoint.sh minio server --console-address ":9101" /data &

# Wait for MinIO server to be up and running
while ! mc alias set gigo-dev http://localhost:9000 gigo-dev gigo-dev > /dev/null 2>&1; do
  echo "Waiting for MinIO to start..."
  sleep 1
done

# Configure MinIO client
mc alias set gigo-dev http://localhost:9000 gigo-dev gigo-dev

# Create the desired bucket
mc mb gigo-dev/gigo-dev

# Allow public read access to the bucket
# mc policy set download gigo-dev/gigo-core

# Keep MinIO server running
wait %1
