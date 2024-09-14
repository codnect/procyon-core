package config

import (
	"codnect.io/procyon-core/runtime/property"
	"io/fs"
	"path/filepath"
	"strings"
)

type Resource interface {
	Name() string
	Location() string
	Profile() string
	Loader() property.SourceLoader
}

type FileResource struct {
	path   string
	file   fs.File
	loader property.SourceLoader
}

func newFileResource(path string, file fs.File, loader property.SourceLoader) *FileResource {
	if strings.TrimSpace(path) == "" {
		panic("cannot create file resource with empty or blank path")
	}

	if file == nil {
		panic("nil file")
	}

	if loader == nil {
		panic("nil loader")
	}

	return &FileResource{
		path:   path,
		file:   file,
		loader: loader,
	}
}

func (r *FileResource) File() fs.File {
	return r.file
}

func (r *FileResource) Location() string {
	return r.path
}

func (r *FileResource) Name() string {
	return filepath.Base(r.Location())
}

func (r *FileResource) Profile() string {
	fileName := filepath.Base(r.Location())
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	nameParts := strings.Split(fileName, "-")
	if len(nameParts) == 1 {
		return ""
	}

	return nameParts[len(nameParts)-1]
}

func (r *FileResource) Loader() property.SourceLoader {
	return r.loader
}
