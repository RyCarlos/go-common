package gen

func NewRoute(m interface{}, templateDir string) *Template {
	return NewTemplate(m, templateDir).
		SetTemplateName("route.template")
}
