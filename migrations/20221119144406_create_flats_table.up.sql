CREATE TABLE IF NOT EXISTS flats(
    id UUID PRIMARY KEY,
    service_id VARCHAR (32) NOT NULL,
    url VARCHAR (128) UNIQUE NOT NULL,
    photo_urls JSONB NOT NULL,
    title VARCHAR (512) NOT NULL,
    description text NOT NULL,
    area smallint NOT NULL,
    rooms smallint NOT NULL,
    floor smallint NOT NULL,
    price integer NOT NULL,
    contact_name VARCHAR (128) NOT NULL,
    phone VARCHAR (32) NOT NULL,
    is_agency boolean NOT NULL,
    published_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);
