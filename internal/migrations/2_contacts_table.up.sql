CREATE TABLE contacts (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      mobile VARCHAR(20) NOT NULL,
      email VARCHAR(255) NOT NULL,
      group_name VARCHAR(255),
      additional TEXT,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);