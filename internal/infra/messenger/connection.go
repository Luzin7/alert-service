package messenger

import "github.com/rabbitmq/amqp091-go"

func MessengerConnection(username string, password string, host string, port string) (*amqp091.Connection, error) {
	connection, err := amqp091.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port + "/")

	return connection, err
}
