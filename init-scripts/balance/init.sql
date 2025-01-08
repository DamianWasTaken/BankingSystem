CREATE TABLE public.user (
    userId INT PRIMARY KEY,
    creditLimit INT,
    creditUsed INT
);

CREATE TABLE public.currencyList (
    currency VARCHAR(3) PRIMARY KEY
);

INSERT INTO public.currencyList (currency) VALUES
('GBP'),
('USD'),
('EUR');

CREATE TABLE public.account_currencies (
    currencyAccountId SERIAL PRIMARY KEY,
    userId INT REFERENCES public.user(userId),
    currency VARCHAR(3) REFERENCES currencyList(currency),
    balance INT
);


