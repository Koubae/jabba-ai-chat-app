CREATE DATABASE IF NOT EXISTS chat_identity CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE chat_identity;

CREATE TABLE IF NOT EXISTS users
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    application_id VARCHAR(255) NOT NULL,
    username       VARCHAR(255) NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,
    created        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY user__application_id_name__unique (application_id, username)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

