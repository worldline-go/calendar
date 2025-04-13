CREATE TABLE if NOT EXISTS calendar_relations (
    id text NOT NULL PRIMARY KEY UNIQUE,
    event_id text NOT NULL,

    -- relation with any
    code integer,
    country varchar(255),

    -- metadata
    updated_at timestamp with time zone default now(),
    updated_by varchar(255) not null default '',

    -- foreign keys
    FOREIGN KEY (event_id) REFERENCES calendar_events (id) ON DELETE CASCADE,

    CONSTRAINT unique_calendar_code UNIQUE (calendar_id, code, country)
);

CREATE INDEX ON calendar_relations (event_id);
CREATE INDEX ON calendar_relations (code);
CREATE INDEX ON calendar_relations (country);
