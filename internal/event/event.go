package event

type Event interface {
	isEvent()
}

type RentalAdded struct {
	PropertyID string
}
type RentalListed struct {
	PropertyID string
}
type RentalLeased struct {
	PropertyID string
	LeaseID    string
}
