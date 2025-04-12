CREATE TABLE if NOT EXISTS holiday (
    id text NOT NULL PRIMARY KEY UNIQUE,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',

    date_from timestamp with time zone,
    date_to timestamp with time zone,

    years text NOT NULL DEFAULT '',

    disabled boolean NOT NULL DEFAULT false,

    updated_at timestamp with time zone default now(),
    updated_by varchar(255) not null default ''
);

-- comments
COMMENT ON COLUMN holiday.years IS
$$Years in which the holiday is valid. Comma separated list of years or ranges.
Example: 2020,2021,2023-2025,2030-*,*,*-1990,1995-*4.
If empty, the holiday is valid for only date_from and date_to.
If the holiday is valid for a range of years, the both range is inclusive.
If the holiday is valid for all years, use *.
If the holiday is valid for a range of years, use * as the start or end of the range.
*4 means every 4 years starting from the first year of the range.$$;

COMMENT ON COLUMN holiday.date_to IS
'The end date of the holiday, exclusive.';

COMMENT ON COLUMN holiday.date_from IS
'The start date of the holiday, inclusive.';

COMMENT ON COLUMN holiday.disabled IS
'If the holiday is disabled, it will not be considered.';
