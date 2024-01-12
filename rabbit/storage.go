package rabbit

import (
	"errors"
	"time"
)

type Storage[T any] struct {
	messages map[string]*T
}

func NewRabbitStorage[T any]() Storage[T] {
	return Storage[T]{messages: make(map[string]*T)}
}

func (s *Storage[T]) GetMessageById(id string) (*T, error) {
	timeout := time.After(time.Second * 30)
	for {
		select {
		case <-timeout:
			return nil, errors.New("InternalError")
		default:
			msg, ok := s.messages[id]
			if !ok {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			delete(s.messages, id)
			return msg, nil
		}
	}
}

func (s *Storage[T]) SaveMessageWithId(id string, msg *T) {
	s.messages[id] = msg
}
