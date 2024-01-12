package staticerr

import "errors"

var (
	ErrorRabbitConnectionFail = errors.New("RabbitConnectionFail")
)
