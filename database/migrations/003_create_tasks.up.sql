CREATE TABLE tasks(
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL,
    user_id INT NOT NULL,
    title VARCHAR(100)NOT NULL,
    description TEXT,
    status ENUM('pending', 'in_progress', 'completed')NOT NULL DEFAULT 'pending',
    progress INT NOT NULL DEFAULT 0,

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CHECK (PROGRESS >=0 AND progress <=100)


);