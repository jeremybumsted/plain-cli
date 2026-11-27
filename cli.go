package main

import (
	"reflect"

	"github.com/alecthomas/kong"
)

// AfterApply propagates root flags into the selected leaf command so shared
// flags behave consistently across command groups.
func (cli *CLI) AfterApply(ctx *kong.Context) error {
	selected := ctx.Selected()
	if selected == nil {
		return nil
	}

	applyRootFlags(selected.Target, cli)
	return nil
}

func applyRootFlags(target reflect.Value, cli *CLI) {
	if !target.IsValid() {
		return
	}
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return
		}
		target = target.Elem()
	}
	if target.Kind() != reflect.Struct {
		return
	}

	if cli.ConfigPath != "" {
		if field := target.FieldByName("ConfigPath"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
			field.SetString(cli.ConfigPath)
		}
	}

	if field := target.FieldByName("Format"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		switch {
		case cli.Quiet:
			field.SetString("quiet")
		case cli.JSON:
			field.SetString("json")
		}
	}
}
