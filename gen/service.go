package gen

func NewService(m interface{}, templateDir string) *Template {
	return NewTemplate(m, templateDir).
		SetFileSuffix("service").
		SetTemplateName("service.template")
}
