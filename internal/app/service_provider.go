package app

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/spv-dev/platform_common/pkg/closer"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/spv-dev/platform_common/pkg/db/pg"
	"github.com/spv-dev/platform_common/pkg/db/transaction"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/client/cache"

	redisClient "github.com/spv-dev/auth/internal/client/cache/redis"
	"github.com/spv-dev/auth/internal/config"
	"github.com/spv-dev/auth/internal/repository"
	cacheRepository "github.com/spv-dev/auth/internal/repository/cache"
	userRepository "github.com/spv-dev/auth/internal/repository/user"
	"github.com/spv-dev/auth/internal/service"
	userService "github.com/spv-dev/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userCache   repository.UserCache
	redisPool   *redis.Pool
	redisClient cache.RedisClient

	userService service.UserService

	userServer *user.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig получение конфигурации подключения к postgres
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

// GRPCConfig получение конфигурации подключения gRPC
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

// RedisConfig получение конфигурации подключения к redis
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %v", err)
		}

		s.redisConfig = cfg
	}
	return s.redisConfig
}

// DBClient получение подключения к БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database : %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager получение объекта менеджера транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

// UserRepository получение объекта доступа к слою repo
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService получение объекта доступа к сервисному слою
func (s *serviceProvider) RedisPool() *redis.Pool {
	if s.redisPool == nil {
		s.redisPool = &redis.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

// RedisClient получение объекта доступа к кэшу
func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redisClient.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// RedisClient получение объекта доступа к кэшу
func (s *serviceProvider) UserCache() repository.UserCache {
	if s.userCache == nil {
		s.userCache = cacheRepository.NewCache(s.RedisClient())
	}

	return s.userCache
}

// UserService получение объекта доступа к сервисному слою
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.UserCache(),
		)
	}

	return s.userService
}

// UserServer получение объекта сервиса
func (s *serviceProvider) UserServer(ctx context.Context) *user.Server {
	if s.userServer == nil {
		s.userServer = user.NewServer(s.UserService(ctx))
	}

	return s.userServer
}
