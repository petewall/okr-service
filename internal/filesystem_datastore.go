package internal

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type FilesystemDatastore struct {
	okrs   []*OKR
	Path   string
	Format string
}

func (d *FilesystemDatastore) Initialize() error {
	log.Debugf("Loading OKRs from file: %s", d.Path)
	data, err := os.ReadFile(d.Path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var okrs []*OKR
	if d.Format == "json" {
		err = json.Unmarshal(data, &okrs)
	} else if d.Format == "yaml" {
		err = yaml.Unmarshal(data, &okrs)
	} else {
		return fmt.Errorf("unsupported datastore format: %s", d.Format)
	}
	if err != nil {
		return fmt.Errorf("failed to parse data store file: %w", err)
	}
	log.Debugf("Found %d OKRs from the file", len(okrs))

	return nil
}

func (d *FilesystemDatastore) save() error {
	var data []byte
	var err error

	if d.Format == "json" {
		data, err = json.Marshal(d.okrs)
	} else if d.Format == "yaml" {
		data, err = yaml.Marshal(d.okrs)
	} else {
		return fmt.Errorf("unsupported datastore format: %s", d.Format)
	}

	if err != nil {
		return fmt.Errorf("failed to encode data store file: %w", err)
	}

	err = os.WriteFile(d.Path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data store file: %w", err)
	}
	return nil
}

func (d *FilesystemDatastore) Add(okr *OKR) error {
	d.okrs = append(d.okrs, okr)
	return d.save()
}

func (d *FilesystemDatastore) GetAll() ([]*OKR, error) {
	return d.okrs, nil
}

func (d *FilesystemDatastore) GetByQuarter(quarter string) ([]*OKR, error) {
	var filtered []*OKR
	for _, okr := range d.okrs {
		if okr.Quarter == quarter {
			filtered = append(filtered, okr)
		}
	}
	return filtered, nil
}
