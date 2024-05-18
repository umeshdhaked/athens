CREATE TABLE invoice
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    order_id   VARCHAR(255) NOT NULL,
    status     VARCHAR(255) NOT NULL,
    user_id    FLOAT        NOT NULL,
    receipt    VARCHAR(255) NOT NULL,
    created_at INT          NOT NULL,
    updated_at INT          NOT NULL,
    deleted_at INT
);