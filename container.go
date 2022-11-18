package tau

type Container struct {
	Image   string `validate:"required"`
	Command string
	id      string
}

func (c *Container) Id() string {
	return c.id
}

func (c *Container) SetId(id string) {
	c.id = id
}
