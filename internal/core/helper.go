package core

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateHelpMessage(entry ModuleEntry) string {
	mod := entry.Factory()
	t := reflect.Indirect(reflect.ValueOf(mod)).Type()

	var sb strings.Builder

	fmt.Fprintf(&sb, "Module: %s\n", entry.Metadata.UniqueName)
	fmt.Fprintf(&sb, "Description: %s\n", entry.Metadata.PresentMessages)
	fmt.Fprintf(&sb, "Tags: %s\n", strings.Join(entry.Metadata.Labels, ", "))
	sb.WriteString("Arguments:\n")

	inputField, ok := t.FieldByName("Input")
	if !ok || inputField.Type.Kind() != reflect.Struct {
		sb.WriteString("  (no input fields found)\n")
		return sb.String()
	}

	formatStructFields(inputField.Type, &sb)
	return sb.String()
}

// Recursively formats struct fields, using only the last name part
func formatStructFields(t reflect.Type, sb *strings.Builder) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			formatStructFields(f.Type, sb)
			continue
		}
		fieldName := f.Name
		tagHelp := f.Tag.Get("help")
		fmt.Fprintf(sb, " - %s: %s (%s)\n", fieldName, tagHelp, f.Type.Name())
	}
}
