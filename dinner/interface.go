package dinner

type Client interface {
	UpdateMenu() error
	IsReady() bool
	Order() error
}
