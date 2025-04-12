package service

import (
	"testing"
	"time"

	"github.com/worldline-go/types"
)

func TestCheckYear(t *testing.T) {
	type args struct {
		year     int
		dateFrom types.Null[types.Time]
		dateTo   types.Null[types.Time]
		years    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid year within range",
			args: args{
				year:     2023,
				dateFrom: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)}),
				dateTo:   types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
				years:    "2020,2021,2023-2025",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Without year range",
			args: args{
				year:   2023,
				dateTo: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "4th year",
			args: args{
				year:     2029,
				dateFrom: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)}),
				dateTo:   types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
				years:    "2020,2021,2025-*4",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckYear(tt.args.year, tt.args.dateFrom, tt.args.dateTo, tt.args.years)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckYear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDate(t *testing.T) {
	type args struct {
		date     types.Time
		dateFrom types.Null[types.Time]
		dateTo   types.Null[types.Time]
		years    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid date within range",
			args: args{
				date:     types.Time{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
				dateFrom: types.NewNull(types.Time{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}),
				dateTo:   types.NewNull(types.Time{Time: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC)}),
				years:    "*-2023",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckDate(tt.args.date, tt.args.dateFrom, tt.args.dateTo, tt.args.years)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDate(t *testing.T) {
	type args struct {
		date     types.Time
		dateFrom types.Null[types.Time]
		dateTo   types.Null[types.Time]
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid date within range",
			args: args{
				date:     types.Time{Time: types.Time{}.Time.AddDate(2023, 0, 0)},
				dateFrom: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)}),
				dateTo:   types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
			},
			want: true,
		},
		{
			name: "Inclusive start date",
			args: args{
				date:     types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)},
				dateFrom: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)}),
				dateTo:   types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
			},
			want: true,
		},
		{
			name: "Exclusive end date",
			args: args{
				date:     types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)},
				dateFrom: types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2020, 0, 0)}),
				dateTo:   types.NewNull(types.Time{Time: types.Time{}.Time.AddDate(2025, 0, 0)}),
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkDate(tt.args.date, tt.args.dateFrom, tt.args.dateTo)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
