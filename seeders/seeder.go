package seeders

type Seeder interface {
	Name() string
	Seed() error
	Truncate() error
}
