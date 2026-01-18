package infraLifecycle

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectDatabase(config *infraConfig.DB) *gorm.DB {
	var dialector gorm.Dialector
	conn := parseConnection(config)
	switch config.Type {
	case infraConfig.DBTypePostgresql:
		dialector = postgres.Open(conn)
	default:
		log.Panicf("database type '%s' is not supported", config.Type)
	}

	databaseLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			// Slow SQL threshold
			SlowThreshold: time.Millisecond * 10,
			// Log level
			LogLevel: logger.Info,
			// Ignore ErrRecordNotFound error for logger
			IgnoreRecordNotFoundError: false,
			// enable color
			Colorful: true,
		},
	)
	database, err := gorm.Open(
		dialector,
		&gorm.Config{
			Logger: databaseLogger,
		})
	if err != nil {
		log.Fatalf("failed to connect to database")
	}
	return database
}

func parseConnection(config *infraConfig.DB) string {
	buffer := &bytes.Buffer{}
	format := getFormat(config.Type)
	for key, value := range config.Conn {
		fmt.Fprintf(buffer, format, key, value)
	}
	return buffer.String()
}

func getFormat(dbType infraConfig.DBType) string {
	switch dbType {
	case infraConfig.DBTypePostgresql:
		return "%s=%s "
	}
	return ""
}
