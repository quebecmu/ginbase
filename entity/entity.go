package entity

type System struct {
	Application *Application `yaml:"application"`
	LogInfo     *LogInfo     `yaml:"logInfo"`
	RedisInfo   *RedisConfig `yaml:"redisInfo"`
	MysqlInfo   *MySQLConfig `yaml:"mysqlInfo"`
}

type Application struct {
	Name string `yaml:"name" json:"name"`
	Port int    `yaml:"port" json:"port"`
}

type LogInfo struct {
	Level      string `yaml:"level" json:"level"`
	Path       string `yaml:"path" json:"path"`
	MaxSize    int    `yaml:"maxSize" json:"maxSize"`
	MaxAge     int    `yaml:"maxAge" json:"maxAge"`
	MaxBackups int    `yaml:"maxBackups" json:"maxBackups"`
}

type MySQLConfig struct {
	Url string `yaml:"url" json:"url"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Network  string `yaml:"network" json:"network"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}
