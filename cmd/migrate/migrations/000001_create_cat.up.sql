CREATE TABLE IF NOT EXISTS cats(
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL CHECK (years_of_experience >= 0),
    breed VARCHAR(255) NOT NULL,
    salary DECIMAL(10,2) NOT NULL CHECK (salary >= 0)
);