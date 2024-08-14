#!/bin/bash

# Name of the systemd service to monitor
service_name="scrapeops-go-proxy-worker.service"

# Maximum memory usage in bytes (12 gigabytes)
max_memory_b=12900000000


# Check if the service is running
if systemctl is-active --quiet $service_name; then
    # Get the memory usage of the service in kilobytes
    memory_usage_b=$(systemctl show -p MemoryCurrent $service_name --value)

    # Check if memory usage exceeds the limit
    if [ -n "$memory_usage_b" ] && [ "$memory_usage_b" -gt "$max_memory_b" ]; then
        echo "Memory usage of $service_name exceeds 12 gigabytes. Restarting the service..."
        systemctl restart $service_name
        if [ $? -eq 0 ]; then
            echo "Service successfully restarted."
        else
            echo "Failed to restart the service."
        fi
    else
        echo "Memory usage of $service_name is within the limit. $memory_usage_b / $max_memory_b"
    fi
else
    echo "$service_name is not running."
fi
