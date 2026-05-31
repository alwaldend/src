package fp

func Return[T1 any](t1 T1) T1 {
	return t1
}

func ReturnE[T1 any](t1 T1) (T1, error) {
	return t1, nil
}

func Return2[T1, T2 any](t1 T1, t2 T2) (T1, T2) {
	return t1, t2
}
