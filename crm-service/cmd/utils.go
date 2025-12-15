package main

import (
	"fmt"
	"strings"
)

func generateCode(prefix string, model interface{}) string {
	var count int64
	db.Model(model).Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1)
}

func formatPrice(price float64) string {
	priceStr := fmt.Sprintf("%.0f", price)
	n := len(priceStr)
	if n <= 3 {
		return priceStr
	}

	var result strings.Builder
	for i, c := range priceStr {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(c)
	}
	return result.String()
}
