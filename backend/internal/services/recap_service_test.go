package services

import (
	"testing"
	"time"
)

func TestCalculateDateRange(t *testing.T) {
	service := &RecapService{}

	// Test date: 2024-01-15 (Monday)
	testDate := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)

	tests := []struct {
		name          string
		filter        string
		date          *time.Time
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "Day filter",
			filter:        "day",
			date:          &testDate,
			expectedStart: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:          "Week filter",
			filter:        "week",
			date:          &testDate,
			expectedStart: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),            // Monday
			expectedEnd:   time.Date(2024, 1, 21, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:          "Month filter",
			filter:        "month",
			date:          &testDate,
			expectedStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, 1, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:          "Invalid filter defaults to day",
			filter:        "invalid",
			date:          &testDate,
			expectedStart: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := service.calculateDateRange(tt.filter, tt.date)

			if !start.Equal(tt.expectedStart) {
				t.Errorf("Expected start date %v, got %v", tt.expectedStart, start)
			}

			if !end.Equal(tt.expectedEnd) {
				t.Errorf("Expected end date %v, got %v", tt.expectedEnd, end)
			}
		})
	}
}

func TestCalculateDateRangeWeekWithSunday(t *testing.T) {
	service := &RecapService{}

	// Test with Sunday (2024-01-14)
	testDate := time.Date(2024, 1, 14, 12, 30, 45, 0, time.UTC)

	start, end := service.calculateDateRange("week", &testDate)

	// Should start from Monday (2024-01-08) and end on Sunday (2024-01-14)
	expectedStart := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, 1, 14, 23, 59, 59, 999999999, time.UTC)

	if !start.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, start)
	}

	if !end.Equal(expectedEnd) {
		t.Errorf("Expected end date %v, got %v", expectedEnd, end)
	}
}
