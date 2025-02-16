package base

import (
	"os"

	itbasisCoreEnv "github.com/itbasis/go-tools-core/env"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
)

const (
	EnvSdkmOriginPrefix = "SDKM_ORIGIN_"
)

func (receiver *basePlugin) PrepareEnvironment(overrideEnv itbasisCoreEnv.Map, saveEnvKeys ...string) itbasisCoreEnv.Map {
	var result = itbasisCoreEnv.SlicesToMap(os.Environ())

	for _, key := range append(saveEnvKeys, itbasisCoreEnv.KeyPath) {
		var (
			sdkmOriginKey     = EnvSdkmOriginPrefix + key
			_, existOriginKey = result[sdkmOriginKey]
			value, exist      = result[key]
		)

		if !existOriginKey && exist {
			result[sdkmOriginKey] = value
		}
	}

	result = itbasisCoreEnv.MergeEnvs(result, overrideEnv)
	result[itbasisCoreEnv.KeyPath] = itbasisCoreOs.CleanPath(result[itbasisCoreEnv.KeyPath], itbasisCoreOs.ExecutableDir())

	return result
}
