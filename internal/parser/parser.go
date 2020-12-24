package parser

import (
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

		if is.Path.Value == `"github.com/gigovich/ddb/dsl"` {
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

		return p.inspectSchemaSpecsFunc(ctx, gd)
	}
}

func (p *Parser) inspectSchemaSpecsFunc(ctx *Context, decl *ast.FuncDecl) bool {
	// check functio name: Schema
	if decl.Name.String() != "Schema" {
		return false
	}

	// check return type: dsl.Table
	if len(decl.Type.Results.List) != 1 {
		return false
	}

	r := decl.Type.Results.List[0]
	se, ok := r.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	xIdent, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}

	_, ok := ctx.SchemaPkgNames[xIdent.Name]
	if !ok || se.Sel.Name != "Table" {
		return false
	}

	// check reciever type: model.Name
	if len(decl.Recv.List) != 1 && len(decl.Recv.List[0].Names) != 1 {
		return false
	}

	if it, ok := decl.Recv.List[0].Type.(*ast.Ident); !ok || it.Name != ctx.ModelType {
		return false
	}

	rc, ok := decl.Recv.List[0].Names[0].(*ast.Ident)
	if !ok {
		return false
	}

	return false
}
