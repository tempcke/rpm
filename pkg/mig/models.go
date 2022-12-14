package mig

type Flow []Step

type Step struct {
	ID   string
	Up   string
	Down string
}
