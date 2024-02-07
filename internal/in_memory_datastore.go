package internal

type InMemoryDatastore struct {
	OKRs []*OKR
}

func (d *InMemoryDatastore) Add(okr *OKR) {
	d.OKRs = append(d.OKRs, okr)
}

func (d *InMemoryDatastore) GetAll() []*OKR {
	return d.OKRs
}

func (d *InMemoryDatastore) GetByQuarter(quarter string) []*OKR {
	var filtered []*OKR
	for _, okr := range d.OKRs {
		if okr.Quarter == quarter {
			filtered = append(filtered, okr)
		}
	}
	return filtered
}
