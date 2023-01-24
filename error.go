package cfg

type sentinelError string

func (s sentinelError) Error() string {
	return string(s)
}

// A NoValueProvidedError denotes that no value was provided.
const NoValueProvidedError sentinelError = "no value provided"
