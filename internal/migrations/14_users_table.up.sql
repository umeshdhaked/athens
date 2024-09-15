CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mobile VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    kyc_done VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    deleted_at INT
);