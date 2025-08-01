package cmd

type Config struct {
	HttpPort            string
	DbHost              string
	DbPort              string
	DbUser              string
	DbPassword          string
	DbName              string
	DbSslMode           string
	EventGoroutineLimit int
}
