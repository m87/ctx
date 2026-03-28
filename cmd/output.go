package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func printOutput(cmd *cobra.Command, data any, textRenderer func() string, shellRenderer func() string) error {
	format := strings.ToLower(strings.TrimSpace(OutputFormat))
	if format == "" {
		format = "text"
	}

	switch format {
	case "text":
		cmd.Println(textRenderer())
		return nil
	case "json":
		encoded, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(encoded))
		return nil
	case "yaml":
		encoded, err := yaml.Marshal(data)
		if err != nil {
			return err
		}
		cmd.Println(string(encoded))
		return nil
	case "shell":
		if shellRenderer != nil {
			cmd.Print(shellRenderer())
			return nil
		}
		cmd.Print(toShell(data))
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", OutputFormat)
	}
}

func toShell(data any) string {
	lines := make([]string, 0)
	flattenShell(reflect.ValueOf(data), "RESULT", &lines)
	if len(lines) == 0 {
		return "RESULT=\n"
	}
	return strings.Join(lines, "\n") + "\n"
}

func flattenShell(value reflect.Value, key string, lines *[]string) {
	if !value.IsValid() {
		*lines = append(*lines, fmt.Sprintf("%s=", key))
		return
	}

	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			*lines = append(*lines, fmt.Sprintf("%s=", key))
			return
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		typeInfo := value.Type()
		for index := 0; index < value.NumField(); index++ {
			if !typeInfo.Field(index).IsExported() {
				continue
			}
			name := strings.ToUpper(typeInfo.Field(index).Name)
			flattenShell(value.Field(index), key+"_"+name, lines)
		}
	case reflect.Map:
		for _, mapKey := range value.MapKeys() {
			name := fmt.Sprintf("%v", mapKey.Interface())
			name = strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
			flattenShell(value.MapIndex(mapKey), key+"_"+name, lines)
		}
	case reflect.Slice, reflect.Array:
		for index := 0; index < value.Len(); index++ {
			flattenShell(value.Index(index), fmt.Sprintf("%s_%d", key, index), lines)
		}
	case reflect.String:
		*lines = append(*lines, fmt.Sprintf("%s=%q", key, value.String()))
	case reflect.Bool:
		*lines = append(*lines, fmt.Sprintf("%s=%t", key, value.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*lines = append(*lines, fmt.Sprintf("%s=%d", key, value.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		*lines = append(*lines, fmt.Sprintf("%s=%d", key, value.Uint()))
	case reflect.Float32, reflect.Float64:
		*lines = append(*lines, fmt.Sprintf("%s=%v", key, value.Float()))
	default:
		*lines = append(*lines, fmt.Sprintf("%s=%q", key, fmt.Sprintf("%v", value.Interface())))
	}
}
