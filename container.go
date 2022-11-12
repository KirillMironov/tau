package tau

type Container struct {
	Id      string `toml:"-"`
	Image   string `validate:"required"`
	Command []string
}

func (c *Container) Start(runtime ContainerRuntime) error {
	containerId, err := runtime.Start(*c)
	if err != nil {
		return err
	}

	c.Id = containerId

	return nil
}

func (c *Container) Remove(runtime ContainerRuntime) error {
	err := runtime.Remove(c.Id)
	if err != nil {
		return err
	}

	c.Id = ""

	return nil
}
