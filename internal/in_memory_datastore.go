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
	for _, okr := range d.OKRs {
		if okr.ID == updatedOKR.ID {
			okr.Quarter = updatedOKR.Quarter
			okr.Category = updatedOKR.Category
			okr.ValueType = updatedOKR.ValueType
			okr.Description = updatedOKR.Description
			okr.Goal = updatedOKR.Goal
			okr.Progress = updatedOKR.Progress
			okr.UpdateMetrics()
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
	filtered := []*OKR{}
	for _, okr := range d.OKRs {
		if okr.Quarter == quarter {
			filtered = append(filtered, okr)
		}
	}
	return filtered, nil
}
