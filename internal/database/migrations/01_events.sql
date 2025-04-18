CREATE TABLE if NOT EXISTS calendar_events (
    id text NOT NULL PRIMARY KEY UNIQUE,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',

    date_from timestamp with time zone NOT NULL,
    date_to timestamp with time zone NOT NULL,
    tz text NOT NULL DEFAULT '',

    rrule text NOT NULL DEFAULT '',

    disabled boolean NOT NULL DEFAULT false,

    updated_at timestamp with time zone default now(),
    updated_by varchar(255) not null default ''
);

-- comments
COMMENT ON COLUMN calendar_events.rrule IS
$$RRULE of https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.10.$$;

COMMENT ON COLUMN calendar_events.date_to IS
'The end date of the event, exclusive.';

COMMENT ON COLUMN calendar_events.date_from IS
'The start date of the event, inclusive.';

COMMENT ON COLUMN calendar_events.disabled IS
'If the event is disabled, it will not be considered.';

COMMENT ON COLUMN calendar_events.tz IS
'Timezone of the event for getting back the original zone.';
