package models

import (
	"github.com/worldline-go/types"
)

type Event struct {
	ID string `db:"id" json:"id" goqu:"skipupdate"`

	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`

	DateFrom types.Time `db:"date_from" json:"date_from" swaggertype:"string"`
	DateTo   types.Time `db:"date_to"   json:"date_to"   swaggertype:"string"`

	RRule    string `db:"rrule"    json:"rrule"`
	Disabled bool   `db:"disabled" json:"disabled"`

	UpdatedAt types.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
}

type Relation struct {
	ID      string             `db:"id"       json:"id"`
	EventID string             `db:"event_id" json:"event_id"`
	Code    types.Null[int64]  `db:"code"     json:"code"     swaggertype:"integer"`
	Country types.Null[string] `db:"country"  json:"country"  swaggertype:"string"`

	UpdatedAt types.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
}
