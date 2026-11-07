package core

import (
	"testing"

	"core/utils"
)

func TestFormatINR(t *testing.T) {
	tests := []struct {
		paise    int64
		expected string
	}{
		{0, "₹0"},
		{100, "₹1"},
		{1000000, "₹10,000"},
		{12345678, "₹1,23,456.78"},
		{-50000, "-₹500"},
	}

	for _, tc := range tests {
		got := utils.FormatINR(tc.paise)
		if got != tc.expected {
			t.Errorf("FormatINR(%d) = %s; want %s", tc.paise, got, tc.expected)
		}
	}
}

func TestParseINRToPaise(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"10000", 1000000},
		{"₹10,000", 1000000},
		{"₹1,23,456.78", 12345678},
		{"50.5", 5050},
	}

	for _, tc := range tests {
		got, err := utils.ParseINRToPaise(tc.input)
		if err != nil {
			t.Fatalf("ParseINRToPaise(%s) error: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ParseINRToPaise(%s) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}

func TestCalculatePercentageMinor(t *testing.T) {
	// 10% of ₹10,000 (1,000,000 paise) = ₹1,000 (100,000 paise)
	// 10% = 1000 basis points
	got := utils.CalculatePercentageMinor(1000000, 1000)
	if got != 100000 {
		t.Errorf("CalculatePercentageMinor = %d; want %d", got, 100000)
	}
}
