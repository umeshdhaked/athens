CREATE TABLE otp (
      id INT AUTO_INCREMENT PRIMARY KEY,
      mobile VARCHAR(255) NOT NULL UNIQUE,
      otp VARCHAR(255),
      exp INT,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);