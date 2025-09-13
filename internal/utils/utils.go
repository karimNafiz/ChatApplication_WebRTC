package utils

// need a function for doing intersection operation between two hash maps

// very simple intersection function
func Intersection(set_a map[string]struct{}, set_b map[string]struct{}) map[string]struct{} {
	intersection := make(map[string]struct{})
	search_set := set_a

	if len(set_a) >= len(set_b) {
		search_set = set_b
	}

	for key, _ := range search_set {

		if _, ok := search_set[key]; !ok {
			continue
		}
		intersection[key] = struct{}{}
	}
	return intersection
}
