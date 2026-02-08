package utils

import (
	"fmt"
	"time"
)

// NepaliDate represents a date in Bikram Sambat (BS) calendar
type NepaliDate struct {
	Year  int
	Month int // 1-12
	Day   int // 1-32 (some months have 32 days)
}

// String returns the date in YYYY-MM-DD format
func (nd NepaliDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", nd.Year, nd.Month, nd.Day)
}

// NepaliMonthName returns the Nepali month name
func (nd NepaliDate) NepaliMonthName() string {
	months := []string{
		"", // 0-index placeholder
		"Baishakh", "Jestha", "Ashad", "Shrawan",
		"Bhadra", "Ashwin", "Kartik", "Mangsir",
		"Poush", "Magh", "Falgun", "Chaitra",
	}
	if nd.Month >= 1 && nd.Month <= 12 {
		return months[nd.Month]
	}
	return "Unknown"
}

// EnglishMonthName returns the approximate English month
func (nd NepaliDate) EnglishMonthName() string {
	months := []string{
		"", // 0-index placeholder
		"April-May", "May-June", "June-July", "July-August",
		"August-September", "September-October", "October-November", "November-December",
		"December-January", "January-February", "February-March", "March-April",
	}
	if nd.Month >= 1 && nd.Month <= 12 {
		return months[nd.Month]
	}
	return "Unknown"
}

// Days in each Nepali month for years 2070-2090 BS
// Note: This is a simplified mapping. In production, use a complete calendar library
var nepaliMonthDays = map[int][]int{
	2080: {31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30}, // 2023-2024 AD
	2081: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, // 2024-2025 AD
	2082: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, // 2025-2026 AD
	2083: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, // 2026-2027 AD
	2084: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, // 2027-2028 AD
	2085: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, // 2028-2029 AD
	2086: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, // 2029-2030 AD
	2087: {30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, // 2030-2031 AD
	2088: {31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, // 2031-2032 AD
	2089: {31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, // 2032-2033 AD
	2090: {31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, // 2033-2034 AD
}

// Reference date: 2080-01-01 BS = 2023-04-14 AD
var referenceBS = NepaliDate{Year: 2080, Month: 1, Day: 1}
var referenceAD = time.Date(2023, 4, 14, 0, 0, 0, 0, time.UTC)

// ADToBS converts Gregorian (AD) date to Bikram Sambat (BS)
func ADToBS(ad time.Time) NepaliDate {
	// Calculate days difference from reference
	daysDiff := int(ad.Sub(referenceAD).Hours() / 24)

	// Start from reference BS date
	year := referenceBS.Year
	month := referenceBS.Month
	day := referenceBS.Day + daysDiff

	// Adjust for days overflow
	for day > 0 {
		monthDays := getDaysInMonth(year, month)
		if day <= monthDays {
			break
		}
		day -= monthDays
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	// Adjust for negative days
	for day <= 0 {
		month--
		if month < 1 {
			month = 12
			year--
		}
		day += getDaysInMonth(year, month)
	}

	return NepaliDate{Year: year, Month: month, Day: day}
}

// BSToAD converts Bikram Sambat (BS) date to Gregorian (AD)
func BSToAD(bs NepaliDate) time.Time {
	// Calculate total days from reference BS to target BS
	totalDays := 0

	// Add/subtract years
	if bs.Year > referenceBS.Year {
		for y := referenceBS.Year; y < bs.Year; y++ {
			totalDays += getTotalDaysInYear(y)
		}
	} else if bs.Year < referenceBS.Year {
		for y := bs.Year; y < referenceBS.Year; y++ {
			totalDays -= getTotalDaysInYear(y)
		}
	}

	// Add/subtract months in target year
	if bs.Month > referenceBS.Month {
		for m := referenceBS.Month; m < bs.Month; m++ {
			totalDays += getDaysInMonth(bs.Year, m)
		}
	} else if bs.Month < referenceBS.Month {
		for m := bs.Month; m < referenceBS.Month; m++ {
			totalDays -= getDaysInMonth(bs.Year, m)
		}
	}

	// Add/subtract days
	totalDays += bs.Day - referenceBS.Day

	// Add to reference AD date
	return referenceAD.AddDate(0, 0, totalDays)
}

// getDaysInMonth returns the number of days in a Nepali month
func getDaysInMonth(year, month int) int {
	if monthDays, ok := nepaliMonthDays[year]; ok && month >= 1 && month <= 12 {
		return monthDays[month-1]
	}
	// Default fallback
	return 30
}

// getTotalDaysInYear returns total days in a Nepali year
func getTotalDaysInYear(year int) int {
	total := 0
	for month := 1; month <= 12; month++ {
		total += getDaysInMonth(year, month)
	}
	return total
}

// GetCurrentNepaliDate returns today's date in BS
func GetCurrentNepaliDate() NepaliDate {
	return ADToBS(time.Now())
}

// GetFiscalYearName returns fiscal year name from BS date
// e.g., 2082-04-01 -> "2082/83"
func GetFiscalYearName(bs NepaliDate) string {
	// Fiscal year starts from Shrawan (month 4)
	if bs.Month >= 4 {
		return fmt.Sprintf("%d/%02d", bs.Year, (bs.Year+1)%100)
	}
	return fmt.Sprintf("%d/%02d", bs.Year-1, bs.Year%100)
}

// GetFiscalYearDates returns start and end dates for a fiscal year
// Returns both BS and AD dates
func GetFiscalYearDates(fiscalYearName string) (startBS NepaliDate, endBS NepaliDate, startAD time.Time, endAD time.Time) {
	// Parse fiscal year name (e.g., "2082/83")
	var year int
	fmt.Sscanf(fiscalYearName, "%d/", &year)

	// Fiscal year starts on Shrawan 1 (month 4, day 1)
	startBS = NepaliDate{Year: year, Month: 4, Day: 1}

	// Fiscal year ends on Ashad 32 (month 3, day 32 of next year)
	endBS = NepaliDate{Year: year + 1, Month: 3, Day: 32}

	// Convert to AD
	startAD = BSToAD(startBS)
	endAD = BSToAD(endBS)

	return
}

// FormatNepaliDate formats a Nepali date in various formats
func FormatNepaliDate(bs NepaliDate, format string) string {
	switch format {
	case "YYYY-MM-DD":
		return bs.String()
	case "DD MMM YYYY":
		return fmt.Sprintf("%d %s %d", bs.Day, bs.NepaliMonthName()[:3], bs.Year)
	case "DD MMMM YYYY":
		return fmt.Sprintf("%d %s %d", bs.Day, bs.NepaliMonthName(), bs.Year)
	default:
		return bs.String()
	}
}

// ParseNepaliDate parses a Nepali date string (YYYY-MM-DD)
func ParseNepaliDate(dateStr string) (NepaliDate, error) {
	var nd NepaliDate
	_, err := fmt.Sscanf(dateStr, "%d-%d-%d", &nd.Year, &nd.Month, &nd.Day)
	if err != nil {
		return NepaliDate{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Validate
	if nd.Month < 1 || nd.Month > 12 {
		return NepaliDate{}, fmt.Errorf("invalid month: %d", nd.Month)
	}

	maxDays := getDaysInMonth(nd.Year, nd.Month)
	if nd.Day < 1 || nd.Day > maxDays {
		return NepaliDate{}, fmt.Errorf("invalid day: %d (max: %d)", nd.Day, maxDays)
	}

	return nd, nil
}
