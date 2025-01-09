CREATE TABLE public.currencyList (
    currency VARCHAR(3) PRIMARY KEY,
    rate decimal(10,2) NOT NULL
);

INSERT INTO public.currencyList (currency, rate) VALUES
('GBP', 1.00),
('USD', 1.25),
('EUR', 1.15);

CREATE TABLE public.account_currencies (
    email VARCHAR(200) NOT NULL,
    currency VARCHAR(3) REFERENCES currencyList(currency) NOT NULL,
    balance decimal(10,2) NOT NULL,
    PRIMARY KEY (Email, currency)
);


