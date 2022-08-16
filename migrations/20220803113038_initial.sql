-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS persons
(
    id serial PRIMARY KEY,
    last_name varchar(200) NOT NULL,
    first_name varchar(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicles
(
    id serial PRIMARY KEY,
    brand varchar(200) NOT NULL,
    model varchar(200) NOT NULL,
    reg_number varchar(20) UNIQUE NOT NULL,
    person_id int NOT NULL REFERENCES persons(id) ON DELETE CASCADE
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vehicles;
DROP TABLE persons;
-- +goose StatementEnd
