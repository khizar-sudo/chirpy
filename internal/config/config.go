package config

import (
	"sync/atomic"

	"github.com/khizar-sudo/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}
