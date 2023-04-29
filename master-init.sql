CREATE USER 'slave_user'@'mysql_slave' IDENTIFIED BY 'slave_password';
GRANT REPLICATION SLAVE ON *.* TO 'slave_user'@'mysql_slave';
FLUSH PRIVILEGES;