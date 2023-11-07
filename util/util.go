package util

func Map[P, Q any](fpq func(P) Q, ps []P) []Q { // Map implementation from https://stackoverflow.com/a/72498530
	result := make([]Q, len(ps))
	for i, p := range ps {
		result[i] = fpq(p)
	}
	return result
}

func Mapi[P, Q any](fpiq func(P, int) Q, ps []P) []Q { // Calls fpiq(p, i) where i is the index of p in ps
	result := make([]Q, len(ps))
	for i, p := range ps {
		result[i] = fpiq(p, i)
	}
	return result
}

func Filter[P any](predicate func(P) bool, ps []P) []P {
	result := make([]P, 0)
	for _, p := range ps {
		if predicate(p) {
			result = append(result, p)
		}
	}
	return result
}

func Filteri[P any](predicate func(P, int) bool, ps []P) []P {
	result := make([]P, 0)
	for i, p := range ps {
		if predicate(p, i) {
			result = append(result, p)
		}
	}
	return result
}
