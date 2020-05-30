package utils

type void struct{}

func Missing(a, b []string) []string {
	ma := make(map[string]void, len(a))
	diffs := []string{}
	for _, ka := range a {
		ma[ka] = void{}
	}
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}
