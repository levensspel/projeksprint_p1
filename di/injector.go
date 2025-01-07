package di

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/database"
	"github.com/levensspel/go-gin-template/domain"
	"github.com/levensspel/go-gin-template/infrastructure/storage"

	userRepository "github.com/levensspel/go-gin-template/repository/user"

	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	// Setup database connection
	do.Provide[*pgxpool.Pool](Injector, database.NewUserRepositoryInject)

	// Setup repositories
	// UserRepository
	do.Provide[userRepository.UserRepository](Injector, userRepository.NewUserRepositoryInject)

	// Setup client
	envMode := os.Getenv("MODE")
	if envMode == "DEBUG" {
		do.Provide[domain.StorageClient](Injector, storage.NewMockStorageClientInject)
	} else {
		do.Provide[domain.StorageClient](Injector, storage.NewS3StorageClientInject)
	}
}
