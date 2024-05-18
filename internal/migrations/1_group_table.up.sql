CREATE TABLE contact_groups (
     id INT AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     user_id INT NOT NULL,
     column_names TEXT,
     created_at INT NOT NULL,
     updated_at INT,
     deleted_at INT DEFAULT NULL
);