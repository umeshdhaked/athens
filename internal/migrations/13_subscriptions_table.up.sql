CREATE TABLE subscription (
      id INT AUTO_INCREMENT PRIMARY KEY,
      pricing_id INT NOT NULL,
      user_id INT NOT NULL,
      type VARCHAR(255) NOT NULL,
      sub_type VARCHAR(255) NOT NULL,
      status VARCHAR(255) NOT NULL,
      added_by VARCHAR(255) NOT NULL,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);