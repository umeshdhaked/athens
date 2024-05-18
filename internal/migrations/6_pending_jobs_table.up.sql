CREATE TABLE pending_jobs (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      type VARCHAR(255) NOT NULL,
      status VARCHAR(255) NOT NULL,
      extra TEXT,
      created_at INT NOT NULL,
      updated_at INT NOT NULL,
      deleted_at INT
);