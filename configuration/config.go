package configuration

type Config struct {
	VersionString      string
	DebugMode          bool
	DataPath           string
	CachePath          string
	FirstLevelDirAsTag bool
}

var config Config

func Init(c Config) {
	config = c
}

func Get() Config {
	return config
}
