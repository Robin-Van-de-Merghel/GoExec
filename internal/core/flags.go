package core

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

func SetupCLI() *cobra.Command {
	var listModules bool
	var listTags bool
	var moduleName string
	var showOptions bool

	rootCmd := &cobra.Command{
		Use:   "goexec [TAGS]...",
		Short: "GoExec network scanner",
		Long:  "GoExec - A network scanner written in go.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract tags from positional arguments
			tags := extractTags(args)

			// Show help if nothing provided
			if len(tags) == 0 && !listTags && !listModules && moduleName == "" {
				cmd.Help()
				return nil
			}

			// If module specified, validate tags
			if moduleName != "" && len(tags) == 0 {
				return fmt.Errorf("module specified but no tags provided")
			}

			// List all tags
			if listTags {
				fmt.Println("Available tags:")
				allTags := collectAllTags()
				for _, t := range allTags {
					fmt.Printf("  - %s\n", t)
				}
				return nil
			}

			// Filter by tags if provided
			matching := FilterModulesByTags(tags)

			// List matching modules
			if listModules {
				if len(matching) == 0 {
					if len(tags) > 0 {
						fmt.Printf("No modules found matching tags: %v\n", tags)
					} else {
						fmt.Println("No modules found")
					}
					return nil
				}
				fmt.Println("Matching modules:")
				for _, m := range matching {
					fmt.Printf("  - %s (tags: %s)\n", m.Metadata.UniqueName, strings.Join(m.Metadata.Labels, ", "))
				}
				return nil
			}

			// If module name specified, run it
			if moduleName != "" {
				// Show options if requested
				if showOptions {
					entry := findModuleByName(moduleName, matching)
					if entry == nil {
						return fmt.Errorf("module '%s' not found", moduleName)
					}
					fmt.Println(GenerateHelpMessage(*entry))
					return nil
				}

				entry := findModuleByName(moduleName, matching)
				if entry == nil {
					return fmt.Errorf("module '%s' not found", moduleName)
				}
				return RunModuleWithLogging(*entry, cmd)
			}

			// Tag provided but no action
			if len(tags) > 0 && !listTags && !listModules {
				return fmt.Errorf("tag(s) %v provided but no action. Use: -L (list modules) or -M MODULE_NAME (run)", tags)
			}

			return nil
		},
	}

	// Custom help function for root
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("GoExec - A dynamic module runner for executing various tools and utilities")
		fmt.Println("\nUsage:")
		fmt.Println("  goexec [TAGS]... [flags]")
		fmt.Println("\nFlags:")
		cmd.Flags().PrintDefaults()
	})

	// Global flags
	rootCmd.Flags().BoolVarP(&listTags, "list-tags", "T", false, "List all available tags")
	rootCmd.Flags().BoolVarP(&listModules, "list-modules", "L", false, "List matching modules")
	rootCmd.Flags().StringVarP(&moduleName, "module", "M", "", "Module name to run")
	rootCmd.Flags().BoolVar(&showOptions, "options", false, "Show module options/help")

	// Check if a module is being requested and add its flags
	addModuleFlagsIfRequested(rootCmd)

	return rootCmd
}

// addModuleFlagsIfRequested checks os.Args for -M flag and adds module flags if found
func addModuleFlagsIfRequested(cmd *cobra.Command) {
	// Find -M and its value in os.Args
	var moduleName string
	var tags []string

	for i, arg := range os.Args[1:] {
		if arg == "-M" || arg == "--module" {
			if i+1 < len(os.Args)-1 {
				moduleName = os.Args[i+2]
			}
			break
		}
	}

	if moduleName != "" {
		// Extract tags
		tags = extractTags(os.Args[1:])

		// Find module
		matching := FilterModulesByTags(tags)
		entry := findModuleByName(moduleName, matching)

		if entry != nil {
			// Add flags from module input
			addFlagsFromStruct(cmd, getModuleInputType(*entry))
		}
	}
}

// extractTags extracts non-flag positional arguments as tags
func extractTags(args []string) []string {
	var tags []string
	skipNext := false

	for _, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Skip flags and their values
		if strings.HasPrefix(arg, "-") {
			// Check if this flag takes a value
			if arg == "-M" || arg == "--module" || arg == "-t" || arg == "--tags" {
				skipNext = true
			}
			continue
		}

		// This is a positional argument (tag)
		tags = append(tags, arg)
	}
	return tags
}

// collectAllTags returns all unique tags
func collectAllTags() []string {
	tagSet := map[string]struct{}{}
	for _, m := range AllModules {
		for _, t := range m.Metadata.Labels {
			tagSet[strings.ToLower(t)] = struct{}{}
		}
	}
	all := []string{}
	for t := range tagSet {
		all = append(all, t)
	}
	return all
}

// findModuleByName finds a module by its unique name (case-insensitive)
func findModuleByName(name string, modules []ModuleEntry) *ModuleEntry {
	nameLower := strings.ToLower(name)
	for i, m := range modules {
		if strings.ToLower(m.Metadata.UniqueName) == nameLower {
			return &modules[i]
		}
	}
	return nil
}

// getModuleInputType returns the Input field type of a module
func getModuleInputType(entry ModuleEntry) reflect.Type {
	mod := entry.Factory()
	t := reflect.Indirect(reflect.ValueOf(mod)).Type()
	inputField, ok := t.FieldByName("Input")
	if !ok {
		return nil
	}
	return inputField.Type
}

// flattenStructFields recursively flattens struct fields and returns a map of fieldName -> (type, tag)
func flattenStructFields(t reflect.Type, prefix string) map[string]reflect.StructField {
	fields := make(map[string]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldName := f.Name

		// If struct, recurse
		if f.Type.Kind() == reflect.Struct {
			nested := flattenStructFields(f.Type, "")
			for k, v := range nested {
				fields[k] = v
			}
			continue
		}

		fields[fieldName] = f
	}
	return fields
}

// Add flags from flattened struct fields
func addFlagsFromStruct(cmd *cobra.Command, t reflect.Type) {
	if t == nil {
		return
	}
	fields := flattenStructFields(t, "")
	for fieldName, field := range fields {
		help := field.Tag.Get("help")
		if help == "" {
			help = fieldName + " (" + field.Type.Name() + ")"
		}
		cmd.Flags().String(fieldName, "", help)
	}
}

// Recursively populate struct from CLI flags (flattened)
func setStructFromFlags(val reflect.Value, t reflect.Type, cmd *cobra.Command, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)
		fieldName := f.Name

		if f.Type.Kind() == reflect.Struct {
			setStructFromFlags(fv, f.Type, cmd, prefix)
			continue
		}

		flagName := fieldName
		if cmd.Flags().Changed(flagName) {
			valStr, _ := cmd.Flags().GetString(flagName)
			fv.SetString(valStr)
		}
	}
}

// ExecuteCLI runs the CLI (call from main.go)
func ExecuteCLI() {
	root := SetupCLI()
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
