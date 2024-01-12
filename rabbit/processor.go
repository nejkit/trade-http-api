package rabbit

import "github.com/rabbitmq/amqp091-go"

type IdFunc[T any] func(*T) string
type ParserFunc[T any] func(amqp091.Delivery) (*T, error)

type Processor[T any] struct {
	idFunc     IdFunc[T]
	parserFunc ParserFunc[T]
	storage    *Storage[T]
}

func NewProcessor[T any](idFunc IdFunc[T], parserFunc ParserFunc[T], storage *Storage[T]) Processor[T] {
	return Processor[T]{idFunc: idFunc, storage: storage, parserFunc: parserFunc}
}

func (p *Processor[T]) processMessage(msg amqp091.Delivery) {
	body, err := p.parserFunc(msg)

	if err != nil {
		msg.Nack(false, false)
		return
	}

	msgId := p.idFunc(body)

	p.storage.SaveMessageWithId(msgId, body)
}
