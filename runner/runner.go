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

func (self *Spec) PersistAt(root *filesystem.Node) (dest *filesystem.Node, err error) {
	imports, specs := self.Parse()
	filename := self.GetFileName()
	specFile, err := root.NewFile(filename)
	if err != nil {
		return nil, err
	}
	specFile.Write([]byte("package main\n\n"))
	specFile.Write(imports)
	specFile.Write([]byte("\nfunc main() {\n"))
	specFile.Write(specs)
	specFile.Write([]byte("\n}\n"))
	specFile.Close()

	dest = root.Join(filename)
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
