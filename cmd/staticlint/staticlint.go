package main

import (
	"fmt"
	"go/ast"
	"strings"

	_ "github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

var (
	goCroticAnalyzers = []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,

		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}
)

func main() {
	mychecks := []*analysis.Analyzer{
		OsExitAnalyzer,
	}

	// Добавляем все анализаторы класса SA пакета staticcheck.io
	// категория проверок SA под кодовым названием staticcheck включает все проверки, связанные с правильностью кода.
	for _, v := range staticcheck.Analyzers {

		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// Добавляем проверку на неправильно выбранный идентификатор
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name == "ST1003" {
			mychecks = append(mychecks, v.Analyzer)
			break
		}
	}
	// Добавляем все анализаторы статического анализатора Go-Critic
	mychecks = append(mychecks, goCroticAnalyzers...)
	multichecker.Main(
		mychecks...,
	)
}

// Проверка на отсутствие конструкции os.Exit() в функции main пакета main
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check calling os.Exit in main Function",
	Run:  runOS,
}

func runOS(pass *analysis.Pass) (interface{}, error) {

	var mainfunc string
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if p, ok := n.(*ast.FuncDecl); ok {
				mainfunc = p.Name.Name
			}

			if c, ok := n.(*ast.CallExpr); ok {

				res := fmt.Sprintf("%s", c.Fun)

				if res == "&{os Exit}" {
					if mainfunc == "main" {

						pass.Reportf(c.Pos(), "using os Exit!")
					}

				}
			}

			return true
		})
	}
	return nil, nil
}