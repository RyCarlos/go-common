package gen

func NewController(m interface{}, templateDir string) *Template {
	return NewTemplate(m, templateDir).
		SetFileSuffix("controller").
		SetTemplateName("controller.template")
}
