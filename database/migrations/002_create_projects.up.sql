CREATE TABLE projects(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100)NOT NULL,
    description TEXT,
    manager_id INT NOT NULL,

    FOREIGN KEY (manager_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);