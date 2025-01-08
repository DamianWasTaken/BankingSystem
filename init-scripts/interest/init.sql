CREATE TABLE public.currencyList (
    currency PRIMARY KEY
);

INSERT INTO public.currencyList (currency) VALUES
('GBP'),
('USD'),
('EUR');

CREATE TABLE public.interest (
    interestId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    accountCurrency VARCHAR(100) UNIQUE NOT NULL REFERENCES public.account_currencies(accountCurrency),
    interestPercent DECIMAL(4,2),
    interestFrequency INT,
    nextInterestDate DATE
);


