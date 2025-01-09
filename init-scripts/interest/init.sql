CREATE TABLE public.currencyList (
    currency VARCHAR(3) PRIMARY KEY
);

INSERT INTO public.currencyList (currency) VALUES
('GBP'),
('USD'),
('EUR');

CREATE TABLE public.interestRate(
    interest DECIMAL(4,2) PRIMARY KEY
);

INSERT INTO public.interestRate (interest) VALUES
(0.00);

CREATE TABLE public.interestUser (
    interestId SERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL,
    accountCurrency VARCHAR(100) UNIQUE NOT NULL REFERENCES public.currencyList(currency),
    interestFrequency INT,
    nextInterestDate DATE
);