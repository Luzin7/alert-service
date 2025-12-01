package smtp

type SMTPClient struct {
	Server   string
	Port     int
	Username string
	Password string
}

func SMTPConnection(server string, port int, username string, password string) (*SMTPClient, error) {
	connection := &SMTPClient{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}

	return connection, nil
}
