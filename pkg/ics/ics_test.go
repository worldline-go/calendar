package ics

import (
	"testing"
	"time"

	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/types"
)

func TestGenerateICS(t *testing.T) {
	type args struct {
		events []models.Event
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "23 Nisan",
			args: args{
				events: []models.Event{
					{
						Name:        "23 Nisan Ulusal Egemenlik ve Çocuk Bayramı",
						Description: "23 Nisan Ulusal Egemenlik ve Çocuk Bayramı",
						DateFrom:    types.Time{Time: time.Date(2023, 4, 23, 0, 0, 0, 0, time.UTC)},
						DateTo:      types.Time{Time: time.Date(2023, 4, 24, 0, 0, 0, 0, time.UTC)},
						RRule:       "FREQ=YEARLY;BYMONTH=4;BYMONTHDAY=23",
						Disabled:    false,
						UpdatedAt:   types.Time{Time: time.Now()},
						UpdatedBy:   "system",
					},
				},
			},
			want:    "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//worldline-go//calendar//EN\r\nBEGIN:VEVENT\r\nSUMMARY:23 Nisan Ulusal Egemenlik ve Çocuk Bayramı\r\nDESCRIPTION:23 Nisan Ulusal Egemenlik ve Çocuk Bayramı\r\nDTSTART;VALUE=DATE:20230423\r\nDTEND;VALUE=DATE:20230424\r\nRRULE:FREQ=YEARLY;BYMONTH=4;BYMONTHDAY=23\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateICS(tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateICS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateICS() = %v, want %v", got, tt.want)
			}
		})
	}
}
