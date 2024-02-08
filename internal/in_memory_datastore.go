package internal

type InMemoryDatastore struct {
	OKRs []*OKR
}

func (d *InMemoryDatastore) Add(okr *OKR) error {
	d.OKRs = append(d.OKRs, okr)
	return nil
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
