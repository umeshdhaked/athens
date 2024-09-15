CREATE TABLE kyc (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    user_name  VARCHAR(255) NOT NULL,
    mobile VARCHAR(255) NOT NULL,
    document_type VARCHAR(255) NOT NULL,
    kyc_doc_link VARCHAR(255) NOT NULL,
    photo_link VARCHAR(255) NOT NULL,
    is_verified VARCHAR(255) NOT NULL,
    comment VARCHAR(255) NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    deleted_at INT
);