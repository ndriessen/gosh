package gitops

type App struct {
	Name       string
	Properties map[string]string
	group      *AppGroup
}
