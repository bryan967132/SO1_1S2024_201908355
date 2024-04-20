CREATE DATABASE IF NOT EXISTS music_data;
USE music_data;

CREATE TABLE IF NOT EXISTS music (
    name VARCHAR(255),
    album VARCHAR(255),
    year VARCHAR(255),
    ranking VARCHAR(255)
);