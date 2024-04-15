CREATE DATABASE IF NOT EXISTS P1SO1;

USE P1SO1;

CREATE TABLE IF NOT EXISTS RAM (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usado FLOAT,
    disponible FLOAT,
    tiempo DATETIME
);

CREATE TABLE IF NOT EXISTS CPU (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usado FLOAT,
    disponible FLOAT,
    tiempo DATETIME
);

SELECT * FROM RAM;
SELECT * FROM CPU;

TRUNCATE TABLE RAM;
TRUNCATE TABLE CPU;