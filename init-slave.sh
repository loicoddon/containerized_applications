#!/bin/bash

# Wait for the master to be ready
sleep 10

mysql -h mysql_server -u root -proot -e "CREATE USER 'slave_user'@'%' IDENTIFIED BY 'slave_password';
GRANT REPLICATION SLAVE ON *.* TO 'slave_user'@'%';
FLUSH PRIVILEGES;"

# Retrieve the master log file and position
MASTER_STATUS=$(mysql -h mysql_server -u root -proot -e "SHOW MASTER STATUS;" 2> /bin/null | awk 'NR==2') #| grep mysql-bin)
MASTER_LOG_FILE=$(echo $MASTER_STATUS | awk '{ print $1 }')
MASTER_LOG_POS=$(echo $MASTER_STATUS | awk '{ print $2 }')
echo "Phase 1 terminé"

# Set up replication on the slave
mysql -h mysql_slave -u root -proot -e "CREATE DATABASE groupie_DB;
CHANGE MASTER TO
  MASTER_HOST='$MASTER_HOST',
  MASTER_USER='slave_user',
  MASTER_PASSWORD='slave_password',
  MASTER_LOG_FILE='$MASTER_LOG_FILE',
  MASTER_LOG_POS=$MASTER_LOG_POS;
START SLAVE;"