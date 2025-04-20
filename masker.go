package logz

import "strings"

func Mask(s string) string {
	return "****"
}

func MaskName(s string) string {
	names := strings.Split(s, " ")
	for i, name := range names {
		if len(name) > 2 {
			names[i] = string(name[0]) + strings.Repeat("*", len(name)-2) + string(name[len(name)-1])
		} else {
			names[i] = strings.Repeat("*", len(name))
		}
	}

	return strings.Join(names, " ")
}

func MaskEmail(s string) string {
	i := strings.Index(s, "@")
	if i == -1 {
		return s
	}

	name, domain := s[:i], s[i:]
	if len(name) > 2 {
		name = string(name[0]) + strings.Repeat("*", len(name)-2) + string(name[len(name)-1])
	} else {
		name = strings.Repeat("*", len(name))
	}

	return name + domain
}
