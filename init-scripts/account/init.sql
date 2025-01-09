CREATE TABLE public.user (
    email VARCHAR(200) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    password VARCHAR(200) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);
