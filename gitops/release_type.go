package gitops

type ReleaseType int

const (
	StageRelease ReleaseType = iota
	ProductRelease
	HotFixRelease
)

func (r ReleaseType) String() string {
	return [...]string{"stage", "product", "hotfix"}[r]
}

func (r ReleaseType) EnumIndex() int {
	return int(r)
}
