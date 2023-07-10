package validator

type EmptyJWTError struct {
}

func (e EmptyJWTError) Error() string {
	return "Empty JWT"
}
