CREATE TABLE pdf_metadata (
    id serial PRIMARY KEY,
    email VARCHAR(255) REFERENCES "Users"("Email") NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    upload_date TIMESTAMPTZ,
    document_id VARCHAR(255) NOT NULL,
    storage_reference VARCHAR(255) NOT NULL
);
