package consumer

type ConsumerID string

type Consumer struct {
	ID      ConsumerID
	Name    string
	APIKey  string
	Enabled bool
}

func (c *Consumer) Disable() {
	c.Enabled = false
}
