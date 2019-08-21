package dinner

type Client interface {
	HealthCheck() error
	Order() error
}
