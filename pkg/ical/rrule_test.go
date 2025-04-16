package ical

import (
	"reflect"
	"testing"
	"time"
)

func TestMatchRRuleAt(t *testing.T) {
	locationNewYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}

	type args struct {
		rrule   *RRule
		dtstart time.Time
		dtend   time.Time
		search  time.Time
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 time.Time
		want2 bool
	}{
		{
			name: "U.S. Presidential Election",
			args: args{
				rrule: &RRule{
					Freq:       "YEARLY",
					Interval:   4,
					ByMonth:    []int{11},
					ByMonthDay: []int{2, 3, 4, 5, 6, 7, 8},
					ByDay:      []string{"TU"},
				},
				dtstart: time.Date(1996, 11, 5, 9, 0, 0, 0, locationNewYork),
				dtend:   time.Date(1996, 11, 6, 0, 0, 0, 0, locationNewYork),
				search:  time.Date(2000, 11, 7, 9, 0, 0, 0, locationNewYork),
			},
			want:  time.Date(2000, 11, 7, 9, 0, 0, 0, locationNewYork),
			want1: time.Date(2000, 11, 8, 0, 0, 0, 0, locationNewYork),
			want2: true,
		},
		{
			name: "Every Thursday in March, forever",
			args: args{
				rrule: &RRule{
					Freq:    "YEARLY",
					ByMonth: []int{3},
					ByDay:   []string{"TH"},
				},
				dtstart: time.Date(1996, 3, 1, 0, 0, 0, 0, locationNewYork),
				dtend:   time.Date(1996, 3, 2, 0, 0, 0, 0, locationNewYork),
				search:  time.Date(1999, 3, 18, 0, 0, 0, 0, locationNewYork),
			},
			want:  time.Date(1999, 3, 18, 0, 0, 0, 0, locationNewYork),
			want1: time.Date(1999, 3, 19, 0, 0, 0, 0, locationNewYork),
			want2: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := MatchRRuleAt(tt.args.rrule, tt.args.dtstart, tt.args.dtend, tt.args.search)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchRRuleAt() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MatchRRuleAt() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("MatchRRuleAt() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
