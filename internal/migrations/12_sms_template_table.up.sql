CREATE TABLE sms_template (
      id INT AUTO_INCREMENT PRIMARY KEY,
      user_id VARCHAR(255) NOT NULL,
      sender_id INT NOT NULL,
      sender_code VARCHAR(255) NOT NULL,
      template_code VARCHAR(255) NOT NULL,
      body TEXT NOT NULL,
      status VARCHAR(255) NOT NULL,
      type VARCHAR(255) NOT NULL,
      length INT NOT NULL,
      language VARCHAR(255) NOT NULL,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);