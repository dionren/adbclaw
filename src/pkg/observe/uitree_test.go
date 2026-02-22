package observe

import (
	"testing"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<hierarchy rotation="0">
  <node index="0" text="" resource-id="" class="android.widget.FrameLayout" package="com.example" content-desc="" checkable="false" checked="false" clickable="false" enabled="true" focusable="false" focused="false" scrollable="false" selected="false" bounds="[0,0][1080,2400]">
    <node index="0" text="Settings" resource-id="com.example:id/title" class="android.widget.TextView" package="com.example" content-desc="" checkable="false" checked="false" clickable="true" enabled="true" focusable="true" focused="false" scrollable="false" selected="false" bounds="[0,120][1080,180]">
    </node>
    <node index="1" text="" resource-id="com.example:id/icon" class="android.widget.ImageView" package="com.example" content-desc="App icon" checkable="false" checked="false" clickable="true" enabled="true" focusable="false" focused="false" scrollable="false" selected="false" bounds="[100,200][200,300]">
    </node>
    <node index="2" text="" resource-id="" class="android.view.View" package="com.example" content-desc="" checkable="false" checked="false" clickable="false" enabled="true" focusable="false" focused="false" scrollable="false" selected="false" bounds="[0,300][1080,400]">
    </node>
    <node index="3" text="Login" resource-id="com.example:id/btn_login" class="android.widget.Button" package="com.example" content-desc="" checkable="false" checked="false" clickable="true" enabled="true" focusable="true" focused="false" scrollable="false" selected="false" bounds="[200,500][880,600]">
    </node>
    <node index="4" text="" resource-id="" class="android.widget.ScrollView" package="com.example" content-desc="" checkable="false" checked="false" clickable="false" enabled="true" focusable="false" focused="false" scrollable="true" selected="false" bounds="[0,600][1080,2400]">
    </node>
  </node>
</hierarchy>`

func TestParseUITree(t *testing.T) {
	tree, err := ParseUITree([]byte(sampleXML))
	if err != nil {
		t.Fatalf("ParseUITree failed: %v", err)
	}

	// The plain View (index 2) has no text/resource-id/content-desc and is not clickable/scrollable,
	// so it should be filtered out. We expect 4 elements.
	if len(tree.Elements) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(tree.Elements))
	}

	// Check first element
	el := tree.Elements[0]
	if el.Text != "Settings" {
		t.Errorf("element 0 text = %q, want %q", el.Text, "Settings")
	}
	if el.ResourceID != "com.example:id/title" {
		t.Errorf("element 0 resource_id = %q, want %q", el.ResourceID, "com.example:id/title")
	}
	if el.Center.X != 540 || el.Center.Y != 150 {
		t.Errorf("element 0 center = (%d,%d), want (540,150)", el.Center.X, el.Center.Y)
	}
	if !el.Clickable {
		t.Error("element 0 should be clickable")
	}

	// Check sequential indexing
	for i, el := range tree.Elements {
		if el.Index != i {
			t.Errorf("element %d has index %d", i, el.Index)
		}
	}
}

func TestFindByText(t *testing.T) {
	tree, _ := ParseUITree([]byte(sampleXML))

	results := tree.FindByText("login")
	if len(results) != 1 {
		t.Fatalf("FindByText('login') returned %d results, want 1", len(results))
	}
	if results[0].Text != "Login" {
		t.Errorf("found text = %q, want %q", results[0].Text, "Login")
	}
}

func TestFindByID(t *testing.T) {
	tree, _ := ParseUITree([]byte(sampleXML))

	results := tree.FindByID("btn_login")
	if len(results) != 1 {
		t.Fatalf("FindByID('btn_login') returned %d results, want 1", len(results))
	}
	if results[0].ResourceID != "com.example:id/btn_login" {
		t.Errorf("found id = %q", results[0].ResourceID)
	}
}

func TestFindByIndex(t *testing.T) {
	tree, _ := ParseUITree([]byte(sampleXML))

	el, err := tree.FindByIndex(0)
	if err != nil {
		t.Fatalf("FindByIndex(0) failed: %v", err)
	}
	if el.Text != "Settings" {
		t.Errorf("element 0 text = %q, want %q", el.Text, "Settings")
	}

	_, err = tree.FindByIndex(99)
	if err == nil {
		t.Error("FindByIndex(99) should return error")
	}
}

func TestFindByContentDesc(t *testing.T) {
	tree, _ := ParseUITree([]byte(sampleXML))

	results := tree.FindByText("App icon")
	if len(results) != 1 {
		t.Fatalf("FindByText('App icon') returned %d results, want 1", len(results))
	}
	if results[0].ContentDesc != "App icon" {
		t.Errorf("found content_desc = %q", results[0].ContentDesc)
	}
}

func TestParseBounds(t *testing.T) {
	tests := []struct {
		input string
		want  Bounds
		err   bool
	}{
		{"[0,120][1080,180]", Bounds{0, 120, 1080, 180}, false},
		{"[100,200][200,300]", Bounds{100, 200, 200, 300}, false},
		{"invalid", Bounds{}, true},
	}

	for _, tt := range tests {
		got, err := parseBounds(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("parseBounds(%q) error = %v, want error = %v", tt.input, err, tt.err)
			continue
		}
		if !tt.err && got != tt.want {
			t.Errorf("parseBounds(%q) = %+v, want %+v", tt.input, got, tt.want)
		}
	}
}
