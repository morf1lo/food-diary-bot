package configs

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type TgBotConfig struct {
	Token string
	Debug bool
}
