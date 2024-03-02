package types

//connect to local Chroma

type ServerConfig struct {
	Host string
	Port int
}

// Add server configuration

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "localhost",
		Port: 8080,
	}
}
