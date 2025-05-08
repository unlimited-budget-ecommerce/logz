package logz

import (
	"maps"
	"strings"
)

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

var replacerMap = make(map[string]func(string) string)

// SetReplacerMap should be called before calling [logz.MaskMap]. Keys are case insensitive.
//
// **This function is unsafe for concurrent calls.**
func SetReplacerMap(m map[string]func(string) string) {
	for k, v := range m {
		replacerMap[strings.ToLower(k)] = v
	}
}

// MaskMap masks field (case insensitive) based on replacerMap.
// To set replacerMap, calls [logz.SetReplacerMap].
func MaskMap(m map[string]any) map[string]any {
	newMap := maps.Clone(m)
	for k, v := range newMap {
		switch v := v.(type) {
		case string:
			if fn, ok := replacerMap[strings.ToLower(k)]; ok {
				newMap[k] = fn(v)
			}
		case map[string]any:
			newMap[k] = MaskMap(v)
		case []any:
			for i := range v {
				vMap, ok := v[i].(map[string]any)
				if !ok {
					break // assuming all items has the same type
				}
				v[i] = MaskMap(vMap)
			}
		}
	}

	return newMap
}
