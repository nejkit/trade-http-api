package providers

type IStorage[T any] interface {
	GetMessageById(id string) (*T, error)
}

type BpsProvider struct {
}
