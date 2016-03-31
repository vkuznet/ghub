package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// Print give string in color
// based on ANSI escape codes
// http://www.wikiwand.com/en/ANSI_escape_code#/Colors
// specialized Go package: https://github.com/fatih/color
func useColor(user string) string {
	code := 31 // red color
	colored := fmt.Sprintf("\x1b[%d;1m%s\x1b[0m", code, user)
	return colored
}

// helper function to handle strings, we either return non-empty string with newline+tab
// or return empty string if input is empty
func tabString(msg string) string {
	if len(msg) != 0 {
		return fmt.Sprintf("%s\n\t", msg)
	}
	return ""
}

func PrintResults(request string, result Results) {
	var tmpl string
	switch request {
	case "search":
		tmpl = `{{.TotalCount}} issues:
{{range .Items}}#{{.Number}} {{.User.Login | useColor}} {{.Title | printf "%.64s" }}
{{end}}`
	case "issues":
		tmpl = `{{.TotalCount}} issues:
{{range .Items}}#{{.Number}} {{.User.Login | useColor }} {{.Title | printf "%.64s" }}
      {{.HTMLURL}}
{{end}}`
	case "issue":
		tmpl = `
{{range .Items}}
#{{.Number}} {{.User.Login | useColor }} {{.Title | printf "%.64s" }}

{{.Body}}
{{end}}`
	case "repos":
		tmpl = `
{{range .Items}}
{{.FullName | useColor}}
	{{.Description | tabString}}{{.Language | printf "Language: %s"}}, {{.Fork | printf "fork: %v"}}
{{end}}`
	default:
		fmt.Printf("%d issues:\n", result.TotalCount)
		for _, item := range result.Items {
			fmt.Printf("#%-5d %9.9s %.55s\n",
				item.Number, item.User.Login, item.Title)
		}
		return
	}
	var report = template.Must(template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Funcs(template.FuncMap{"useColor": useColor}).
		Funcs(template.FuncMap{"tabString": tabString}).
		Parse(tmpl))

	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}
