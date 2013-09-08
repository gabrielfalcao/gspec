package dsl
import (
	"fmt"
	"net/rpc"
)
var indent = 0
type TestCreator func(int) Expectation
type TestCallback func()

type TestResult struct {
	Description string
	Traceback error
	Level int
}
type SpecContext struct {
	ServerAddress string
}
type SpecDescription struct {
	Name string
}

func (self *TestResult) Passed() bool {
	return self.Traceback == nil
}

func (self SpecContext) Effectively(description string, spec TestCallback) {
	client, err := rpc.DialHTTP("tcp", self.ServerAddress)
	if err != nil {
		fmt.Println("\033[31mDIAL ERROR:", err, "\033[0m")
		return
	}

	meta := SpecDescription{description}
	client.Call("Runner.ReportBeginning", &meta, nil)

	defer func(name string, level int){
		var tb error
		if err := recover(); err != nil {
			tb = err.(error)
		} else {
			tb = nil
		}

		result := &TestResult{name, tb, level}

		client.Call("Runner.ReportResult", result, nil)
		if err != nil {
			fmt.Println("\033[31mDIAL ERROR:", err, "\033[0m")
		}
	}(description, 0)

	spec()
}

func (self SpecContext) It(description string, run_da_spec TestCallback) {
	self.Effectively(fmt.Sprintf("It %s", description), run_da_spec)
}
func (self SpecContext) Given(description string, run_da_spec TestCallback) {
	self.Effectively(fmt.Sprintf("Given %s", description), run_da_spec)
}
func (self SpecContext) When(description string, run_da_spec TestCallback) {
	self.Effectively(fmt.Sprintf("When %s", description), run_da_spec)
}
func (self SpecContext) Then(description string, run_da_spec TestCallback) {
	self.Effectively(fmt.Sprintf("Then %s", description), run_da_spec)
}
func (self SpecContext) And(description string, run_da_spec TestCallback) {
	self.Effectively(fmt.Sprintf("And %s", description), run_da_spec)
}
func (self SpecContext) Describe(description string, run_da_suite TestCallback) {
	self.Effectively(fmt.Sprintf("Describe %s", description), run_da_suite)
}
