package models

import "time"

// CalculateAge returns the completed years between dob and now.
func CalculateAge(dob time.Time, now time.Time) int {
	years := now.Year() - dob.Year()

	// If birthday hasn't occurred yet this year, subtract one.
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}



