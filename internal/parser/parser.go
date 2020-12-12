package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser golang file
type Parser struct {
	filePath string
	fileSet  *token.FileSet
}

// New parser instance
func New(filePath string) *Parser {
	return &Parser{
		filePath: filePath,
	}
}

// Parse go source file
func (p *Parser) Parse() (*Context, error) {
	p.fileSet = token.NewFileSet()

	ctx := &Context{}
	ctx.DefList = make(map[string]*ModelDef)
	ctx.ModelImport = "model"
	ctx.FieldImport = "field"

	// parse file
	parsed, err := parser.ParseFile(p.fileSet, p.filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	ctx.PkgName = parsed.Name.Name
	ctx.SchemaPkgNames = make(map[string]struct{})

	ast.Inspect(parsed, p.InspectImports(ctx))
	ast.Inspect(parsed, p.InspectTypes(ctx))
	ast.Inspect(parsed, p.InspectFuncs(ctx))

	return ctx, nil
}

// InspectImports declarations
func (p *Parser) InspectImports(ctx *Context) func(ast.Node) bool {
	return func(n ast.Node) bool {
		is, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}

		if is.Path.Value == `"github.com/gigovich/ddb/schema"` {
			if is.Name != nil {
				ctx.SchemaPkgNames[is.Name.String()] = struct{}{}
			} else {
				ctx.SchemaPkgNames["schema"] = struct{}{}
			}
		}
		return true
	}
}

// InspectTypes declarations
func (p *Parser) InspectTypes(ctx *Context) func(ast.Node) bool {
	return func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			return true
		}
		return false
	}
}

// InspectFuncs declarations
func (p *Parser) InspectFuncs(ctx *Context) func(ast.Node) bool {
	return func(n ast.Node) bool {
		gd, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		return p.inspectSchemaSpecs(ctx, gd)
	}
}

func (p *Parser) inspectSchemaSpecs(ctx *Context, decl *ast.FuncDecl) bool {
	if decl.Name.String() != "Schema" {
		return false
	}

	for _, r := range decl.Type.Results.List {
		se, ok := r.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		xIdent, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}

		if _, ok := ctx.SchemaPkgNames[xIdent.Name]; ok && se.Sel.Name == "Table" {
			fmt.Println(">>>>>>>>>")
			return true
		}
	}

	return true
}
