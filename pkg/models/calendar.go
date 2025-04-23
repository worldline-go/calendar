package models

import (
	"github.com/worldline-go/types"
)

type Event struct {
	ID string `db:"id" json:"id" goqu:"skipupdate"`

	Name        string             `db:"name"        json:"name"`
	Description string             `db:"description" json:"description"`
	EventGroup  types.Null[string] `db:"event_group" json:"event_group"`

	DateFrom types.Time `db:"date_from" json:"date_from" swaggertype:"string"`
	DateTo   types.Time `db:"date_to"   json:"date_to"   swaggertype:"string"`
	Tz       string     `db:"tz"        json:"tz"`
	AllDay   bool       `db:"all_day"   json:"all_day"`

	RRule    string `db:"rrule"    json:"rrule"`
	Disabled bool   `db:"disabled" json:"disabled"`

	UpdatedAt types.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
}

type Relation struct {
	Entity string `db:"entity" json:"entity"`

	EventID    types.Null[string] `db:"event_id"    json:"event_id"    swaggertype:"string"`
	EventGroup types.Null[string] `db:"event_group" json:"event_group" swaggertype:"string"`

	UpdatedAt types.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
}
