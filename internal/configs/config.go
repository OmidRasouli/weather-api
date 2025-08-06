// config.go
// Contains core configuration structures

package configs // Change to singular for package name - more idiomatic Go

var global *Config // Lowercase for package-level variable

// Config is the main configuration struct that holds all app settings.
// It is populated from the environment variables.
type Config struct {
	Server      ServerConfig      // Server-related settings (e.g., port)
	Database    DatabaseConfig    // Database connection settings
	OpenWeather OpenWeatherConfig // OpenWeather API settings
	Redis       RedisConfig       // Add Redis configuration
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port int // Port on which the HTTP server will listen
}

// DatabaseConfig holds the configuration for connecting to a database.
// It includes host, port, user credentials, database name, and SSL mode.
type DatabaseConfig struct {
	Host     string // Database host address
	Port     int    // Database port
	User     string // Username for database authentication
	Password string // Password for database authentication
	DBName   string // Name of the database to connect to
	SSLMode  string // SSL mode (e.g., disable, require)
}

// OpenWeatherConfig holds the configuration for accessing the OpenWeather API.
// It includes the API key required for authentication.
type OpenWeatherConfig struct {
	APIKey string // API key for OpenWeather service
}

// RedisConfig holds configuration for Redis cache
type RedisConfig struct {
	Host     string // Redis server host
	Port     int    // Redis server port
	Password string // Password for Redis authentication
	DB       int    // Redis database number
	TTL      int    // Time-to-live in seconds for cached items
}

// Access methods - in idiomatic Go, we'd typically not use "Get" prefixes
func (c *Config) GetServer() ServerConfig {
	return c.Server
}

func (c *Config) GetDatabase() DatabaseConfig {
	return c.Database
}

func (c *Config) GetOpenWeather() OpenWeatherConfig {
	return c.OpenWeather
}

func (c *Config) GetRedis() RedisConfig {
	return c.Redis
}

// Global returns the singleton config instance
func Global() *Config {
	return global
}

// SetGlobal sets the global configuration instance
func SetGlobal(cfg *Config) {
	global = cfg
}
