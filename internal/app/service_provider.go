package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/gomodule/redigo/redis"
	"github.com/spv-dev/platform_common/pkg/closer"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/spv-dev/platform_common/pkg/db/pg"
	"github.com/spv-dev/platform_common/pkg/db/transaction"

	"github.com/spv-dev/auth/internal/api/access"
	"github.com/spv-dev/auth/internal/api/auth"
	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/client/cache"
	redisClient "github.com/spv-dev/auth/internal/client/cache/redis"
	"github.com/spv-dev/auth/internal/client/kafka"
	kafkaProducer "github.com/spv-dev/auth/internal/client/kafka/producer"
	"github.com/spv-dev/auth/internal/config"
	"github.com/spv-dev/auth/internal/repository"
	cacheRepository "github.com/spv-dev/auth/internal/repository/cache"
	userRepository "github.com/spv-dev/auth/internal/repository/user"
	"github.com/spv-dev/auth/internal/service"
	userService "github.com/spv-dev/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	redisConfig         config.RedisConfig
	httpConfig          config.HTTPConfig
	swaggerConfig       config.SwaggerConfig
	kafkaProducerConfig config.KafkaProducerConfig
	authConfig          config.AuthConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userCache   repository.UserCache
	redisPool   *redis.Pool
	redisClient cache.RedisClient

	producer kafka.Producer
	sender   sarama.SyncProducer

	userService service.UserService

	userServer *user.Server

	authServer   *auth.Server
	accessServer *access.Server
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

// HTTPConfig получение конфигурации подключения http
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = cfg
	}
	return s.httpConfig
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

// SwaggerConfig получение конфигурации подключения к redis
func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = cfg
	}
	return s.swaggerConfig
}

// AuthConfig получение конфигурации подключения Auth
func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config: %v", err)
		}

		s.authConfig = cfg
	}
	return s.authConfig
}

func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka producer config: %s", err.Error())
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
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

// RedisPool получение объекта доступа к сервисному слою
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
			s.Producer(),
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

// AuthServer получение объекта сервиса
func (s *serviceProvider) AuthServer(_ context.Context) *auth.Server {
	if s.authServer == nil {
		s.authServer = auth.NewServer(s.AuthConfig())
	}

	return s.authServer
}

// AuthServer получение объекта сервиса
func (s *serviceProvider) AccessServer(_ context.Context) *access.Server {
	if s.accessServer == nil {
		s.accessServer = access.NewServer()
	}

	return s.accessServer
}

func (s *serviceProvider) Producer() kafka.Producer {
	if s.producer == nil {
		s.producer = kafkaProducer.NewProducer(
			s.Sender(),
		)
		closer.Add(s.producer.Close)
	}

	return s.producer
}

func (s *serviceProvider) Sender() sarama.SyncProducer {
	if s.sender == nil {
		sender, err := sarama.NewSyncProducer(
			s.KafkaProducerConfig().Brokers(),
			s.KafkaProducerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create kafka sender: %v", err)
		}

		s.sender = sender
	}

	return s.sender
}
