package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func parseEnvInt(envMap map[string]string, key string, errorMsg string) int {
	val, err := strconv.Atoi(envMap[key])
	if err != nil {
		log.Fatalf("%s: %v", errorMsg, err)
	}
	return val
}

func parseEnvDuration(envMap map[string]string, key string, errorMsg string) time.Duration {
	seconds := parseEnvInt(envMap, key, errorMsg)
	return time.Duration(int64(seconds) * int64(math.Pow10(9)))
}

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load env failed: %v", err)
	}

	return &config{
		app: &app{
			host:         envMap["APP_HOST"],
			port:         parseEnvInt(envMap, "APP_PORT", "load port failed"),
			name:         envMap["APP_NAME"],
			version:      envMap["APP_VERSION"],
			readTimeout:  parseEnvDuration(envMap, "APP_READ_TIMEOUT", "load read timeout failed"),
			writeTimeout: parseEnvDuration(envMap, "APP_WRITE_TIMEOUT", "load write timeout failed"),
			bodyLimit:    parseEnvInt(envMap, "APP_BODY_LIMIT", "load body limit failed"),
		},
		db: &db{
			host:        envMap["DB_HOST"],
			port:        parseEnvInt(envMap, "DB_PORT", "load database port failed"),
			username:    envMap["DB_USER"],
			password:    envMap["DB_PASSWORD"],
			database:    envMap["DB_NAME"],
			maxPoolSize: parseEnvInt(envMap, "DB_MAX_POOL_SIZE", "load max pool size failed"),
		},
		jwt: &jwt{
			secretKey:       envMap["JWT_SECRET_KEY"],
			accessExpiresAt: parseEnvInt(envMap, "JWT_ACCESS_EXPIRES", "load access expires at failed"),
		},
	}
}

type IConfig interface {
	App() IAppConfig
	DB() IDBConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	Host() string
	Port() int
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
}

func (c *config) App() IAppConfig {
	return c.app
}

func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }
func (a *app) BodyLimit() int              { return a.bodyLimit }

type IDBConfig interface {
	Url() string
	MaxPoolSize() int
	Host() string
	Port() int
	Database() string
	Username() string
	Password() string
}

type db struct {
	host        string
	port        int
	username    string
	password    string
	database    string
	maxPoolSize int
}

func (c *config) DB() IDBConfig {
	return c.db
}

func (d *db) Url() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin&maxPoolSize=%d",
		d.username,
		d.password,
		d.host,
		d.port,
		d.database,
		d.maxPoolSize)
}
func (d *db) MaxPoolSize() int { return d.maxPoolSize }
func (d *db) Host() string     { return d.host }
func (d *db) Port() int        { return d.port }
func (d *db) Database() string { return d.database }
func (d *db) Username() string { return d.username }
func (d *db) Password() string { return d.password }

type IJwtConfig interface {
	SecretKey() []byte
	AccessExpiresAt() int
	SetJwtAccessExpires(t int)
}

type jwt struct {
	secretKey       string
	accessExpiresAt int
}

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

func (j *jwt) SecretKey() []byte         { return []byte(j.secretKey) }
func (j *jwt) AccessExpiresAt() int      { return j.accessExpiresAt }
func (j *jwt) SetJwtAccessExpires(t int) { j.accessExpiresAt = t }
