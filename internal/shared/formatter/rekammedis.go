package formatter

import "strings"

func FormatNoRekamMedis(noRM string) string {
	noRM = strings.TrimSpace(noRM)
	if len(noRM) == 8 && strings.HasPrefix(noRM, "00") {
		return noRM[2:]
	}
	return noRM
}
