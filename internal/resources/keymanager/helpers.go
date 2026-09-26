package keymanager

import "strings"

func GetReferenceID(ref string) string {
	ref = strings.TrimSuffix(ref, "/")
	if i := strings.LastIndex(ref, "/"); i >= 0 {
		return ref[i+1:]
	}
	return ref
}
