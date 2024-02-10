package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type FilesystemDatastore struct {
	cache  *InMemoryDatastore
	Path   string
	Format string
}

func (d *FilesystemDatastore) Initialize() error {
	var okrs []*OKR

	log.Debugf("Loading OKRs from file: %s", d.Path)
	data, err := os.ReadFile(d.Path)
	if err == nil {
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
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to read file: %w", err)
	}

	d.cache = &InMemoryDatastore{
		OKRs: okrs,
	}
	return nil
}

func (d *FilesystemDatastore) save() error {
	var data []byte
	var err error

	if d.Format == "json" {
		data, err = json.Marshal(d.cache.OKRs)
	} else if d.Format == "yaml" {
		data, err = yaml.Marshal(d.cache.OKRs)
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
	err := d.cache.Add(okr)
	if err != nil {
		return err
	}
	return d.save()
}

func (d *FilesystemDatastore) Update(updatedOKR *OKR) error {
	err := d.cache.Update(updatedOKR)
	if err != nil {
		return err
	}
	return d.save()
}

func (d *FilesystemDatastore) Delete(id string) error {
	err := d.cache.Delete(id)
	if err != nil {
		return err
	}
	return d.save()
}

func (d *FilesystemDatastore) Get(id string) (*OKR, error) {
	return d.cache.Get(id)
}

func (d *FilesystemDatastore) GetAll() ([]*OKR, error) {
	return d.cache.GetAll()
}

func (d *FilesystemDatastore) GetByQuarter(quarter string) ([]*OKR, error) {
	return d.cache.GetByQuarter(quarter)
}
