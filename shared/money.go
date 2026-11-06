package shared

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func FormatINR(paise int64) string {
	sign := ""
	if paise < 0 {
		sign = "-"
		paise = -paise
	}

	rupees := paise / 100
	remainder := paise % 100

	rupeesStr := strconv.FormatInt(rupees, 10)
	n := len(rupeesStr)

	var formatted string
	if n <= 3 {
		formatted = rupeesStr
	} else {
		last3 := rupeesStr[n-3:]
		rest := rupeesStr[:n-3]
		var parts []string
		for len(rest) > 2 {
			parts = append([]string{rest[len(rest)-2:]}, parts...)
			rest = rest[:len(rest)-2]
		}
		if len(rest) > 0 {
			parts = append([]string{rest}, parts...)
		}
		parts = append(parts, last3)
		formatted = strings.Join(parts, ",")
	}

	if remainder == 0 {
		return fmt.Sprintf("%s₹%s", sign, formatted)
	}
	return fmt.Sprintf("%s₹%s.%02d", sign, formatted, remainder)
}

func ParseINRToPaise(val string) (int64, error) {
	clean := strings.TrimSpace(val)
	clean = strings.ReplaceAll(clean, "₹", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	if clean == "" {
		return 0, nil
	}

	parts := strings.Split(clean, ".")
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid whole amount: %w", err)
	}

	var paise int64
	if len(parts) > 1 {
		fractionStr := parts[1]
		if len(fractionStr) == 1 {
			fractionStr += "0"
		} else if len(fractionStr) > 2 {
			fractionStr = fractionStr[:2]
		}
		f, err := strconv.ParseInt(fractionStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid decimal: %w", err)
		}
		paise = f
	}

	if whole < 0 {
		return whole*100 - paise, nil
	}
	return whole*100 + paise, nil
}

func CalculatePercentageMinor(baseMinor int64, basisPoints int64) int64 {
	val := float64(baseMinor*basisPoints) / 10000.0
	return int64(math.Round(val))
}
