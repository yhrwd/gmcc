package nbt

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type PathNodeType int

// PathNode types
const (
	PathNodeRoot        PathNodeType = iota // {tags} - root compound
	PathNodeChild                           // name - child tag
	PathNodeCompound                        // name{tags} - compound child
	PathNodeAllElements                     // [] - all list elements
	PathNodeIndex                           // [index] - specific list element
	PathNodeFilter                          // [{tags}] - filtered list elements
)

// PathNode represents a single node in an NBT path
type PathNode struct {
	Type    PathNodeType
	Name    string         // for child, compound
	Index   int            // for index node
	Pattern map[string]any // for root, compound, filter
}

// Path represents a parsed NBT path
type Path struct {
	Nodes []PathNode
}

// ParsePath parses an NBT path string
func ParsePath(s string) (*Path, error) {
	p := &pathParser{input: s, pos: 0}
	return p.parse()
}

// QueryPath queries data using an NBT path string
func QueryPath(data any, path string) ([]any, error) {
	p, err := ParsePath(path)
	if err != nil {
		return nil, err
	}
	return p.Execute(data)
}

// QueryPathOne queries data and returns a single result
func QueryPathOne(data any, path string) (any, error) {
	results, err := QueryPath(data, path)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("multiple results found: %d", len(results))
	}
	return results[0], nil
}

// Execute executes the path on data
func (p *Path) Execute(data any) ([]any, error) {
	results := []any{data}

	for _, node := range p.Nodes {
		var newResults []any
		for _, item := range results {
			items, err := executeNode(node, item)
			if err != nil {
				return nil, err
			}
			newResults = append(newResults, items...)
		}
		results = newResults
		if len(results) == 0 {
			break
		}
	}

	return results, nil
}

// String returns the path string representation
func (p *Path) String() string {
	var sb strings.Builder
	for i, node := range p.Nodes {
		if i > 0 && node.Type != PathNodeAllElements && node.Type != PathNodeIndex && node.Type != PathNodeFilter {
			sb.WriteByte('.')
		}
		sb.WriteString(node.String())
	}
	return sb.String()
}

func (n PathNode) String() string {
	switch n.Type {
	case PathNodeRoot:
		if n.Pattern != nil && len(n.Pattern) > 0 {
			return FormatSNBT(n.Pattern)
		}
		return "{}"
	case PathNodeChild:
		return quoteName(n.Name)
	case PathNodeCompound:
		if n.Pattern != nil && len(n.Pattern) > 0 {
			return quoteName(n.Name) + FormatSNBT(n.Pattern)
		}
		return quoteName(n.Name) + "{}"
	case PathNodeAllElements:
		return "[]"
	case PathNodeIndex:
		return fmt.Sprintf("[%d]", n.Index)
	case PathNodeFilter:
		if n.Pattern != nil && len(n.Pattern) > 0 {
			return "[" + FormatSNBT(n.Pattern) + "]"
		}
		return "[{}]"
	default:
		return ""
	}
}

type pathParser struct {
	input string
	pos   int
}

func (p *pathParser) parse() (*Path, error) {
	path := &Path{}

	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("empty path")
	}

	// Check for root node
	if p.peek() == '{' {
		node, err := p.parseRootNode()
		if err != nil {
			return nil, err
		}
		path.Nodes = append(path.Nodes, node)
	}

	for p.pos < len(p.input) {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}

		// Skip dot separator
		if p.peek() == '.' {
			p.pos++
			p.skipWhitespace()
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		path.Nodes = append(path.Nodes, node)
	}

	return path, nil
}

func (p *pathParser) parseRootNode() (PathNode, error) {
	if p.peek() != '{' {
		return PathNode{}, fmt.Errorf("expected '{'")
	}
	p.pos++

	pattern, err := p.parseCompoundPattern()
	if err != nil {
		return PathNode{}, err
	}

	return PathNode{Type: PathNodeRoot, Pattern: pattern}, nil
}

func (p *pathParser) parseNode() (PathNode, error) {
	if p.pos >= len(p.input) {
		return PathNode{}, fmt.Errorf("unexpected end of path")
	}

	ch := p.peek()

	// List operations
	if ch == '[' {
		return p.parseListNode()
	}

	// Quoted or unquoted name
	name, err := p.parseName()
	if err != nil {
		return PathNode{}, err
	}

	// Check for compound pattern
	p.skipWhitespace()
	if p.pos < len(p.input) && p.peek() == '{' {
		p.pos++
		pattern, err := p.parseCompoundPattern()
		if err != nil {
			return PathNode{}, err
		}
		return PathNode{Type: PathNodeCompound, Name: name, Pattern: pattern}, nil
	}

	return PathNode{Type: PathNodeChild, Name: name}, nil
}

func (p *pathParser) parseListNode() (PathNode, error) {
	if p.peek() != '[' {
		return PathNode{}, fmt.Errorf("expected '['")
	}
	p.pos++

	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return PathNode{}, fmt.Errorf("unclosed '['")
	}

	ch := p.peek()

	// All elements: []
	if ch == ']' {
		p.pos++
		return PathNode{Type: PathNodeAllElements}, nil
	}

	// Filter: [{tags}]
	if ch == '{' {
		p.pos++
		pattern, err := p.parseCompoundPattern()
		if err != nil {
			return PathNode{}, err
		}
		if p.peek() != ']' {
			return PathNode{}, fmt.Errorf("expected ']'")
		}
		p.pos++
		return PathNode{Type: PathNodeFilter, Pattern: pattern}, nil
	}

	// Index: [number]
	index, err := p.parseNumber()
	if err != nil {
		return PathNode{}, err
	}

	p.skipWhitespace()
	if p.peek() != ']' {
		return PathNode{}, fmt.Errorf("expected ']'")
	}
	p.pos++

	return PathNode{Type: PathNodeIndex, Index: index}, nil
}

func (p *pathParser) parseCompoundPattern() (map[string]any, error) {
	pattern := make(map[string]any)

	p.skipWhitespace()
	if p.peek() == '}' {
		p.pos++
		return pattern, nil
	}

	// Parse the pattern content using SNBT parser
	// For simplicity, we use a simplified version
	start := p.pos - 1 // include the opening brace
	braceCount := 1

	for p.pos < len(p.input) && braceCount > 0 {
		ch := p.input[p.pos]
		if ch == '{' {
			braceCount++
		} else if ch == '}' {
			braceCount--
		} else if ch == '"' || ch == '\'' {
			// Skip quoted strings
			quote := ch
			p.pos++
			for p.pos < len(p.input) && p.input[p.pos] != quote {
				if p.input[p.pos] == '\\' {
					p.pos++
				}
				p.pos++
			}
		}
		p.pos++
	}

	if braceCount > 0 {
		return nil, fmt.Errorf("unclosed '{'")
	}

	snbt := p.input[start:p.pos]

	// Parse SNBT pattern
	result, err := ParseSNBT(snbt)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("pattern must be a compound")
	}

	return m, nil
}

func (p *pathParser) parseName() (string, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return "", fmt.Errorf("expected name")
	}

	ch := p.peek()

	// Quoted name
	if ch == '"' || ch == '\'' {
		return p.parseQuotedName()
	}

	// Unquoted name
	return p.parseUnquotedName()
}

func (p *pathParser) parseQuotedName() (string, error) {
	quote := p.peek()
	p.pos++

	var sb strings.Builder
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == quote {
			p.pos++
			return sb.String(), nil
		}
		if ch == '\\' && p.pos+1 < len(p.input) {
			p.pos++
			sb.WriteByte(p.input[p.pos])
		} else {
			sb.WriteByte(ch)
		}
		p.pos++
	}

	return "", fmt.Errorf("unclosed quote")
}

func (p *pathParser) parseUnquotedName() (string, error) {
	var sb strings.Builder

	// Allow namespaced IDs (minecraft:stone)
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == '.' || ch == '[' || ch == '{' || unicode.IsSpace(rune(ch)) {
			break
		}
		sb.WriteByte(ch)
		p.pos++
	}

	name := sb.String()
	if name == "" {
		return "", fmt.Errorf("expected name")
	}

	return name, nil
}

func (p *pathParser) parseNumber() (int, error) {
	p.skipWhitespace()

	start := p.pos
	if p.pos < len(p.input) && p.input[p.pos] == '-' {
		p.pos++
	}

	for p.pos < len(p.input) && unicode.IsDigit(rune(p.input[p.pos])) {
		p.pos++
	}

	if start == p.pos {
		return 0, fmt.Errorf("expected number")
	}

	n, err := strconv.Atoi(p.input[start:p.pos])
	return n, err
}

func (p *pathParser) peek() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *pathParser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func quoteName(name string) string {
	needsQuote := false
	for _, ch := range name {
		if ch == '.' || ch == ' ' || ch == '"' || ch == '\'' || ch == '[' || ch == ']' || ch == '{' || ch == '}' {
			needsQuote = true
			break
		}
	}

	if !needsQuote {
		return name
	}

	var sb strings.Builder
	sb.WriteByte('"')
	for _, ch := range name {
		if ch == '"' || ch == '\\' {
			sb.WriteByte('\\')
		}
		sb.WriteRune(ch)
	}
	sb.WriteByte('"')
	return sb.String()
}

func executeNode(node PathNode, data any) ([]any, error) {
	type nodeExecutor func(PathNode, any) ([]any, error)
	executors := map[PathNodeType]nodeExecutor{
		PathNodeRoot:        executeRoot,
		PathNodeChild:       executeChild,
		PathNodeCompound:    executeCompound,
		PathNodeAllElements: executeAllElements,
		PathNodeIndex:       executeIndex,
		PathNodeFilter:      executeFilter,
	}

	if executor, ok := executors[node.Type]; ok {
		return executor(node, data)
	}
	return nil, fmt.Errorf("Unknown path node type: %d", node.Type)
}

func executeRoot(node PathNode, data any) ([]any, error) {
	if node.Pattern == nil {
		return []any{data}, nil
	}
	matched, err := matchPattern(data, node.Pattern)
	if err != nil {
		return nil, err
	}
	if matched {
		return []any{data}, nil
	}
	return []any{}, nil
}

func executeChild(node PathNode, data any) ([]any, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected compound, got %T", data)
	}

	value, exists := m[node.Name]
	if !exists {
		return []any{}, nil
	}

	return []any{value}, nil
}

func executeCompound(node PathNode, data any) ([]any, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected compound, got %T", data)
	}

	value, exists := m[node.Name]
	if !exists {
		return []any{}, nil
	}

	if node.Pattern == nil {
		return []any{value}, nil
	}

	matched, err := matchPattern(value, node.Pattern)
	if err != nil {
		return nil, err
	}
	if matched {
		return []any{value}, nil
	}

	return []any{}, nil
}

func executeAllElements(node PathNode, data any) ([]any, error) {
	list, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list, got %T", data)
	}

	return list, nil
}

func executeIndex(node PathNode, data any) ([]any, error) {
	list, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list, got %T", data)
	}

	if node.Index < 0 {
		node.Index += len(list)
	}

	if node.Index < 0 || node.Index >= len(list) {
		return []any{}, nil
	}

	return []any{list[node.Index]}, nil
}

func executeFilter(node PathNode, data any) ([]any, error) {
	list, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list, got %T", data)
	}

	var results []any

	for _, item := range list {
		matched, err := matchPattern(item, node.Pattern)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results, item)
		}
	}

	return results, nil
}

func matchPattern(data any, pattern map[string]any) (bool, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return false, nil
	}

	for key, expected := range pattern {
		actual, exists := m[key]
		if !exists {
			return false, nil
		}

		expectedMap, expectedIsMap := expected.(map[string]any)
		if expectedIsMap {
			matched, err := matchPattern(actual, expectedMap)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
			continue
		}

		if !reflect.DeepEqual(actual, expected) {
			return false, nil
		}
	}

	return true, nil
}
