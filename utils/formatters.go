package utils

import "strings"

func StandardizeCuit(cuit string) string {
	if len(cuit) > 11 {
		cuit = strings.ReplaceAll(cuit, "-", "")
	}
	return cuit
}
