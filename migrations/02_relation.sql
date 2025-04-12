CREATE TABLE if NOT EXISTS holiday_relation (
    id text NOT NULL PRIMARY KEY UNIQUE,
    holiday_id text NOT NULL,

    -- relation with any
    code integer,
    country varchar(255),

    -- metadata
    updated_at timestamp with time zone default now(),
    updated_by varchar(255) not null default '',

    -- foreign keys
    FOREIGN KEY (holiday_id) REFERENCES holiday (id) ON DELETE CASCADE,

    CONSTRAINT unique_holiday_code UNIQUE (holiday_id, code, country)
);

CREATE INDEX ON holiday_relation (holiday_id);
CREATE INDEX ON holiday_relation (code);
CREATE INDEX ON holiday_relation (country);
