CREATE TABLE promo_phones (
      id INT AUTO_INCREMENT PRIMARY KEY,
      mobile VARCHAR(255) NOT NULL,
      is_already_contacted VARCHAR(255),
      comment TEXT,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);