package rest

type Check interface {
	Check() (bool, string)
}
