package chat

//Error - default chat error
type Error struct {
	s string
}

//NewError init new chat error
func NewError(message string) Error {
	return Error{
		s: message,
	}
}

func (err Error) Error() string {
	return err.s
}
