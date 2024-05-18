CREATE TABLE sms_audit (
       id INT AUTO_INCREMENT PRIMARY KEY,
       user_id VARCHAR(255) NOT NULL,
       credits_consumed INT NOT NULL,
       template_id INT NOT NULL,
       sender_code VARCHAR(255) NOT NULL,
       contact_id INT NOT NULL,
       status VARCHAR(255) NOT NULL,
       triggered_mode VARCHAR(255) NOT NULL,
       created_at INT NOT NULL,
       updated_at INT NOT NULL,
       deleted_at INT
);