package gspec
import (
	"fmt"
	"strings"
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
