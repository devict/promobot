package channels

type Channel interface {
	Send(string) error
	Type() string
	Name() string
}
