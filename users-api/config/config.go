package config

import "time"

const (
	MySQLHost     = "localhost"
	MySQLPort     = "3306"
	MySQLDatabase = "users-api"
	MySQLUsername = "root"
	MySQLPassword = "root"
	CacheDuration = 30 * time.Second
	MemcachedHost = "localhost"
	MemcachedPort = "11211"
	JWTKey        = "ThisIsAnExampleJWTKey!"
	JWTDuration   = 24 * time.Hour
)
