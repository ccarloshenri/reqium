package filesystem

import "os"

type OSFileReader struct{}

func NewOSFileReader() *OSFileReader {
	return &OSFileReader{}
}

func (r *OSFileReader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
