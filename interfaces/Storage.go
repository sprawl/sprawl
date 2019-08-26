package interfaces

type Storage interface {
	SetDbPath(dbPath string)
	Run() error
	Close()
	Get(key []byte) ([]byte, error)
	Put(key []byte, data []byte) error
	Delete(key []byte) error
}
