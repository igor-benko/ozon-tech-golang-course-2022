package entity

type Person struct {
	ID        uint64 `db:"id"`
	LastName  string `db:"last_name"`
	FirstName string `db:"first_name"`
}

type PersonPage struct {
	Persons []Person
	Total   uint64
}

type PersonFilter struct {
	Offset uint64
	Limit  uint64
	Order  string
}
