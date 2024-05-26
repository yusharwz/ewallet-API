CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fullname VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    image_url VARCHAR(255),
    pin VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    verification_code VARCHAR(35),
    expired_code TIMESTAMP WITHOUT TIME ZONE,
    phone_number VARCHAR(17) NOT NULL UNIQUE,
    roles VARCHAR(25) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    balance DECIMAL(15, 2) DEFAULT 0.00 CHECK (balance >= 0),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE payment_method (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('debit', 'credit')),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description VARCHAR(100),
    status VARCHAR(10) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES transactions(id),
    from_wallet_id UUID REFERENCES wallets(id),
    to_wallet_id UUID REFERENCES wallets(id),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE topup_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES transactions(id),
    payment_method_id UUID REFERENCES payment_method(id),
    payment_url VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON transactions(user_id);
CREATE INDEX idx_from_wallet_id ON wallet_transactions(from_wallet_id);
CREATE INDEX idx_to_wallet_id ON wallet_transactions(to_wallet_id);

INSERT INTO payment_method (id, payment_name)
VALUES
    ('087f9751-1dfc-474d-bdee-07ce44b1fe7a', 'Mandiri'),
    ('089e8004-2428-41f9-bf06-856082bb83d3', 'QRIS'),
    ('0eaad501-e44d-46e2-902a-9325c6c6c5eb', 'Indomaret'),
    ('0fafc78f-ebbf-421d-bc89-3246ce6198ad', 'CIMB NIaga'),
    ('220309af-cd3b-40e5-b353-6754c66f3831', 'Kredivo'),
    ('29690f9f-c6c4-4fda-acac-be91555b1f94', 'Akulaku'),
    ('2bed0329-499e-43b5-9b99-583b203ea102', 'BNI'),
    ('3863b99e-9909-486c-8ec1-b7a3162c9f97', 'BRI'),
    ('76954351-6cb3-496d-8866-d7f5772a04fe', 'Permata Bank'),
    ('91b75dee-155e-4ac3-9bfd-f8bed82b6189', 'ShopeePay/SPayLater'),
    ('9fa520e0-d10b-4be1-a6d7-e8b6fc635c5c', 'Debit/CreditCard'),
    ('b25a226e-82ab-4d29-a68e-6957fb7e21a9', 'Alfa Group'),
    ('cf51fa64-1686-4fee-a4e1-ea13c939f99b', 'BCA'),
    ('f9569b06-a389-4685-b3cc-89b13a111214', 'Gopay/GopayLater');