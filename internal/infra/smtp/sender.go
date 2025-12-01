package smtp

type Connection struct {
	client *SMTPClient
}

func NewConnection(client *SMTPClient) *Connection {
	return &Connection{client: client}
}

func (c *Connection) Send(to, subject, body string) error {
	// TODO: Implementar client e lógica de envio de email de vdd
	return nil
}
