package gitops

import "errors"

type ReleaseType int

const (
	StageRelease ReleaseType = iota + 1
	ProductRelease
	HotFixRelease
)

var UnsupportedReleaseTypeErr = errors.New("unsupported release type")

func (r ReleaseType) String() string {
	return [...]string{"stage", "product", "hotfix"}[r-1]
}

func (r ReleaseType) EnumIndex() int {
	return int(r)
}

func NewReleaseType(value string) (ReleaseType, error) {
	switch value {
	case "stage":
		return StageRelease, nil
	case "product":
		return ProductRelease, nil
	case "hotfix":
		return HotFixRelease, nil
	}
	return 0, UnsupportedReleaseTypeErr
}
