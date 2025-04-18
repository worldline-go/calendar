package ical

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/types"
)

func TestGenerateICS(t *testing.T) {
	tzIstanbul, _ := time.LoadLocation("Europe/Istanbul")
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
						DateFrom:    types.Time{Time: time.Date(2023, 4, 23, 0, 0, 0, 0, tzIstanbul)},
						DateTo:      types.Time{Time: time.Date(2023, 4, 24, 0, 0, 0, 0, tzIstanbul)},
						RRule:       "FREQ=YEARLY;BYMONTH=4;BYMONTHDAY=23",
						Disabled:    false,
						UpdatedAt:   types.Time{Time: time.Now()},
						UpdatedBy:   "system",
					},
				},
			},
			want: "BEGIN:VCALENDAR\r\n" +
				"VERSION:2.0\r\n" +
				"PRODID:-//worldline-go//calendar//EN\r\n" +
				"BEGIN:VEVENT\r\n" +
				"UID:\r\n" +
				"CATEGORIES:Holidays\r\n" +
				"CLASS:public\r\n" +
				"STATUS:CONFIRMED\r\n" +
				"SUMMARY:23 Nisan Ulusal Egemenlik ve Çocuk Bayramı\r\n" +
				"DESCRIPTION:23 Nisan Ulusal Egemenlik ve Çocuk Bayramı\r\n" +
				"X-MICROSOFT-CDO-ALLDAYEVENT:TRUE\r\n" +
				"DTSTART;VALUE=DATE:20230423\r\n" +
				"DTEND;VALUE=DATE:20230424\r\n" +
				"RRULE:FREQ=YEARLY;BYMONTH=4;BYMONTHDAY=23\r\n" +
				"TRANSP:TRANSPARENT\r\n" +
				"END:VEVENT\r\n" +
				"END:VCALENDAR\r\n",
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

func TestParseICS(t *testing.T) {
	tzIstanbul, _ := time.LoadLocation("Europe/Istanbul")
	type args struct {
		data []byte
		tz   string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Event
		wantErr bool
	}{
		{
			name: "23 Nisan",
			args: args{
				data: []byte(`
BEGIN:VEVENT
SUMMARY:Atatürk'ü Anma\, Gençlik ve Spor Günü
DTSTART;VALUE=DATE:20240519
DTEND;VALUE=DATE:20240520
DTSTAMP:20241008T090751Z
UID:f6d4e8a07317c9779f0fa9ea3152f722-2024
CATEGORIES:Holidays
CLASS:public
DESCRIPTION:National holiday -  Türkiye'de pek çok kişi her yıl 19 May
 ıs'ta Atatürk Anma\, Gençlik ve Spor Günü'nü spor etkinliklerine kat
 ılarak ve bu gün 1919 yılında başlayan Kurtuluş Savaşı'nı hatırl
 ayarak kutlamaktadır.
LAST-MODIFIED:20241008T090751Z
TRANSP:transparent
END:VEVENT
`),
				tz: "Europe/Istanbul",
			},
			want: []models.Event{
				{
					ID:          "f6d4e8a07317c9779f0fa9ea3152f722-2024",
					Name:        "Atatürk'ü Anma, Gençlik ve Spor Günü",
					Description: "National holiday -  Türkiye'de pek çok kişi her yıl 19 Mayıs'ta Atatürk Anma, Gençlik ve Spor Günü'nü spor etkinliklerine katılarak ve bu gün 1919 yılında başlayan Kurtuluş Savaşı'nı hatırlayarak kutlamaktadır.",
					DateFrom:    types.Time{Time: time.Date(2024, 5, 19, 0, 0, 0, 0, tzIstanbul)},
					DateTo:      types.Time{Time: time.Date(2024, 5, 20, 0, 0, 0, 0, tzIstanbul)},
					TZone:       "Europe/Istanbul",
					RRule:       "",
					Disabled:    false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz, err := time.LoadLocation(tt.args.tz)
			if err != nil {
				t.Fatalf("time.LoadLocation() error = %v", err)
			}

			got, err := ParseICS(bytes.NewReader(tt.args.data), tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseICS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseICS() = \n%#v\n, want \n%#v\n", got, tt.want)
			}
		})
	}
}
