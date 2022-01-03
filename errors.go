package cfg

type sentinelError string

func (s sentinelError) Error() string {
	return string(s)
}

const NoValueProvidedError sentinelError = "no value provided"
