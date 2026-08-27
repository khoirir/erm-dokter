package formatter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func FormatHarga(val any) string {
	var harga float64
	switch v := val.(type) {
	case float64:
		harga = v
	case float32:
		harga = float64(v)
	case int:
		harga = float64(v)
	case int64:
		harga = float64(v)
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return v
		}
		harga = f
	default:
		return ""
	}

	isNegative := harga < 0
	if isNegative {
		harga = -harga
	}

	intPart := int64(math.Floor(harga))
	fracPart := int64(math.Round((harga - float64(intPart)) * 100))
	if fracPart == 100 {
		intPart++
		fracPart = 0
	}

	strInt := fmt.Sprintf("%d", intPart)
	var formattedInt strings.Builder
	n := len(strInt)

	for i, c := range strInt {
		if i > 0 && (n-i)%3 == 0 {
			formattedInt.WriteRune('.')
		}
		formattedInt.WriteRune(c)
	}

	result := fmt.Sprintf("%s,%02d", formattedInt.String(), fracPart)
	if isNegative {
		return "-Rp " + result
	}
	return "Rp " + result
}
