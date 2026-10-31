package fp

func Return[T1 any](v1 T1) T1 {
	return v1
}

func ReturnI[T1, T2 any](v1 T2) func(T1) T2 {
	return func(_ T1) T2 {
		return v1
	}
}

func ReturnE[T1 any](v1 T1) (T1, error) {
	return v1, nil
}

func Return2[T1, T2 any](v1 T1, v2 T2) (T1, T2) {
	return v1, v2
}
