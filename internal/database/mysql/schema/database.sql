CREATE DATABASE IF NOT EXISTS chatting;

CREATE TABLE IF NOT EXISTS room (
    `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL UNIQUE,
    `createdAt` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chat (
    `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `room` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `message` varchar(255) NOT NULL,
    `when` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS serverInfo (
    `ip` varchar(255) PRIMARY KEY NOT NULL,
    `available` bool NOT NULL
);