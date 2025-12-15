package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	testCases := []struct {
		name string
		dob  time.Time
		now  time.Time
		want int
	}{
		{
			name: "birthday passed this year",
			dob:  time.Date(1990, 1, 10, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			want: 35,
		},
		{
			name: "birthday today",
			dob:  time.Date(2000, 12, 15, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
			want: 25,
		},
		{
			name: "birthday upcoming this year",
			dob:  time.Date(1995, 12, 31, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			want: 29,
		},
		{
			name: "future dob",
			dob:  time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CalculateAge(tc.dob, tc.now); got != tc.want {
				t.Fatalf("CalculateAge() = %d, want %d", got, tc.want)
			}
		})
	}
}



