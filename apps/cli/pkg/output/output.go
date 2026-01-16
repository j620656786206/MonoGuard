// Package output provides formatted output utilities
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Formatter handles output formatting
type Formatter struct {
	Format string // "text" or "json"
}

// NewFormatter creates a new output formatter
func NewFormatter(format string) *Formatter {
	// Normalize format
	normalized := strings.ToLower(format)
	if normalized == "" {
		normalized = "text"
	}
	return &Formatter{Format: normalized}
}

// Print outputs data in the configured format to stdout
func (f *Formatter) Print(data interface{}) error {
	return f.PrintTo(os.Stdout, data)
}

// PrintTo outputs data in the configured format to the given writer
func (f *Formatter) PrintTo(w io.Writer, data interface{}) error {
	if f.Format == "json" {
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(b))
	} else {
		// For text format, provide structured output
		switch v := data.(type) {
		case string:
			fmt.Fprintln(w, v)
		case map[string]interface{}:
			for key, val := range v {
				fmt.Fprintf(w, "%s: %v\n", capitalize(key), val)
			}
		case map[string]string:
			for key, val := range v {
				fmt.Fprintf(w, "%s: %s\n", capitalize(key), val)
			}
		default:
			// Use reflection for structs
			val := reflect.ValueOf(data)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.Kind() == reflect.Struct {
				typ := val.Type()
				for i := 0; i < val.NumField(); i++ {
					field := typ.Field(i)
					value := val.Field(i)
					fmt.Fprintf(w, "%s: %v\n", field.Name, value.Interface())
				}
			} else {
				fmt.Fprintf(w, "%+v\n", data)
			}
		}
	}
	return nil
}

// capitalize returns the string with first letter capitalized
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
