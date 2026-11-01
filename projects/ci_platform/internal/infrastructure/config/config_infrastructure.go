package infraConfig

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// server config
// snake_case keys in nested structs do not work if you use `json`,
// you need to use `mapstructure`
// https://github.com/spf13/viper/issues/125
type App struct {
	Server *Server       `mapstructure:"server"`
	DB     *DB           `mapstructure:"database"`
	Redis  *redis.Config `mapstructure:"redis"`
	Fiber  *fiber.Config `mapstructure:"fiber"`
	// have to keep it here, otherwise viper.UnmarshalExact will throw an error
	// because there is s CLI flag 'config', but the struct does not have a
	// corresponding field
	Path string `mapstructure:"config"`
}

var Default = App{
	Server: &Server{
		Host:         "127.0.0.1",
		Port:         "8080",
		RequestLimit: 100,
		JWT: &JWT{
			Secret:     "",
			Expiration: 72,
		},
	},
	DB: &DB{
		Type: DBTypePostgresql,
		Conn: map[string]string{
			"host":     "127.0.0.1",
			"port":     "9920",
			"user":     "admin",
			"password": "admin",
			"dbname":   "ci_platform",
			"sslmode":  "disable",
		},
	},
	Redis: &redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "admin",
		Database: 0,
		Reset:    false,
	},
	Fiber: &fiber.Config{},
}

type Server struct {
	// server port
	Port string `mapstructure:"port"`
	// server host
	Host string `mapstructure:"host"`
	// Max number of recent connections during `Expiration` seconds (1 minute)
	// before sending a 429 response
	RequestLimit int  `mapstructure:"request_limit"`
	JWT          *JWT `mapstructure:"jwt"`
}

type JWT struct {
	// secret signs jwt claims
	Secret string `mapstructure:"secret"`
	// expiration affects how fast jwt claims expire (hours)
	Expiration int `mapstructure:"expiration" validate:"required,min=1"`
}

// jwt middleware requires an array of bytes
func (c *JWT) GetSecret() []byte { return []byte(c.Secret) }

type DB struct {
	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	Conn map[string]string `mapstructure:"connection"`
	Type DBType            `mapstructure:"type"`
}

type DBType string

const DBTypePostgresql = "postgresql"

// populate the config file using command line, environment, and a set of
// predefined paths
func (a *App) Populate() *App {
	// enable automatic load of environment variables that match
	// e.g., TIANYI_CONFIG for '--config' flag
	viper.SetEnvPrefix("ci_platform")
	viper.AutomaticEnv()
	if fromCmd := viper.GetString("config"); fromCmd != "" {
		// get config path from cmd or from environment
		viper.SetConfigFile(fromCmd)
	} else {
		// default config paths
		homeDirectory, errHome := os.UserHomeDir()
		configDirectory, errConfig := os.UserConfigDir()
		cacheDirectory, errCache := os.UserCacheDir()
		cobra.CheckErr(errHome)
		cobra.CheckErr(errConfig)
		cobra.CheckErr(errCache)
		paths := []string{
			".",
			"ci_platform",
			filepath.Join(homeDirectory, "ci_platform"),
			filepath.Join(configDirectory, "ci_platform"),
			filepath.Join(cacheDirectory, "ci_platform"),
			"/opt/ci_platform",
		}
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
		viper.SetConfigType("yml")
		viper.SetConfigName("ci_platform")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config file: %v", err)
	}
	configPath := viper.ConfigFileUsed()
	if err := viper.UnmarshalExact(a); err != nil {
		log.Fatalf(
			"could not unmarshal config file '%s': %v",
			configPath,
			err,
		)
	}
	log.Println("Using config file:", configPath)
	return a
}
