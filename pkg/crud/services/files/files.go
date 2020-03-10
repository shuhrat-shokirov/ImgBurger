package files

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type FilesSvc struct {
	mediaPath string
}

func NewFilesSvc(mediaPath string) *FilesSvc {
	if mediaPath == "" {
		panic(errors.New("media path can't be nil")) // <- be accurate
	}

	return &FilesSvc{mediaPath: mediaPath}
}

func (receiver *FilesSvc) Save(src io.Reader, contentType string) (name string, err error) {
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", err
	}

	if len(extensions) == 0 {
		return "", errors.New("invalid extension")
	}

	uuidV4 := uuid.New().String()
	name = fmt.Sprintf("%s%s", uuidV4, extensions[0])
	path := filepath.Join(receiver.mediaPath, name)

	dst, _ := os.Create(path)
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return name, nil
}
