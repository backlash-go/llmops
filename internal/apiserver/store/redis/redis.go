package redis

import (
	"crypto/tls"
	"strconv"
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v7"
)

var (
	once    sync.Once
	S       *datastore
	initErr error
)

// Config defines redis connection options.
type Config struct {
	Host                  string
	Port                  int
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

// RStore defines redis store behavior.
type RStore interface {
	Rdb() goredis.UniversalClient
	Close() error
}

type datastore struct {
	db goredis.UniversalClient
}

// NewRedisStore creates a redis store singleton.
func NewRedisStore(config *Config) (RStore, error) {
	once.Do(func() {
		store := &datastore{db: newClient(config)}
		if err := store.db.Ping().Err(); err != nil {
			initErr = err
			_ = store.Close()

			return
		}

		S = store
	})

	return S, initErr
}

// Rdb returns the redis client.
func (d *datastore) Rdb() goredis.UniversalClient {
	return d.db
}

// Close closes the redis client.
func (d *datastore) Close() error {
	if d == nil || d.db == nil {
		return nil
	}

	return d.db.Close()
}

func newClient(config *Config) goredis.UniversalClient {
	if config == nil {
		config = &Config{}
	}

	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	opts := &goredis.UniversalOptions{
		Addrs:        redisAddrs(config),
		MasterName:   config.MasterName,
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}

	if config.MaxIdle > 0 {
		minIdleConns := config.MaxIdle
		if minIdleConns > poolSize {
			minIdleConns = poolSize
		}
		opts.MinIdleConns = minIdleConns
	}

	if opts.MasterName != "" {
		return goredis.NewFailoverClient(opts.Failover())
	}
	if config.EnableCluster {
		return goredis.NewClusterClient(opts.Cluster())
	}

	return goredis.NewClient(opts.Simple())
}

func redisAddrs(config *Config) []string {
	if len(config.Addrs) != 0 {
		return config.Addrs
	}
	if config.Port != 0 {
		return []string{config.Host + ":" + strconv.Itoa(config.Port)}
	}

	return nil
}
