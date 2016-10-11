package gomu

// Error encapsulates a name, an error and whether there is a custom error message or not.
type Error struct {
	Name                     string
	Err                      error
	CustomErrorMessageExists bool
}

func (e Error) Error() string {
	if e.CustomErrorMessageExists {
		return e.Err.Error()
	}
	return e.Name + ": " + e.Err.Error()
}

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Errors returns itself.
func (e Errors) Errors() []error {
	return e
}

func (e Errors) Error() (errStr string) {
	for _, err := range e {
		errStr += err.Error() + ";"
	}
	return
}
