package tau

type Container struct {
	Id      string `toml:"-"`
	Image   string `validate:"required"`
	Command string
}
