package nbt

import (
	"testing"
)

func TestParsePath_Simple(t *testing.T) {
	tests := []struct {
		path     string
		wantErr  bool
		nodeType int
	}{
		{"foo", false, PathNodeChild},
		{"foo.bar", false, PathNodeChild},
		{"foo[0]", false, PathNodeIndex},
		{"foo[]", false, PathNodeAllElements},
		{"foo[{}]", false, PathNodeFilter},
		{"foo{bar:1b}", false, PathNodeCompound},
		{"{}", false, PathNodeRoot},
		{"{foo:1b}", false, PathNodeRoot},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			p, err := ParsePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(p.Nodes) > 0 {
				if p.Nodes[0].Type != tt.nodeType && p.Nodes[len(p.Nodes)-1].Type != tt.nodeType {
					t.Errorf("ParsePath(%q) node type = %v, want %v", tt.path, p.Nodes[0].Type, tt.nodeType)
				}
			}
		})
	}
}

func TestParsePath_QuotedNames(t *testing.T) {
	tests := []string{
		`"a.b"`,
		`'a.b'`,
		`"a b"`,
		`"a[0]"`,
		`foo."bar baz"`,
		`foo.'a.b'`,
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			p, err := ParsePath(path)
			if err != nil {
				t.Errorf("ParsePath(%q) error = %v", path, err)
				return
			}
			if len(p.Nodes) == 0 {
				t.Errorf("ParsePath(%q) returned empty path", path)
			}
		})
	}
}

func TestQueryPath_Simple(t *testing.T) {
	data := map[string]any{
		"foo": map[string]any{
			"bar": int32(42),
		},
		"list": []any{int32(1), int32(2), int32(3)},
	}

	tests := []struct {
		path    string
		wantLen int
		wantVal any
	}{
		{"foo", 1, nil},
		{"foo.bar", 1, int32(42)},
		{"list", 1, nil},
		{"list[]", 3, nil},
		{"list[0]", 1, int32(1)},
		{"list[2]", 1, int32(3)},
		{"list[-1]", 1, int32(3)},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			results, err := QueryPath(data, tt.path)
			if err != nil {
				t.Errorf("QueryPath(%q) error = %v", tt.path, err)
				return
			}
			if len(results) != tt.wantLen {
				t.Errorf("QueryPath(%q) got %d results, want %d", tt.path, len(results), tt.wantLen)
			}
			if tt.wantVal != nil && len(results) > 0 {
				if results[0] != tt.wantVal {
					t.Errorf("QueryPath(%q) got %v, want %v", tt.path, results[0], tt.wantVal)
				}
			}
		})
	}
}

func TestQueryPath_Compound(t *testing.T) {
	data := map[string]any{
		"VillagerData": map[string]any{
			"type":  "plains",
			"level": int32(5),
		},
	}

	results, err := QueryPath(data, "VillagerData{type:\"plains\"}")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestQueryPath_Filter(t *testing.T) {
	data := map[string]any{
		"Inventory": []any{
			map[string]any{"Slot": int8(0), "id": "diamond"},
			map[string]any{"Slot": int8(1), "id": "iron"},
			map[string]any{"Slot": int8(2), "id": "gold"},
		},
	}

	results, err := QueryPath(data, "Inventory[{Slot:1b}]")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestQueryPath_Root(t *testing.T) {
	data := map[string]any{
		"Invisible": int8(1),
		"name":      "test",
	}

	results, err := QueryPath(data, "{Invisible:1b}")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	results, err = QueryPath(data, "{Invisible:0b}")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestQueryPath_Nested(t *testing.T) {
	data := map[string]any{
		"foo": map[string]any{
			"bar": []any{
				map[string]any{"baz": int32(1)},
				map[string]any{"baz": int32(2)},
			},
		},
	}

	results, err := QueryPath(data, "foo.bar[].baz")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestQueryPath_NamespacedID(t *testing.T) {
	data := map[string]any{
		"components": map[string]any{
			"minecraft:written_book_content": map[string]any{
				"pages": []any{
					map[string]any{"raw": "page1"},
					map[string]any{"raw": "page2"},
				},
			},
		},
	}

	results, err := QueryPath(data, "components.minecraft:written_book_content.pages[0].raw")
	if err != nil {
		t.Fatalf("QueryPath error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0] != "page1" {
		t.Errorf("expected 'page1', got %v", results[0])
	}
}

func TestPath_String(t *testing.T) {
	tests := []string{
		"foo",
		"foo.bar",
		"foo[0]",
		"foo[]",
		"{}",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			p, err := ParsePath(path)
			if err != nil {
				t.Fatalf("ParsePath error: %v", err)
			}
			got := p.String()
			if got != path {
				t.Errorf("Path.String() = %q, want %q", got, path)
			}
		})
	}
}

func TestPath_RoundTrip(t *testing.T) {
	original := "{foo:1b}"
	p, err := ParsePath(original)
	if err != nil {
		t.Fatalf("ParsePath error: %v", err)
	}
	results, err := QueryPath(map[string]any{"foo": int8(1)}, p.String())
	if err != nil {
		t.Fatalf("QueryPath with formatted path error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
