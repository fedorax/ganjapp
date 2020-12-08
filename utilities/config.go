package utilities

// Config struct
type Config struct {
	EnvLoaded   bool
	CookieKey   string
	JWTKey      string
	JWTAudience string
}

// AppConfig contains the configuration struct
var AppConfig = Config{EnvLoaded: LoadEnv(), CookieKey: getCookieKey(), JWTKey: getJWTKey(), JWTAudience: getJWTAudience()}

func getCookieKey() string {
	// Try and use a key from the environment, otherwise generate one...
	// Note that this is not ideal, as each time the server restarts,
	// all the sessions will be invalidated!
	return GetEnv("GANJAPP_COOKIE_KEY", GetRandomString(40))
}

func getJWTKey() string {
	// Try and use a key from the environment, otherwise generate one...
	// Note that this is not ideal, as each time the server restarts,
	// all the sessions will be invalidated!
	return GetEnv("GANJAPP_JWT_KEY", GetRandomString(40))
}

func getJWTAudience() string {
	// Try and use a key from the environment, otherwise generate one...
	// Note that this is not ideal, as each time the server restarts,
	// all the sessions will be invalidated!
	return GetEnv("GANJAPP_JWT_AUDIENCE", GetRandomString(20)+".servers.ganj.app")
}
