package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser golang file
type Parser struct {
	filePath    string
	modelParser *modelParser
	fileSet     *token.FileSet
}

// New parser instance
func New(filePath string) *Parser {
	return &Parser{
		modelParser: newModelParser(),
		filePath:    filePath,
	}
}

// Parse go source file
func (p *Parser) Parse() (*Context, error) {
	p.fileSet = token.NewFileSet()

	// parse file
	parsed, err := parser.ParseFile(p.fileSet, p.filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	ctx.PkgName = parsed.Name.Name

	ast.Inspect(parsed, p.InspectTypes(ctx))
	ast.Inspect(parsed, p.InspectVars(ctx))

	return ctx, nil
}

// InspectTypes declarations
func (p *Parser) InspectTypes(ctx *Context) func(ast.Node) bool {
	return func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			return true
		}

		return p.modelParser.InspectTypeSpecs(ctx, gd.Specs)
	}
}

// InspectVars declarations
func (p *Parser) InspectVars(ctx *Context) func(ast.Node) bool {
	return func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			return true
		}
		return p.modelParser.InspectVarsSpecs(ctx, gd.Specs)
	}
}
