package service

type Service interface {
	ReadStream() ([]string, error)
	GetName() string
}
