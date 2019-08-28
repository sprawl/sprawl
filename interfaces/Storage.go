package interfaces

// Storage defines a database interface that works with Sprawl
type Storage interface {
	SetDbPath(dbPath string)
	Run() error
	Close()
	Get(key []byte) ([]byte, error)
	Put(key []byte, data []byte) error
	Delete(key []byte) error
	GetAll() (map[string]string, error)
	GetAllWithPrefix(prefix string) (map[string]string, error)
	DeleteAll() error
	DeleteAllWithPrefix(prefix string) error
}
