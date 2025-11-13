package utils

import (
	"encoding/csv"
	"strings"
)

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

// WriteAsCSV writes a slice of strings as a CSV string.
func WriteAsCSV(vals []string) (string, error) {
	b := &strings.Builder{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}
