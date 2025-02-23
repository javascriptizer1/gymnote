#!/bin/bash
(crontab -l 2>/dev/null; echo "0 3 */3 * * docker compose run -f ../../docker-compose.yaml --rm clickhouse-backup") | crontab -
