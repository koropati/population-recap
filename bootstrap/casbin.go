package bootstrap

import (
	"log"
	"sync"

	"github.com/casbin/casbin"
)

var (
	casbinOnce     sync.Once
	casbinEnforcer *casbin.Enforcer
)

// NewRedisClient creates a new Redis client connection
func NewCasbinEnforcer(config *Config) *casbin.Enforcer {
	casbinOnce.Do(func() {

		authEnforcer, err := casbin.NewEnforcerSafe(config.CasbinModelPath, config.CasbinPolicyPath)
		if err != nil {
			log.Fatal(err)
		}
		casbinEnforcer = authEnforcer
	})

	return casbinEnforcer
}
