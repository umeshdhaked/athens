CREATE TABLE sms_campaign (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      scheduled_at INT NOT NULL,
      status VARCHAR(255) NOT NULL,
      user_id VARCHAR(255) NOT NULL,
      template_id INT NOT NULL,
      sender_code VARCHAR(255) NOT NULL,
      group_name VARCHAR(255) NOT NULL,
      type VARCHAR(255) NOT NULL,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);