package deps

import (
	"llmops/internal/apiserver/store/mysql"
	"llmops/internal/apiserver/store/redis"
)

// Dependencies contains clients and stores shared by apiserver layers.
type Dependencies struct {
	MySQL mysql.Factory
	Redis redis.RStore
}
