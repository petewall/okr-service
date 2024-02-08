package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

type Datastore interface {
	Add(*OKR) error
	GetAll() ([]*OKR, error)
	GetByQuarter(quarter string) ([]*OKR, error)
}

const (
	DatastoreTypeFilesystem = "fs"
	DatastoreTypeMemory     = "memory"
)

func InitializeDatastore() (Datastore, error) {
	dsType := viper.GetString("datastore.type")
	if dsType == DatastoreTypeFilesystem {
		path := viper.GetString("datastore.path")
		if path == "" {
			return nil, fmt.Errorf("missing datastore path")
		}

		ds := &FilesystemDatastore{
			Path:   path,
			Format: viper.GetString("datastore.format"),
		}
		return ds, ds.Initialize()
	} else if dsType == DatastoreTypeMemory {
		return &InMemoryDatastore{}, nil
	}
	return nil, fmt.Errorf("unsupproted datastore type: %s", dsType)
}
