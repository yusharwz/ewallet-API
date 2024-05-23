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
    phone_number VARCHAR(17) NOT NULL,
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
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
);a

CREATE TABLE topup_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES transactions(id),
    payment_method_id UUID REFERENCES payment_method(id),
    payment_url VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Membuat indeks untuk tabel transactions
CREATE INDEX idx_user_id ON transactions(user_id);
-- Membuat indeks untuk tabel wallet_transactions
CREATE INDEX idx_from_wallet_id ON wallet_transactions(from_wallet_id);
CREATE INDEX idx_to_wallet_id ON wallet_transactions(to_wallet_id);
-- Membuat indeks untuk tabel topup_transactions
CREATE INDEX idx_from_wallet_id_topup ON topup_transactions(from_wallet_id);
