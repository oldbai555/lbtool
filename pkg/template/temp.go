package template

type Function struct {
	Package     string `json:"package"`
	ModelName   string `json:"model_name"`
	Variable    string `json:"variable"`
	Description string `json:"description"`
	Template    string `json:"template"`
}

var teplDemo1 = `
package {{.Package}}

// TestGenTemplate {{.Description}}
func TestGenTemplate() string {
	return fmt.Sprintf("%s", "{{.ModelName}}")
}`
