CREATE TABLE account_deletion_requests (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP NULL DEFAULT NULL,
    reviewed_by INT NULL,

    PRIMARY KEY (id),

    KEY user_id (user_id),
    KEY reviewed_by (reviewed_by),

    CONSTRAINT fk_deletion_request_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_deletion_request_admin
        FOREIGN KEY (reviewed_by)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT deletion_request_status_check
        CHECK (status IN ('pending', 'approved', 'rejected'))
);