package internal

type Datastore interface {
	Add(*OKR)
	GetAll() []*OKR
	GetByQuarter(quarter string) []*OKR
}
