package entity

type Vehicle struct {
	ID        uint64 `db:"id"`
	Brand     string `db:"brand"`
	Model     string `db:"model"`
	RegNumber string `db:"reg_number"`
	PersonID  uint64 `db:"person_id"`
}

type VehiclePage struct {
	Vehicles []Vehicle
	Total    uint64
}

type VehicleFilter struct {
	Offset uint64
	Limit  uint64
	Order  string
}
