CREATE DATABASE IF NOT EXISTS auth_service;
CREATE DATABASE IF NOT EXISTS entity_service;
CREATE DATABASE IF NOT EXISTS incident_service;
CREATE DATABASE IF NOT EXISTS location_service;
CREATE DATABASE IF NOT EXISTS monitoring_service;
CREATE DATABASE IF NOT EXISTS notification_service;
CREATE DATABASE IF NOT EXISTS chat_service;
CREATE DATABASE IF NOT EXISTS user_service;

GRANT ALL PRIVILEGES ON auth_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON entity_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON incident_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON location_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON monitoring_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON notification_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON chat_service.* TO 'appuser'@'%';
GRANT ALL PRIVILEGES ON user_service.* TO 'appuser'@'%';