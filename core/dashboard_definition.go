package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	dashboardColumnCount    = 16
	dashboardMaxWidgets     = 200
	dashboardMaxQueryLength = 8000
	dashboardMaxTextLength  = 20000
)

type dashboardDefinition struct {
	Version int               `json:"version"`
	Query   string            `json:"query"`
	Widgets []dashboardWidget `json:"widgets"`
}

type dashboardWidget struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Layout     dashboardWidgetLayout  `json:"layout"`
	Query      string                 `json:"query"`
	Properties map[string]interface{} `json:"properties"`
}

type dashboardWidgetLayout struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func validateDashboardDefinition(raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		return invalidDefinition("dashboard definition must be valid JSON")
	}
	root, ok := value.(map[string]interface{})
	if !ok {
		return invalidDefinition("dashboard definition must be an object")
	}
	if len(root) == 0 {
		return nil
	}
	if _, ok := root["query"].(string); !ok {
		return invalidDefinition("dashboard query must be a string")
	}
	version, ok := root["version"].(json.Number)
	if !ok || version.String() != "1" {
		return invalidDefinition("unsupported dashboard definition version")
	}
	var definition dashboardDefinition
	if err := json.Unmarshal(raw, &definition); err != nil {
		return invalidDefinition("dashboard definition is invalid")
	}
	if len(definition.Query) > dashboardMaxQueryLength {
		return invalidDefinition("dashboard query exceeds 8000 characters")
	}
	if _, ok := root["widgets"].([]interface{}); !ok {
		return invalidDefinition("dashboard widgets must be an array")
	}
	if len(definition.Widgets) > dashboardMaxWidgets {
		return invalidDefinition("dashboard cannot contain more than 200 widgets")
	}
	ids := make(map[string]struct{}, len(definition.Widgets))
	for i, widget := range definition.Widgets {
		if widget.ID == "" {
			return invalidDefinition("widget ID is required")
		}
		if _, exists := ids[widget.ID]; exists {
			return invalidDefinition("widget IDs must be unique")
		}
		ids[widget.ID] = struct{}{}
		if widget.Type != "text" {
			return invalidDefinition(fmt.Sprintf("unknown widget type %q", widget.Type))
		}
		widgetValue, ok := rootWidgetAt(root, i)
		if !ok {
			return invalidDefinition("dashboard widget is invalid")
		}
		if _, ok := widgetValue["query"].(string); !ok {
			return invalidDefinition("widget query must be a string")
		}
		if _, ok := widgetValue["properties"].(map[string]interface{}); !ok {
			return invalidDefinition("widget properties must be an object")
		}
		layout := widget.Layout
		if layout.X < 0 || layout.X+layout.Width > dashboardColumnCount || layout.Y < 0 || layout.Width < 2 || layout.Height < 2 {
			return invalidDefinition("widget layout is outside the grid or below its minimum size")
		}
		if len(widget.Query) > dashboardMaxQueryLength {
			return invalidDefinition("widget query exceeds 8000 characters")
		}
		text, ok := widget.Properties["text"].(string)
		if !ok {
			return invalidDefinition("text widget properties.text must be a string")
		}
		if len(text) > dashboardMaxTextLength {
			return invalidDefinition("text widget text exceeds 20000 characters")
		}
		for j := 0; j < i; j++ {
			if dashboardLayoutsOverlap(layout, definition.Widgets[j].Layout) {
				return invalidDefinition("widgets cannot overlap")
			}
		}
	}
	return nil
}

func rootWidgetAt(root map[string]interface{}, index int) (map[string]interface{}, bool) {
	widgets, ok := root["widgets"].([]interface{})
	if !ok || index >= len(widgets) {
		return nil, false
	}
	widget, ok := widgets[index].(map[string]interface{})
	return widget, ok
}

func dashboardLayoutsOverlap(a, b dashboardWidgetLayout) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X && a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}

func invalidDefinition(description string) error {
	return &InvalidDashboardError{Description: description}
}
