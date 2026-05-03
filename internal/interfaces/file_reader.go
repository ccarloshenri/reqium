package interfaces

type FileReader interface {
	Read(path string) ([]byte, error)
}
