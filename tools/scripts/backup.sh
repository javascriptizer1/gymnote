#!/bin/bash

DATE=$(date +%Y-%m-%d_%H-%M)
BACKUP_NAME="gymnote_backup_$DATE"

clickhouse-backup create $BACKUP_NAME
clickhouse-backup upload $BACKUP_NAME