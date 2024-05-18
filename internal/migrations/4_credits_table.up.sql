CREATE TABLE credits (
     id INT AUTO_INCREMENT PRIMARY KEY,
     user_id INT NOT NULL,
     balance INT NOT NULL,
     balance_left INT NOT NULL,
     created_at INT NOT NULL,
     updated_at INT NOT NULL,
     deleted_at INT
);