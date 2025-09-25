package gen

func NewDto(m interface{}, templateDir string) *Template {
	return NewTemplate(m, templateDir).
		SetFileSuffix("controller").
		SetTemplateName("dto\\dto.template").
		SetFileSuffix("dto")
}
