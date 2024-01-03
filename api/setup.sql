SELECT user FROM mysql.user WHERE user = 'golang';

CREATE USER golang;
ALTER USER golang WITH PASSWORD 'golang';
CREATE DATABASE IF NOT EXISTS golang;
GRANT ALL PRIVILEGES ON DATABASE golang TO golang;
