package config

const LocalEnv = "local"
const ProdEnv = "prod"

var allEnvs map[string]struct{} = map[string]struct{}{
	LocalEnv: {},
	ProdEnv:  {},
}
