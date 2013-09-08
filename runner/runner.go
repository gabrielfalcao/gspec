package runner
import (
	"os"
	"fmt"
	"regexp"
	"os/exec"
	"go/token"
	"crypto/md5"
	"net"
	"net/rpc"
	"net/http"
	"github.com/gabrielfalcao/gspec/filesystem"
	"github.com/gabrielfalcao/gspec/dsl"
	. "github.com/gabrielfalcao/gspec/scanner"
)

type Runner struct {
	RootNode *filesystem.Node
	WorkingDirectory string
	ServerAddress string
}

type SpecSet struct {
	Nodes []filesystem.Node
	TokenFiles token.FileSet
	currentPos int
	Parent Runner
}

func NewSpecSet (runner Runner, nodes []filesystem.Node) (set SpecSet) {
	fset := token.NewFileSet()
	set = SpecSet{nodes, *fset, 0, runner}
	return
}

type Spec struct {
	Node *filesystem.Node
	Parent *SpecSet
	Runner *Runner
}

func (self *Spec) PersistAt(root *filesystem.Node) (dest *filesystem.Node, err error) {
	imports, specs := self.Parse()
	filename := self.GetFileName()
	specFile, err := root.NewFile(filename)
	if err != nil {
		return nil, err
	}
	specFile.Write([]byte("package main\n\n"))
	specFile.Write([]byte("import gspec_dsl \"github.com/gabrielfalcao/gspec/dsl\"\n"))
	specFile.Write(imports)
	specFile.Write([]byte("\nfunc main() {\n"))
	specFile.Write([]byte("        hub := gspec_dsl.SpecContext{\""+self.Runner.ServerAddress+"\"}\n"))
	specFile.Write(self.Translate(specs))
	specFile.Write([]byte("\n}\n"))
	specFile.Close()

	dest = root.Join(filename)
	return
}

func (self *Spec) Translate (original []byte) (translated []byte) {
	re := regexp.MustCompilePOSIX("(It|Describe|Given|When|Then)[(]")
	source := string(original)
	destination := re.ReplaceAllString(source, "hub.$1(")
	translated = []byte(destination)
	return
}

func (self *Spec) Parse () (imports []byte, describes[]byte) {
	imports, describes = ParseFile(self.Node, self.Parent.TokenFiles)
	return
}
func (self *Spec) Hash() (sum string) {
	hash := md5.New()
	hash.Write([]byte(self.Node.Path()))
	sum = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

func (self *Spec) GetFileName() string {
	return fmt.Sprintf("gspec_%s.go", self.Hash())
}

func (self *SpecSet) Length() int {
	return len(self.Nodes)
}

func (self *SpecSet) Next() (child *Spec) {
	if (self.Length() > self.currentPos) {
		node := &self.Nodes[self.currentPos]
		child = &Spec{node, self, &self.Parent}
		self.currentPos++
		return
	}
	return nil
}

func ReportError(message string, params ...interface{}){
	fmt.Printf("\033[0;31mERROR: \033[1;37m%s\033[0m\n", fmt.Sprintf(message, params...))
}

func (self *Runner) ReportResult(result *dsl.TestResult, reply *bool) error{
	if result.Passed() {
		fmt.Printf("\033[A\033[1;32m%s\033[0m\n", result.Description)
	} else {
		fmt.Printf("\033[A\033[1;31m%s\033[0m\n", result.Description)
		fmt.Println("\033[0;31m" + result.Traceback.Error() + "\033[0m")
	}
	return nil
}

func (self *Runner) ReportBeginning(spec *dsl.SpecDescription, reply *bool) error{
	fmt.Printf("\033[0;37m%s\033[0m\n", spec.Name)
	return nil
}

func (self *Runner) Run() {
	var err error
	if (err != nil) {
		ReportError("%s", err)
		return
	}

	rpc.Register(self)
	rpc.HandleHTTP()

	handler, err := net.Listen("tcp", self.ServerAddress)
	if err != nil {
		fmt.Println("\033[31m", err, "\033[0m")
		return
	}

	go http.Serve(handler, nil)

	curnode := filesystem.Node{self.WorkingDirectory}
	files, err := self.RootNode.ListFiles()

	if (err != nil) {
		ReportError("path \"%s\" does not exist", os.Args[1])
		return
	}

	specs := NewSpecSet(*self, files)

	for spec := specs.Next(); spec != nil; spec = specs.Next() {
		var destinationNode *filesystem.Node

		destinationNode, err := spec.PersistAt(self.RootNode)
		if err != nil {
			fmt.Printf("\n\033[37mignoring \033[33m%s\033[0m because %s\n", curnode.RelPath(*spec.Node), err)
			continue
		}

		cmd := exec.Command("go", "run", destinationNode.Path())

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		destinationNode.Delete()
		if err != nil {
			ReportError("exec: %s %s", err, destinationNode.Path())
			return
		}
	}
}
