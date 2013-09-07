package runner
import (
	"os"
	"os/exec"
	"fmt"
	"crypto/md5"
	"go/token"
	"github.com/gabrielfalcao/gspec/filesystem"
	. "github.com/gabrielfalcao/gspec/scanner"
)

type SpecSet struct {
	Nodes []filesystem.Node
	TokenFiles token.FileSet
	currentPos int
}

func NewSpecSet (nodes []filesystem.Node) (set SpecSet) {
	fset := token.NewFileSet()
	set = SpecSet{nodes, *fset, 0}
	return
}

type Spec struct {
	Node *filesystem.Node
	Parent *SpecSet
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

func (self *SpecSet) Length () int {
	return len(self.Nodes)
}

func (self *SpecSet) Next () (child *Spec) {
	if (self.Length() > self.currentPos) {
		node := &self.Nodes[self.currentPos]
		child = &Spec{node, self}
		self.currentPos++
		return
	}
	return nil
}



type Runner struct {
	RootNode *filesystem.Node
	WorkingDirectory string
}

func ReportError(message string, params ...interface{}){
	fmt.Printf("\033[0;31mERROR: \033[1;37m%s\033[0m\n", fmt.Sprintf(message, params...))
}


func (self *Runner) Run() {
	var err error
	if (err != nil) {
		ReportError("%s", err)
		return
	}

	curnode := filesystem.Node{self.WorkingDirectory}
	files, err := self.RootNode.ListFiles()

	if (err != nil) {
		ReportError("path \"%s\" does not exist", os.Args[1])
		return
	}

	specs := NewSpecSet(files)

	for spec := specs.Next(); spec != nil; spec = specs.Next() {
		imports, specs := spec.Parse()
		if len(specs) == 0 {
			continue
		}

		specFileName := spec.GetFileName()
		specFile, err := self.RootNode.NewFile(specFileName)
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
		specnode := self.RootNode.Join(specFileName)

		fmt.Printf("\n\033[37mGSpec \033[0mis running %s...\n", curnode.RelPath(*spec.Node))

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
