package cache

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
	GetAll() (map[string][]byte, error)
	Delete(key string) error
	FindValue(value string) (string, error)
}
