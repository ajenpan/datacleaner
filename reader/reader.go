package reader

type Reader interface {
	Read() ([]byte, error)
	Close() error
}
