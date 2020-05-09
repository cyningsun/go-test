package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"strconv"
)
const tpl = `
package {{.pkg}}

var {{.typ}}Map = map[uint32] {{.typ}} {
	{{range $name,$val :=.vars}}{{$name}}: {{$val}},
	{{end}}
}
`

var (
	pkgInfo *build.Package
)
var (
	t = flag.String("type", "", "require")
)

func parse(pkg *build.Package, typ string) map[uint64]string {
	//解析当前目录下包信息
	codeMap := make(map[uint64]string, 0)
	fset := token.NewFileSet()
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			// 只需要const
			if !ok || decl.Tok != token.VAR {
				return true
			}
			for _, spec := range decl.Specs {
				vspec, ok := spec.(*ast.ValueSpec)
				if !ok  {
					continue
				}
				key,val,ok := parseNode(typ, vspec) 
				if !ok {
					continue
				}
				codeMap[key] = val
			}
			return true
		})
	}
	return codeMap
}

func parseNode(typ string, vspec *ast.ValueSpec) (uint64,string,bool){
	if len(vspec.Values) != 1 {
		return 0,"",false
	}
	if len(vspec.Names) < 1 || vspec.Names[0].Name == "" {
		return 0,"",false
	}
	name := vspec.Names[0].Name

	vspecVal, ok := vspec.Values[0].(*ast.CompositeLit)
	if !ok || vspecVal.Type == nil {
		return 0,"",false
	}
	ident, ok := vspecVal.Type.(*ast.Ident)
	if !ok || ident.Name != typ {
		return 0,"",false
	}
	if !ok || len(vspecVal.Elts) < 1 {
		return 0,"",false
	}
	code, ok := vspecVal.Elts[0].(*ast.KeyValueExpr)
	if !ok {
		return 0,"",false
	}
	codeVal, ok := code.Value.(*ast.BasicLit)
	if !ok || codeVal.Value == "" {
		return 0,"",false
	}
	val, err := strconv.ParseUint(codeVal.Value, 10, 64)
	if err != nil {
		return 0,"",false
	}
	return val, name, true
}

func render(pkg, typ string,  vars map[uint64]string) []byte {
	data := map[string]interface{}{
		"pkg":  pkg,
		"typ":  typ,
		"vars": vars,
	}
	//利用模板库，生成代码文件
	t, err := template.New("").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		log.Fatal(err)
	}
	//格式化
	src, err := format.Source(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return src
}

func main() {
	flag.Parse()
	if len(*t) == 0 {
		log.Fatal("-type require")
	}
	pkgInfo, err := build.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}
	pkgName := os.Getenv("GOPACKAGE")
	if pkgName == "" {
		pkgName = pkgInfo.Name
	}

	codeMap := parse(pkgInfo, *t)
	buffer := render(pkgName, *t, codeMap)
	//保存到文件
	filename := ""
	if filename == "" {
		baseName := fmt.Sprintf("%s_map.go", *t)
		filename = filepath.Join(".", strings.ToLower(baseName))
	}
	err = ioutil.WriteFile(filename, buffer, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}