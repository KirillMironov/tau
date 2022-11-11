package tau

import "github.com/BurntSushi/toml"

type Pod struct {
	Kind       string      `validate:"eq=pod"`
	Name       string      `validate:"required"`
	Containers []Container `validate:"required,dive,required"`
}

func NewPod(data []byte) (pod Pod, err error) {
	err = toml.Unmarshal(data, &pod)
	if err != nil {
		return Pod{}, err
	}

	err = validate.Struct(&pod)
	if err != nil {
		return Pod{}, err
	}

	return pod, nil
}
