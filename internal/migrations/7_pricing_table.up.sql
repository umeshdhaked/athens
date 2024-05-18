CREATE TABLE pricing (
     id INT AUTO_INCREMENT PRIMARY KEY,
     category VARCHAR(255) NOT NULL,
     sub_category VARCHAR(255) NOT NULL,
     pricing_type VARCHAR(255) NOT NULL,
     rates FLOAT NOT NULL,
     state VARCHAR(255) NOT NULL,
     created_at INT NOT NULL,
     updated_at INT NOT NULL,
     deleted_at INT
);