package fp

type EmptyEither = Either[struct{}]

func EmptyEitherOf(err error) EmptyEither {
	if err == nil {
		return Right(struct{}{})
	}
	return Left[struct{}](err)
}

func EmptyLeft(err error) EmptyEither {
	return Left[struct{}](err)
}

func EmptyRight() EmptyEither {
	return Right(struct{}{})
}
