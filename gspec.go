package gspec
import (
	"os"
	"os/exec"
	"fmt"
	"strings"
	"crypto/md5"
	"go/token"
	"falcao.it/gspec/filesystem"
	. "falcao.it/gspec/scanner"
)
const INDENTATION_LEVEL = 2

var indent int

type Expectation struct {
	source interface{}
	positive bool
	To *Expectation
	Not *Expectation
}

func (self *Expectation) Equal(other interface{}) {
	// expected := fmt.Sprintf("%v", other)
	// got := fmt.Sprintf("%v", self.source)
	if self.source != other {
		panic(fmt.Sprintf("expected \"%v\" to equal \"%v\"", self.source, other))
	}
}

type TestCreator func(int) Expectation
type TestCallback func()

type TestResult struct {
	error error
	failure error
}
func (self *TestResult) Failure() error {
	return self.failure
}
func (self *TestResult) Error() error {
	return self.error
}
func (self *TestResult) Passed() bool {
	return self.error == nil && self.failure == nil
}

func (self *TestResult) Failed() bool {
	return self.failure != nil || self.error != nil
}

func Expect(source interface{}) Expectation{
	// starting with a positive expectation
	x := Expectation{source, true, nil, nil}
	// and creating a link to self
	x.To = &x
	// and creating a respective negative expectation
	x.Not = &Expectation{source, false, nil, nil}
	// and its self link as well
	x.Not.To = x.Not
	return x
}
func Effectively(spec TestCallback) {
	indent += INDENTATION_LEVEL
	spec()
	indent -= INDENTATION_LEVEL
}

func GetIndentation() string{
	return strings.Repeat(" ", indent)
}
func ShowCallback(name, description string) {
	indentation := GetIndentation()

	fmt.Printf("%s\033[32m%s %s\033[0m\n", indentation, name, description)
}
func It(description string, run_da_spec TestCallback) {
	ShowCallback("It", description)
	Effectively(run_da_spec)
}
func Given(description string, run_da_spec TestCallback) {
	ShowCallback("Given", description)
	Effectively(run_da_spec)
}
func When(description string, run_da_spec TestCallback) {
	ShowCallback("When", description)
	Effectively(run_da_spec)
}
func Then(description string, run_da_spec TestCallback) {
	ShowCallback("Then", description)
	Effectively(run_da_spec)
}
func And(description string, run_da_spec TestCallback) {
	ShowCallback("And", description)
	Effectively(run_da_spec)
}

func Describe(description string, run_da_suite TestCallback) {
	ShowCallback("Describe", description)
	indent = 0
	Effectively(run_da_suite)
}

func ReportError(message string, params ...interface{}){
	fmt.Printf("\033[0;31mERROR: \033[1;37m%s\033[0m\n", fmt.Sprintf(message, params...))
}

func Run() {
	var here filesystem.Node
	var err error
	if len(os.Args) == 1 {
		here, err = filesystem.GetNode(".")
	} else {
		here, err = filesystem.GetNode(os.Args[1])
	}
	if (err != nil) {
		ReportError("%s", err)
		return
	}

	wd, _ := os.Getwd()
	curnode := filesystem.Node{wd}
	files, err := here.ListFiles()
	if (err != nil) {
		ReportError("path \"%s\" does not exist", os.Args[1])
		return
	}
	fset := token.NewFileSet()

	for _, file := range files {
		imports, specs := ParseFile(file, *fset)
		if len(specs) == 0 {
			continue
		}
		hash := md5.New()
		hash.Write([]byte(file.Path()))

		specFileName := fmt.Sprintf("gspec_%x.go", hash.Sum(nil))
		specFile, err := here.NewFile(specFileName)
		if err != nil {
			ReportError("%s", err)
			return

		}
		specFile.Write([]byte("package main\n\n"))
		specFile.Write(imports)
		specFile.Write([]byte("\nfunc main() {\n"))
		specFile.Write(specs)
		specFile.Write([]byte("\n}\n"))
		specFile.Close()
		specnode := here.Join(specFileName)

		fmt.Printf("\033[35mHawkEye\033[1;37m is running \033[0;33m%s\033[0m...\n", curnode.RelPath(file))

		cmd := exec.Command("go", "run", specnode.Path())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		specnode.Delete()
		if err != nil {
			ReportError("exec: %s %s", err, specnode.Path())
			return
		}
	}
}
