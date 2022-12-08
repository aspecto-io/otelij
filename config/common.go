package config

const (
	TYPE = "TYPE" //span, metrics, logs (For future usage)
)

func GetType() Type {
	return TRACES
}
