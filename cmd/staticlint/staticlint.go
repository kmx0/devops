package main

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// gocrA := goanalysis.DummyRun()

	mychecks := []*analysis.Analyzer{
		// gocrA,
		OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {

		if strings.HasPrefix(v.Name, "SA") {
			mychecks = append(mychecks, v)
		}
	}

	// if v, ok := staticcheck.Analyzers; ok {
	// 	mychecks = append(mychecks, v)
	// }
	for _, v := range staticcheck.Analyzers {
		if v.Name == "ST1003" {
			mychecks = append(mychecks, v)
			break
		}
	}
	multichecker.Main(
		mychecks...,
	)
}

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check calling os.Exit in main Function",
	Run:  runOS,
}

func runOS(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if c, ok := n.(*ast.CallExpr); ok {

				res := fmt.Sprintf("%s", c.Fun)

				if res == "&{os Exit}" {
					pass.Reportf(c.Pos(), "using os Exit!")

				}
			}

			return true
		})
	}
	return nil, nil
}
