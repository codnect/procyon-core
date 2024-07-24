package component

type NotFoundError struct {
	Name string
}

func (ne *NotFoundError) Error() string {
	return ""
}
