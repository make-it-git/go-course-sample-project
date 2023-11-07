package config

import (
	"fmt"
	"strings"
)

type UnknownEnvErr struct {
	env       string
	supported string
}

func NewUnknownEnvErr(env string) error {
	keys := make([]string, len(allEnvs))
	i := 0
	for k, _ := range allEnvs {
		keys[i] = k
		i++
	}

	return UnknownEnvErr{
		env:       env,
		supported: strings.Join(keys, ", "),
	}
}

func (e UnknownEnvErr) Error() string {
	return fmt.Sprintf("Unknown env: %s, supported envs: %s", e.env, e.supported)
}
