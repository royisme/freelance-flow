package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// FormatAmount formats a float amount with the currency code.
func FormatAmount(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

// FormatDate formats a date string according to the user's settings.
func FormatDate(raw string, dateFormat string, timezone string) string {
	if raw == "" {
		return ""
	}
	layouts := []string{"2006-01-02", time.RFC3339}
	var parsed time.Time
	var err error
	for _, l := range layouts {
		parsed, err = time.Parse(l, raw)
		if err == nil {
			break
		}
	}
	if err != nil {
		return raw
	}
	// default timezone if empty
	if timezone == "" {
		timezone = "UTC"
	}
	loc, locErr := time.LoadLocation(timezone)
	if locErr == nil {
		parsed = parsed.In(loc)
	}
	if dateFormat == "" {
		dateFormat = "2006-01-02"
	}
	return parsed.Format(dateFormat)
}

// GetCurrencySymbol returns the currency symbol for a given currency code.
func GetCurrencySymbol(currency string) string {
	symbols := map[string]string{
		"CAD": "$",
		"USD": "$",
		"CNY": "¥",
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
	}
	if sym, ok := symbols[currency]; ok {
		return sym
	}
	return "$" // Default to $
}

// FormatServiceType converts a snake_case service type to Title Case.
func FormatServiceType(serviceType string) string {
	if serviceType == "" {
		return ""
	}
	description := strings.ReplaceAll(serviceType, "_", " ")
	words := strings.Fields(description)
	for i, w := range words {
		if len(w) > 0 {
			// Capitalize first letter
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// ApplyTemplate replaces placeholders in a template string with values from the data struct.
// It supports placeholders like {{key}} matching struct fields case-insensitively or by json tag.
//go:noinline
func ApplyTemplate(tmplStr string, data interface{}) string {
	if tmplStr == "" {
		return ""
	}

	// 1. Build a map of keys to values from data
	val := reflect.ValueOf(data)
	//nolint:govet // reflect.Ptr comparison triggers inline rule incorrectly
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		// If not a struct, can't easily map keys. Return as is or support map?
		// For now assume struct as used in InvoiceService
		return tmplStr
	}

	dataMap := make(map[string]string)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Use JSON tag as key if available, otherwise lowercase field name
		key := strings.ToLower(field.Name)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				key = parts[0]
			}
		}

		// Store value as string
		dataMap[key] = fmt.Sprintf("%v", fieldVal.Interface())
		// Also store strict lowercase name just in case
		dataMap[strings.ToLower(field.Name)] = fmt.Sprintf("%v", fieldVal.Interface())
	}

	// 2. Regex replace {{key}}
	// Matches {{ key }} or {{key}}
	re := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}\}`)
	return re.ReplaceAllStringFunc(tmplStr, func(match string) string {
		// Extract key
		// match is like "{{number}}"
		parts := re.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}
		key := strings.ToLower(parts[1])

		if val, ok := dataMap[key]; ok {
			return val
		}
		return match // Keep original if not found
	})
}
