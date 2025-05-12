package core

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

type DynamicRouter struct {
	apiHandlers []interface{}
	Routes      []RouteConfig
}

func (b *DynamicRouter) add(method, path string, handler Handler) {
	b.Routes = append(b.Routes, RouteConfig{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (b *DynamicRouter) SetHandlers(handlers ...interface{}) {
	b.apiHandlers = append(b.apiHandlers, handlers...)
}

func (b *DynamicRouter) RoutersPath(dirs ...string) {
	if len(b.apiHandlers) == 0 {
		log.Println("no registered handlers found; did you forget SetHandlers()?")
		return
	}

	ctxType := reflect.TypeOf((*Context)(nil)).Elem()
	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.go"))
		if err != nil {
			log.Printf("failed to scan folder %s: %v", dir, err)
			continue
		}

		for _, file := range files {
			routes := ParseRoutesFromFile(file)

			for _, route := range routes {
				found := false

				for _, handler := range b.apiHandlers {
					val := reflect.ValueOf(handler)
					method := val.MethodByName(route.Handler)
					if !method.IsValid() {
						log.Printf("method %s not found in handler %T", route.Handler, handler)
						continue
					}

					// Kiểm tra kiểu tham số đầu vào và kiểu trả về
					methodType := method.Type()
					if methodType.NumIn() != 1 || methodType.In(0) != ctxType || methodType.NumOut() != 0 {
						log.Printf("method %s in %T must be func(Context)", route.Handler, handler)
						continue
					}
					// Tạo handler thực thi
					h := Handler(func(ctx Context) {
						results := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
						if len(results) > 0 && !results[0].IsNil() {
							if err, ok := results[0].Interface().(error); ok {
								log.Printf("handler %s returned error: %v", route.Handler, err)
							}
						}
					})

					b.add(route.Method, route.Path, h)
					found = true
					break
				}

				if !found {
					log.Printf("handler method %s not found for path %s", route.Handler, route.Path)
				}
			}
		}
	}
}

type ParseRoute struct {
	Path    string
	Method  string
	Handler string
}

func ParseRoutesFromFile(filename string) []ParseRoute {
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	var routes []ParseRoute
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}
		for _, comment := range fn.Doc.List {
			if strings.HasPrefix(comment.Text, "// @route") {
				// Tách comment thành các phần tử
				parts := strings.Fields(comment.Text)
				if len(parts) != 4 {
					log.Printf("Invalid @route comment format: %s", comment.Text)
					continue
				}
				// Phần tử đầu tiên là @route, phần tử thứ 2 là method, phần tử thứ 3 là path
				method := parts[2]
				path := parts[3]

				routes = append(routes, ParseRoute{
					Method:  method,
					Path:    path,
					Handler: fn.Name.Name,
				})
			}
		}
	}
	return routes
}
