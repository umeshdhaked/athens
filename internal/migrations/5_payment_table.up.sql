CREATE TABLE payments
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    amount           FLOAT        NOT NULL,
    amount_paid      FLOAT        NOT NULL,
    mobile           VARCHAR(255) NOT NULL,
    user_id          FLOAT        NOT NULL,
    order_created_at FLOAT        NOT NULL,
    amount_due       FLOAT        NOT NULL,
    currency         VARCHAR(255) NOT NULL,
    receipt          VARCHAR(255) NOT NULL,
    order_id         VARCHAR(255) NOT NULL,
    offer_id         VARCHAR(255) NOT NULL,
    entity           VARCHAR(255) NOT NULL,
    attempts         FLOAT        NOT NULL,
    status           VARCHAR(255) NOT NULL,
    created_at       INT          NOT NULL,
    updated_at       INT          NOT NULL,
    deleted_at       INT
);