CREATE USER 'slave_user'@'%' IDENTIFIED BY 'slave_password';
GRANT REPLICATION SLAVE ON *.* TO 'slave_user'@'%';
FLUSH PRIVILEGES;