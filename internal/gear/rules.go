package gear

type RulesData struct {
	SpecRelPath string // ".gear/<name>.spec"
}

func Render(d RulesData) ([]byte, error)
