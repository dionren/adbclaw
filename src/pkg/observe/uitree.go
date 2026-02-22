package observe

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/llm-net/adbclaw/pkg/adb"
)

// Bounds represents the bounding box of a UI element.
type Bounds struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

// Point represents a coordinate.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Element is a filtered, indexed UI element from the hierarchy.
type Element struct {
	Index       int    `json:"index"`
	Class       string `json:"class"`
	ResourceID  string `json:"resource_id"`
	Text        string `json:"text"`
	ContentDesc string `json:"content_desc"`
	Bounds      Bounds `json:"bounds"`
	Center      Point  `json:"center"`
	Clickable   bool   `json:"clickable"`
	Scrollable  bool   `json:"scrollable"`
	Focusable   bool   `json:"focusable"`
	Enabled     bool   `json:"enabled"`
	Selected    bool   `json:"selected"`
	Checked     bool   `json:"checked"`
	PackageName string `json:"package,omitempty"`
}

// UITree holds the parsed UI hierarchy.
type UITree struct {
	Elements []Element `json:"elements"`
}

// xmlNode represents a node in the uiautomator XML dump.
type xmlNode struct {
	XMLName     xml.Name  `xml:"node"`
	Index       string    `xml:"index,attr"`
	Text        string    `xml:"text,attr"`
	ResourceID  string    `xml:"resource-id,attr"`
	Class       string    `xml:"class,attr"`
	Package     string    `xml:"package,attr"`
	ContentDesc string    `xml:"content-desc,attr"`
	Checkable   string    `xml:"checkable,attr"`
	Checked     string    `xml:"checked,attr"`
	Clickable   string    `xml:"clickable,attr"`
	Enabled     string    `xml:"enabled,attr"`
	Focusable   string    `xml:"focusable,attr"`
	Focused     string    `xml:"focused,attr"`
	Scrollable  string    `xml:"scrollable,attr"`
	Selected    string    `xml:"selected,attr"`
	Bounds      string    `xml:"bounds,attr"`
	Children    []xmlNode `xml:"node"`
}

type xmlHierarchy struct {
	XMLName  xml.Name  `xml:"hierarchy"`
	Rotation string    `xml:"rotation,attr"`
	Nodes    []xmlNode `xml:"node"`
}

var boundsRe = regexp.MustCompile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)

func parseBounds(s string) (Bounds, error) {
	m := boundsRe.FindStringSubmatch(s)
	if len(m) != 5 {
		return Bounds{}, fmt.Errorf("invalid bounds format: %s", s)
	}
	l, _ := strconv.Atoi(m[1])
	t, _ := strconv.Atoi(m[2])
	r, _ := strconv.Atoi(m[3])
	b, _ := strconv.Atoi(m[4])
	return Bounds{Left: l, Top: t, Right: r, Bottom: b}, nil
}

// isSignificant returns true if a node is "meaningful" for agent interaction.
func isSignificant(n *xmlNode) bool {
	if n.Text != "" || n.ResourceID != "" || n.ContentDesc != "" {
		return true
	}
	if n.Clickable == "true" || n.Scrollable == "true" {
		return true
	}
	return false
}

// DumpUITree runs "uiautomator dump" and returns the parsed XML.
func DumpUITree(cmd adb.Commander) (*UITree, error) {
	// Dump UI hierarchy to a file on device, then read it
	result, err := cmd.Shell("uiautomator", "dump", "/dev/tty")
	if err != nil {
		return nil, fmt.Errorf("uiautomator dump failed: %w", err)
	}

	xmlData := result.Stdout
	// uiautomator dump /dev/tty outputs the XML directly to stdout,
	// but may prefix with "UI hierchary dumped to: /dev/tty" or similar.
	// Find the XML start.
	idx := strings.Index(xmlData, "<?xml")
	if idx < 0 {
		// Try finding hierarchy tag directly
		idx = strings.Index(xmlData, "<hierarchy")
	}
	if idx < 0 {
		return nil, fmt.Errorf("uiautomator dump returned no XML data: %s", truncate(xmlData, 200))
	}
	xmlData = xmlData[idx:]

	return ParseUITree([]byte(xmlData))
}

// ParseUITree parses uiautomator XML into a UITree with indexed elements.
func ParseUITree(data []byte) (*UITree, error) {
	var h xmlHierarchy
	if err := xml.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("xml parse error: %w", err)
	}

	var elements []Element
	var walk func(nodes []xmlNode)
	walk = func(nodes []xmlNode) {
		for i := range nodes {
			n := &nodes[i]
			if isSignificant(n) {
				bounds, err := parseBounds(n.Bounds)
				if err == nil {
					el := Element{
						Index:       len(elements),
						Class:       n.Class,
						ResourceID:  n.ResourceID,
						Text:        n.Text,
						ContentDesc: n.ContentDesc,
						Bounds:      bounds,
						Center: Point{
							X: (bounds.Left + bounds.Right) / 2,
							Y: (bounds.Top + bounds.Bottom) / 2,
						},
						Clickable:   n.Clickable == "true",
						Scrollable:  n.Scrollable == "true",
						Focusable:   n.Focusable == "true",
						Enabled:     n.Enabled == "true",
						Selected:    n.Selected == "true",
						Checked:     n.Checked == "true",
						PackageName: n.Package,
					}
					elements = append(elements, el)
				}
			}
			walk(n.Children)
		}
	}
	walk(h.Nodes)

	return &UITree{Elements: elements}, nil
}

// FindByIndex returns the element at the given index.
func (t *UITree) FindByIndex(index int) (*Element, error) {
	if index < 0 || index >= len(t.Elements) {
		return nil, fmt.Errorf("index %d out of range (0-%d)", index, len(t.Elements)-1)
	}
	return &t.Elements[index], nil
}

// FindByText returns elements whose text contains the query (case-insensitive).
func (t *UITree) FindByText(query string) []Element {
	query = strings.ToLower(query)
	var results []Element
	for _, el := range t.Elements {
		if strings.Contains(strings.ToLower(el.Text), query) ||
			strings.Contains(strings.ToLower(el.ContentDesc), query) {
			results = append(results, el)
		}
	}
	return results
}

// FindByID returns elements whose resource-id contains the query.
func (t *UITree) FindByID(query string) []Element {
	query = strings.ToLower(query)
	var results []Element
	for _, el := range t.Elements {
		if strings.Contains(strings.ToLower(el.ResourceID), query) {
			results = append(results, el)
		}
	}
	return results
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
