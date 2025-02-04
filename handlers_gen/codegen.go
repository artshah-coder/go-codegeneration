package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type decodeTpl[T int | string] struct {
	FieldName string
	ParamName string
	Required  bool
	Default   string
	Enum      []T
	Min       int
	Max       int
}

type Method struct {
	MethodName string
	MethodType string `json:"method"`
	InType     string
	OutType    ast.Expr
	URL        string `json:"url"`
	Auth       bool   `json:"auth"`
}

const (
	postTmpl        = "templates/postTmpl.gotext"
	getTmpl         = "templates/getTmpl.gotext"
	multiplexerTmpl = "templates/multiplexerTmpl.gotext"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(out, `// GENERATED CODE`)
	fmt.Fprintln(out, `// DO NOT EDIT`)
	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import (`)
	fmt.Fprintln(out, "\t"+`"encoding/json"`)
	fmt.Fprintln(out, "\t"+`"fmt"`)
	fmt.Fprintln(out, "\t"+`"io"`)
	fmt.Fprintln(out, "\t"+`"net/http"`)
	fmt.Fprintln(out, "\t"+`"net/url"`)
	fmt.Fprintln(out, "\t"+`"slices"`)
	fmt.Fprintln(out, "\t"+`"strconv"`)
	fmt.Fprintln(out, `)`+"\n")

	methods := make(map[string][]Method)

	for _, f := range node.Decls {
		if g, ok := f.(*ast.GenDecl); ok {
			for _, spec := range g.Specs {
				currType, ok := spec.(*ast.TypeSpec)
				if !ok {
					fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
					continue
				}

				currStruct, ok := currType.Type.(*ast.StructType)
				if !ok {
					fmt.Printf("SKIP %#T is not ast.StructType\n", currStruct)
					continue
				}

				needDecl := true

			FIELDS_LOOP:
				for _, field := range currStruct.Fields.List {

					if field.Tag != nil {
						tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
						tagValue := tag.Get("apivalidator")
						if tagValue == "" {
							continue FIELDS_LOOP
						}

						if needDecl {
							fmt.Printf("process struct %s\n", currType.Name.Name)
							fmt.Printf("\tgenerating Decode method\n")

							fmt.Fprintf(out, "// %s validator\n", currType.Name.Name)
							fmt.Fprint(out, "func (out *"+currType.Name.Name+") ")
							fmt.Fprintln(out, "Decode(in map[string][]string) error {")
							fmt.Fprintln(out, "\tvar val []string\n\tvar exists bool")
							needDecl = false
						}

						fieldName := field.Names[0].Name
						fieldType := field.Type.(*ast.Ident).Name

						options := strings.Split(tagValue, ",")
						switch fieldType {
						case "int":
							decodeOps := decodeTpl[int]{
								FieldName: fieldName,
								ParamName: strings.ToLower(fieldName),
								Min:       math.MinInt,
								Max:       math.MaxInt,
							}
							for _, opt := range options {
								if val, ok := strings.CutPrefix(opt, "paramname="); ok {
									decodeOps.ParamName = val
									continue
								}
								if _, ok := strings.CutPrefix(opt, "required"); ok {
									decodeOps.Required = true
									continue
								}
								if val, ok := strings.CutPrefix(opt, "default="); ok {
									decodeOps.Default = val
									continue
								}
								if val, ok := strings.CutPrefix(opt, "enum="); ok {
									vals := strings.Split(val, "|")
									for _, v := range vals {
										i, err := strconv.Atoi(v)
										if err != nil {
											log.Panic(err)
										}
										decodeOps.Enum = append(decodeOps.Enum, i)
									}
									continue
								}
								if val, ok := strings.CutPrefix(opt, "min="); ok {
									decodeOps.Min, err = strconv.Atoi(val)
									if err != nil {
										log.Panic(err)
									}
									continue
								}
								if val, ok := strings.CutPrefix(opt, "max="); ok {
									decodeOps.Max, err = strconv.Atoi(val)
									if err != nil {
										log.Panic(err)
									}
									continue
								}
							}
							fmt.Fprintf(out, "\tval, exists = in[\"%s\"]\n", decodeOps.ParamName)
							if decodeOps.Required {
								fmt.Fprintf(out,
									"\tif !exists {\n\t\treturn fmt.Errorf(\"%s must me not empty\")\n\t}\n",
									strings.ToLower(decodeOps.FieldName))
							}
							fmt.Fprintf(out, "\tif exists {\n")
							fmt.Fprintf(out, "\t\t%s, err := strconv.Atoi(val[0])\n",
								strings.ToLower(decodeOps.FieldName))
							fmt.Fprintf(out, "\t\tif err != nil {\n")
							fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s must be int\")\n\t\t}\n",
								strings.ToLower(decodeOps.FieldName))
							fmt.Fprintf(out, "\t\tout.%s = %s\n",
								decodeOps.FieldName, strings.ToLower(decodeOps.FieldName))
							if decodeOps.Default != "" {
								def, err := strconv.Atoi(decodeOps.Default)
								if err != nil {
									log.Panic(err)
								}
								fmt.Fprintf(out, "\t\tif val[0] == \"\" {\n\t\t\tout.%s = %d\n\t\t} else {\n",
									decodeOps.FieldName, def)
								if len(decodeOps.Enum) != 0 {
									enum := "["
									for i, item := range decodeOps.Enum {
										if i != len(decodeOps.Enum)-1 {
											enum += fmt.Sprintf("%d , ", item)
											continue
										}
										enum += fmt.Sprintf("%d]", item)
									}
									fmt.Fprintf(out, "\t\t\topts := %#v\n", decodeOps.Enum)
									fmt.Fprintf(out, "\t\t\tif !slices.Contains(opts, out.%s) {\n", decodeOps.FieldName)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s must be one of %s\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), enum)
								}
								if decodeOps.Min > math.MinInt {
									fmt.Fprintf(out, "\t\t\tif out.%s < %d {\n", decodeOps.FieldName, decodeOps.Min)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s must be >= %d\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Min)
								}
								if decodeOps.Max < math.MaxInt {
									fmt.Fprintf(out, "\t\t\tif out.%s > %d {\n", decodeOps.FieldName, decodeOps.Max)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s must be <= %d\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Max)
								}
								fmt.Fprintln(out, "\t\t}")
							} else {
								if len(decodeOps.Enum) != 0 {
									enum := "["
									for i, item := range decodeOps.Enum {
										if i != len(decodeOps.Enum)-1 {
											enum += fmt.Sprintf("%d , ", item)
											continue
										}
										enum += fmt.Sprintf("%d]", item)
									}
									fmt.Fprintf(out, "\t\topts := %#v\n", decodeOps.Enum)
									fmt.Fprintf(out, "\t\tif !slices.Contains(opts, out.%s) {\n", decodeOps.FieldName)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s must be one of %s\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), enum)
								}
								if decodeOps.Min > math.MinInt {
									fmt.Fprintf(out, "\t\tif out.%s < %d {\n", decodeOps.FieldName, decodeOps.Min)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s must be >= %d\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Min)
								}
								if decodeOps.Max < math.MaxInt {
									fmt.Fprintf(out, "\t\tif out.%s > %d {\n", decodeOps.FieldName, decodeOps.Max)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s must be <= %d\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Max)
								}
							}
							fmt.Fprintln(out, "\t}")
						case "string":
							decodeOps := decodeTpl[string]{
								FieldName: fieldName,
								ParamName: strings.ToLower(fieldName),
								Min:       math.MinInt,
								Max:       math.MaxInt,
							}
							for _, opt := range options {
								if val, ok := strings.CutPrefix(opt, "paramname="); ok {
									decodeOps.ParamName = val
									continue
								}
								if _, ok := strings.CutPrefix(opt, "required"); ok {
									decodeOps.Required = true
									continue
								}
								if val, ok := strings.CutPrefix(opt, "default="); ok {
									decodeOps.Default = val
									continue
								}
								if val, ok := strings.CutPrefix(opt, "enum="); ok {
									decodeOps.Enum = strings.Split(val, "|")
									continue
								}
								if val, ok := strings.CutPrefix(opt, "min="); ok {
									decodeOps.Min, err = strconv.Atoi(val)
									if err != nil {
										log.Panic(err)
									}
									continue
								}
								if val, ok := strings.CutPrefix(opt, "max="); ok {
									decodeOps.Max, err = strconv.Atoi(val)
									if err != nil {
										log.Panic(err)
									}
									continue
								}
							}
							fmt.Fprintf(out, "\tval, exists = in[\"%s\"]\n", decodeOps.ParamName)
							if decodeOps.Required {
								fmt.Fprintf(out,
									"\tif !exists {\n\t\treturn fmt.Errorf(\"%s must me not empty\")\n\t}\n",
									strings.ToLower(decodeOps.FieldName))
							}
							fmt.Fprintf(out, "\tif exists {\n")
							fmt.Fprintf(out, "\t\tout.%s = val[0]\n", decodeOps.FieldName)
							if decodeOps.Default != "" {
								fmt.Fprintf(out, "\t\tif out.%s == \"\" {\n\t\t\tout.%s = \"%s\"\n\t\t} else {\n",
									decodeOps.FieldName, decodeOps.FieldName, decodeOps.Default)
								if len(decodeOps.Enum) != 0 {
									enum := "["
									for i, item := range decodeOps.Enum {
										if i != len(decodeOps.Enum)-1 {
											enum += item + ", "
											continue
										}
										enum += item + "]"
									}
									fmt.Fprintf(out, "\t\t\topts := %#v\n", decodeOps.Enum)
									fmt.Fprintf(out, "\t\t\tif !slices.Contains(opts, out.%s) {\n", decodeOps.FieldName)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s must be one of %s\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), enum)
								}
								if decodeOps.Min > math.MinInt {
									fmt.Fprintf(out, "\t\t\tif len(out.%s) < %d {\n", decodeOps.FieldName, decodeOps.Min)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s len must be >= %d\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Min)
								}
								if decodeOps.Max < math.MaxInt {
									fmt.Fprintf(out, "\t\t\tif len(out.%s) > %d {\n", decodeOps.FieldName, decodeOps.Max)
									fmt.Fprintf(out, "\t\t\t\treturn fmt.Errorf(\"%s len must be <= %d\")\n\t\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Max)
								}
								fmt.Fprintln(out, "\t\t}")
							} else {
								if len(decodeOps.Enum) != 0 {
									enum := "["
									for i, item := range decodeOps.Enum {
										if i != len(decodeOps.Enum)-1 {
											enum += item + ", "
											continue
										}
										enum += item + "]"
									}
									fmt.Fprintf(out, "\t\topts := %#v\n", decodeOps.Enum)
									fmt.Fprintf(out, "\t\tif !slices.Contains(opts, out.%s) {\n", decodeOps.FieldName)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s must be one of %s\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), enum)
								}
								if decodeOps.Min > math.MinInt {
									fmt.Fprintf(out, "\t\tif len(out.%s) < %d {\n", decodeOps.FieldName, decodeOps.Min)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s len must be >= %d\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Min)
								}
								if decodeOps.Max < math.MaxInt {
									fmt.Fprintf(out, "\t\tif len(out.%s) > %d {\n", decodeOps.FieldName, decodeOps.Max)
									fmt.Fprintf(out, "\t\t\treturn fmt.Errorf(\"%s len must be <= %d\")\n\t\t}\n",
										strings.ToLower(decodeOps.FieldName), decodeOps.Max)
								}
							}
							fmt.Fprintln(out, "\t}")
						}
					}
				}
				if !needDecl {
					fmt.Fprintln(out, "\treturn nil\n}\n")
				}
			}
		} else if fd, ok := f.(*ast.FuncDecl); ok {
			if fd.Doc == nil {
				fmt.Printf("SKIP function %#v: doesn't have comments\n", fd.Name.Name)
				continue
			}

			needCodegen := false
			var rcvType string
			for _, comment := range fd.Doc.List {
				if jsonStr, ok := strings.CutPrefix(comment.Text, "// apigen:api"); ok {
					needCodegen = true
					rcvType = fmt.Sprintf("%s", fd.Recv.List[0].Type.(*ast.StarExpr).X)
					inType := fd.Type.Params.List[len(fd.Type.Params.List)-1].Type.(*ast.Ident).Name
					outType := fd.Type.Results.List[0].Type.(*ast.StarExpr).X
					method := Method{
						MethodName: fd.Name.Name,
						InType:     inType,
						OutType:    outType,
					}
					err := json.Unmarshal([]byte(strings.TrimSpace(jsonStr)), &method)
					if err != nil {
						log.Panic(err)
					}
					methods[rcvType] = append(methods[rcvType], method)
					break
				}
			}
			if !needCodegen {
				fmt.Printf("SKIP function %#v: doesn't have apigen mark\n", fd.Name.Name)
				continue
			}

			fmt.Printf("process function %s\n", methods[rcvType][len(methods[rcvType])-1].MethodName)
			fmt.Printf("\tgenerating handler%s\n",
				methods[rcvType][len(methods[rcvType])-1].MethodName)
			fmt.Fprintf(out, "// handler for %s method\n",
				methods[rcvType][len(methods[rcvType])-1].MethodName)
			fmt.Fprintf(out, "func (srv *%s) handler%s(w http.ResponseWriter, r *http.Request) {\n",
				rcvType, methods[rcvType][len(methods[rcvType])-1].MethodName)

			tmpl := ""
			switch methods[rcvType][len(methods[rcvType])-1].MethodType {
			case http.MethodPost:
				tmpl = postTmpl
			case http.MethodGet:
				fallthrough
			case "":
				tmpl = getTmpl
			default:
				log.Panic("Unknown method type")
			}
			t := template.Must(template.ParseGlob(tmpl))
			t.Execute(out, methods[rcvType][len(methods[rcvType])-1])
			fmt.Fprintln(out, "\n")
		}
	}

	fmt.Printf("generate code for multiplexers\n")
	for mux, ms := range methods {
		fmt.Printf("\tgenerating code for %s\n", mux)
		fmt.Fprintf(out, "// multiplexer for %s\n", mux)
		fmt.Fprintf(out, "func (srv *%s) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n",
			mux)
		t := template.Must(template.ParseGlob(multiplexerTmpl))
		t.Execute(out, ms)
		fmt.Fprintln(out, "\n")
	}
}
