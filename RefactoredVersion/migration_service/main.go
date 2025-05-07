package main

import (
	"context"
	"github.com/SwanHtetAungPhyo/migration_service/providers/config"
	"github.com/SwanHtetAungPhyo/migration_service/providers/database"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			func(cfg *config.Config) string {
				return cfg.DB.DSN
			},
			database.NewDB,
		),
		fx.Invoke(func(db *gorm.DB, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Println("Starting database migration...")
					return database.Migrate(db)
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Closing database connection...")
					sqlDB, err := db.DB()
					if err != nil {
						return err
					}
					return sqlDB.Close()
				},
			})
		}),
	)

	app.Run()
}
