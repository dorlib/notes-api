CREATE DATABASE IF NOT EXISTS notesapp;

USE notesapp;

CREATE TABLE IF NOT EXISTS notes (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
                                     content TEXT NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
