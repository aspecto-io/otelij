package config

type ParameterError struct {
	error string
}

func (pe ParameterError) Error() string {
	return pe.error
}
