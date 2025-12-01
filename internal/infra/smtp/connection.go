package smtp

type SMTPClient struct {
	Server   string
	Port     int
	Username string
	Password string
}

func SMTPConnection(server string, port int, username string, password string) (*Connection, error) {
	client := &SMTPClient{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}

	return NewConnection(client), nil
}
