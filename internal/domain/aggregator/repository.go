package aggregator

// RepositoryReader is a interface for read application
type RepositoryReader interface {
	Find(string) (*Entity, error)
}

// RepositoryWriter is a interface for write application
type RepositoryWriter interface {
	Insert(Entity) (*Entity, error)
	Delete(string) error
	Update(Entity) (*Entity, error)
}

// Repository is a interface for read 'n write application
type Repository interface {
	RepositoryReader
	RepositoryWriter
}
