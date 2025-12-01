package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	conn    *amqp.Connection
	handler *Handler
}

func NewWorker(conn *amqp.Connection, handler *Handler) *Worker {
	return &Worker{
		conn:    conn,
		handler: handler,
	}
}

func (w *Worker) Start(queueName string) {
	ch, _ := w.conn.Channel()
	defer ch.Close()

	msgs, _ := ch.Consume(
		queueName, "", false, false, false, false, nil,
	)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Recebi msg: %s", d.MessageId)

			err := w.handler.Handle(d.Body)

			if err != nil {
				log.Printf("Erro processando: %v", err)
				d.Nack(false, false)
			} else {
				d.Ack(false)
			}
		}
	}()

	log.Printf("Esperando mensagens na fila %s...", queueName)

	<-forever
}
