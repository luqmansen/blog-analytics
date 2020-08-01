package configs

type Configuration struct {
	Server   Server
	Database Database
	Client   []string
}

type Database struct {
	URI      string
	Database string
	Timeout  int
}

type Server struct {
	Port string
}

type Client struct {
	Url []string
}
