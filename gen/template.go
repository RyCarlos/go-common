package gen

import (
	"bytes"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type Config struct {
	Model       *Model
	FileName    string // 文件名
	packageName string // 包名
	FileSuffix  string // 后缀
	CreatedAt   string // 文件创建时间
	SaveDirPath string // 文件保存目录
	Overwrite   bool   // 是否可覆盖
}

type Template struct {
	*Config
	TemplateName string
	TemplateDir  string
	content      bytes.Buffer // 模板内容
	IsRender     bool
}

func (t *Template) SetModel(m interface{}) *Template {
	t.Model = NewModel(m)
	return t
}

func (t *Template) SetPackageName(name string) *Template {
	t.packageName = name
	return t
}

func (t *Template) SetFileSuffix(suffix string) *Template {
	t.FileSuffix = suffix
	t.SetFileName(t.FileName)
	return t
}

func (t *Template) GetPackageName() string {
	return t.packageName
}

func (t *Template) SetSaveDirPath(path string) *Template {
	t.SaveDirPath = path
	return t
}

func (t *Template) SetTemplateName(name string) *Template {
	t.TemplateName = name
	return t
}

func (t *Template) SetOverwrite(overwrite bool) *Template {
	t.Overwrite = overwrite
	return t
}

func (t *Template) SetFileName(name string) *Template {
	t.FileName = name
	return t
}

func (t *Template) GetFileName() string {
	if t.FileSuffix != "" {
		return t.FileName + "_" + t.FileSuffix + ".go"
	}
	return t.FileName + ".go"
}

func (t *Template) render(data interface{}) *Template {
	if !t.IsRender {
		files, err := template.ParseFiles(filepath.Join(t.TemplateDir, t.TemplateName))
		if err != nil {
			panic(fmt.Errorf("解析模板失败: %s", err))
		}
		err = files.Execute(&t.content, data)
		if err != nil {
			panic(fmt.Errorf("渲染模板失败: %s", err))
		}
	}
	t.IsRender = true
	return t
}

func (t *Template) Build() {
	t.SetPackageName(filepath.Base(t.SaveDirPath)).render(t).Save(t.content)
}

func (t *Template) Save(content bytes.Buffer) bool {
	filePath := filepath.Join(t.SaveDirPath, t.GetFileName())
	// 是否开启覆盖模式
	if !t.Overwrite && fileutil.IsExist(filePath) {
		fmt.Printf("--------未开启文件覆盖模式，文件[%s]已经存在,写入失败！\r\n", filePath)
		return false
	}
	err := os.MkdirAll(t.SaveDirPath, 0755)
	if err != nil {
		panic(fmt.Errorf("创建目录[%s]失败: %s", t.SaveDirPath, err))
	}

	file, err := os.Create(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("关闭文件失败: %s", err))
		}
	}(file)
	if err != nil {
		panic(fmt.Errorf("保存文件失败: %s", err))
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		panic(fmt.Errorf("模板内容写入文件失败: %s", err))
	}

	fmt.Printf("--------Api文件[%s]保存成功,保存目录[%s]\r\n", t.GetFileName(), t.SaveDirPath)
	return true
}

func NewTemplate(m interface{}, templateDir string) *Template {
	model := NewModel(m)
	config := &Config{
		Model:     model,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Overwrite: false,
	}
	t := &Template{
		Config:      config,
		TemplateDir: templateDir,
		content:     bytes.Buffer{},
		IsRender:    false,
	}
	t.SetFileName(model.GetNameSnakeCase())
	return t
}
