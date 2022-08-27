package template

type Function struct {
	Package     string `json:"package"`
	Description string `json:"description"`
	ModelName   string `json:"model_name"`
	Template    string `json:"template"`
}

var teplDemo1 = `
package {{.Package}}

// TestGenTemplate {{.Description}}
func TestGenTemplate() string {
	return fmt.Sprintf("%s", "{{.ModelName}}")
}`
