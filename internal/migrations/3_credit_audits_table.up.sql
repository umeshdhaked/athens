CREATE TABLE credit_audits (
       id INT AUTO_INCREMENT PRIMARY KEY,
       category VARCHAR(255) NOT NULL,
       sub_category VARCHAR(255) NOT NULL,
       deducted_amount INT NOT NULL,
       added_amount INT NOT NULL,
       credit_id INT NOT NULL,
       user_id INT NOT NULL,
       payment_order_id VARCHAR(255),
       created_at INT NOT NULL,
       updated_at INT NOT NULL,
       deleted_at INT
);