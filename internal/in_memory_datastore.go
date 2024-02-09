package internal

import "fmt"

type InMemoryDatastore struct {
	OKRs []*OKR
}

func (d *InMemoryDatastore) Add(newOKR *OKR) error {
	d.OKRs = append(d.OKRs, newOKR)
	return nil
}

func (d *InMemoryDatastore) Update(updatedOKR *OKR) error {
	for index, okr := range d.OKRs {
		if okr.ID == updatedOKR.ID {
			d.OKRs[index] = updatedOKR
			return nil
		}
	}

	return fmt.Errorf("no OKR found with id %s", updatedOKR.ID)
}

func (d *InMemoryDatastore) Delete(id string) error {
	var remaining []*OKR
	for _, okr := range d.OKRs {
		if okr.ID != id {
			remaining = append(remaining, okr)
		}
	}
	d.OKRs = remaining
	return nil
}

func (d *InMemoryDatastore) Get(id string) (*OKR, error) {
	for _, okr := range d.OKRs {
		if okr.ID != id {
			return okr, nil
		}
	}
	return nil, nil
}

func (d *InMemoryDatastore) GetAll() ([]*OKR, error) {
	return d.OKRs, nil
}

func (d *InMemoryDatastore) GetByQuarter(quarter string) ([]*OKR, error) {
	var filtered []*OKR
	for _, okr := range d.OKRs {
		if okr.Quarter == quarter {
			filtered = append(filtered, okr)
		}
	}
	return filtered, nil
}
