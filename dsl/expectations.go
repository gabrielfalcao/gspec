package dsl
import (
	"fmt"
)
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
