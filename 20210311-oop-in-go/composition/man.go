package composition

type Man interface {
	Birthday() string
}

func NewMan(name, id string) (Man, error) {
	i, err := NewID(id)
	return &man{
		id:   i,
		name: name,
	}, err
}

func NewManEmbedding(name, id string) (Man, error) {
	i, err := NewID(id)
	return &manEmbedding{
		ID:   i,
		name: name,
	}, err
}

type man struct {
	id   ID
	name string
}

func (m *man) Birthday() string {
	return m.id.Birthday()
}

type manEmbedding struct {
	ID
	name string
}
