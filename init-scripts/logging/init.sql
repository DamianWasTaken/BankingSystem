CREATE TABLE public.accountStatus (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created TIMESTAMP NOT NULL
);
CREATE TABLE public.balanceTransactions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    transactionType VARCHAR(50) NOT NULL,
    created TIMESTAMP
);
CREATE TABLE public.interestConfiguration (
    id SERIAL PRIMARY KEY,
    interestRate DECIMAL(10, 2) NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE public.interestUserApplication (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    interestRate DECIMAL(10, 2) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    created TIMESTAMP NOT NULL,
    outcome VARCHAR(6) NOT NULL
);