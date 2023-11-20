package cfg

type Cfg struct {
	RedisAddress string
	ServerPort   int
}

func LoadConfig() Cfg {
	cfg := Cfg{
		RedisAddress: "localhost:6379",
		ServerPort:   6000,
	}
	return cfg
}
