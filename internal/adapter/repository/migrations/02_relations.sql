CREATE TABLE if NOT EXISTS calendar_relations (
    entity text NOT NULL,

    -- relation with any
    event_group text,
    event_id text,

    -- metadata
    updated_at timestamp with time zone default now(),
    updated_by varchar(255) not null default '',

    -- foreign keys
    FOREIGN KEY (event_id) REFERENCES calendar_events (id) ON DELETE CASCADE,

    CONSTRAINT unique_calendar_entity UNIQUE (entity, event_group, event_id)
);
