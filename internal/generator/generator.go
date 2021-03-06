package generator

type Options struct {
	PkgName string
}

func Generate(schemaString string, o Options) string {
	return "package " + o.PkgName + "\n"
}
