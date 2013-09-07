package filesystem
import (
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
)

type Node struct {
	GivenPath string
}
type NodeList []Node


func (self Node) IsGoFile() (bool){
	return strings.HasSuffix(self.GivenPath, ".go")
}
func (self Node) Path() (abs string){
	abs, _ = filepath.Abs(self.GivenPath)
	return
}
func (self Node) IsDir() (bool){
	info, err := self.Stat()
	if err != nil {
		panic("FUDEU MUITO!!!!:\n\n" )
	}
	return info.IsDir()
}
func (self Node) Folder() (parent *Node){
	if self.IsDir() {
		return &self
	} else {
		other := self.Join("..")
		return other
	}
}
func (self Node) Parent() (parent *Node){
	folder := self.Folder()
	parent = folder.Join("..")
	return
}
func (self Node) Open() (file *os.File, err error) {
	return os.Open(self.Path())
}
func (self Node) Delete() error{
	return os.Remove(self.Path())
}

func (self Node) Stat() (info os.FileInfo, err error){
	file, err := self.Open()
	if err != nil {
		return
	}

	return file.Stat()
}
func (self Node) Read() (content []byte) {
	file, err := self.Open()
	if err != nil {
		return []byte(err.Error())
	}
	content, _ = ioutil.ReadAll(file)
	return
}
func (self Node) NewFile(name string) (file *os.File, err error) {
	newNode := self.Join(name)
	return os.Create(newNode.Path())
}
func (self Node) RelPath(other Node) (rel string){
	rel, _ = filepath.Rel(self.Path(), other.Path())
	return
}
func (self Node) Join(other string) *Node {
	node, err := GetNode(filepath.Join(self.Path(), other))
	if err != nil {
		panic(err)
	}
	return &node
}

func MergeList (self NodeList, other NodeList) (result NodeList){
	result = make(NodeList, len(self))
	copy(result, self)
	next := len(result)
	copy(result[next:], other)
	return
}

func (self Node) ListFiles() (result NodeList, err error){
	items, err := ioutil.ReadDir(self.Path())
	if err != nil {
		return nil, err
	}

	result = make(NodeList, 0)

	for _, item := range items {
		child := self.Join(item.Name())
		result = append(result, *child)
		if item.IsDir() {
			children, err := child.ListFiles()
			if err != nil {
				return result, err
			}
			result = MergeList(result, children)
		}
	}
	return
}
func GetNode(path string) (ret Node, err error) {
	return Node{path}, nil
}
