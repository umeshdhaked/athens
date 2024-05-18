CREATE TABLE cron_processing (
     id INT AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     batch INT NOT NULL,
     in_progress INT NOT NULL,
     status VARCHAR(255) NOT NULL,
     last_evaluated_id INT NOT NULL,
     created_at INT NOT NULL,
     updated_at INT NOT NULL,
     deleted_at INT
);