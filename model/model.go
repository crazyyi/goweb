package model

type Model struct {
	db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) People() ([]*Person, error) {
	return m.SelectPeople()
}

func (m *Model) CreateRecord(p *Person) (int, error) {
	return m.CreateNewRecord(p)
}

func (m *Model) DeleteRow(id int64) (int64, error) {
	return m.DeleteRowAt(id)
}