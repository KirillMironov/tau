package resources

const (
	KindContainer Kind = "container"
	KindPod       Kind = "pod"
)

type Kind string

func (k Kind) String() string {
	return string(k)
}
