package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
)

// isSchedulingSubcommand checks if the given command is a valid scheduling subcommand
func isSchedulingSubcommand(cmd string) bool {
	validSubcommands := []string{
		"tolerations", "affinity", "nodeselector",
		"resources", "topology", "priority", "runtime",
	}
	for _, v := range validSubcommands {
		if cmd == v {
			return true
		}
	}
	return false
}

// isHelpFlag checks if the argument is a help flag
func isHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "-help"
}

// containsHelpFlag checks if any argument in the slice is a help flag
func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if isHelpFlag(arg) {
			return true
		}
	}
	return false
}

// preprocessArgs expands short flags with attached values like -oyaml to -o yaml
// This mimics kubectl behavior where flags can be written without space
func preprocessArgs(args []string) []string {
	// Short flags that take string values
	stringFlags := map[string]bool{
		"-o": true,
		"-n": true,
		"-l": true,
	}

	// Short boolean flags (for combining like -Ac)
	boolFlags := map[string]bool{
		"-A": true,
		"-c": true,
	}

	var result []string

	for _, arg := range args {
		// Skip if not a flag
		if !strings.HasPrefix(arg, "-") {
			result = append(result, arg)
			continue
		}

		// Skip long flags (they work as expected)
		if strings.HasPrefix(arg, "--") {
			result = append(result, arg)
			continue
		}

		// Check if it's a short string flag with attached value (e.g., -oyaml)
		expanded := false
		for flag := range stringFlags {
			if strings.HasPrefix(arg, flag) && len(arg) > len(flag) {
				// Check if the character after flag is not '=' (handled by flag package)
				if arg[len(flag)] != '=' {
					// Split: -oyaml -> -o yaml
					result = append(result, flag, arg[len(flag):])
					expanded = true
					break
				}
			}
		}

		if expanded {
			continue
		}

		// Check for combined boolean flags (e.g., -Ac -> -A -c)
		if len(arg) > 2 {
			allBool := true
			for i := 1; i < len(arg); i++ {
				shortFlag := "-" + string(arg[i])
				if !boolFlags[shortFlag] {
					allBool = false
					break
				}
			}
			if allBool {
				for i := 1; i < len(arg); i++ {
					result = append(result, "-"+string(arg[i]))
				}
				continue
			}
		}

		// No expansion needed
		result = append(result, arg)
	}

	return result
}

func main() {
	// Check for help with no arguments or just help flag
	if len(os.Args) < 2 || isHelpFlag(os.Args[1]) {
		printUsage()
		os.Exit(0)
	}

	// Parse command type
	cmdType := os.Args[1]

	// Handle completion command
	if cmdType == "completion" {
		handleCompletion(os.Args[2:])
		os.Exit(0)
	}
	var subCommand string
	var resourceType string
	var argsOffset int

	// Check if it's scheduling command
	if cmdType == "scheduling" {
		// Check for help: kubectl getinfo scheduling --help
		if len(os.Args) < 3 || isHelpFlag(os.Args[2]) {
			printSchedulingUsage("")
			os.Exit(0)
		}

		// Check if second argument is a subcommand
		if isSchedulingSubcommand(os.Args[2]) {
			subCommand = os.Args[2]
			// Check for help: kubectl getinfo scheduling <subcommand> --help
			if len(os.Args) < 4 || isHelpFlag(os.Args[3]) {
				printSchedulingUsage(subCommand)
				os.Exit(0)
			}
			// Check if any remaining args contain help
			if containsHelpFlag(os.Args[4:]) {
				printSchedulingUsage(subCommand)
				os.Exit(0)
			}
			resourceType = os.Args[3]
			argsOffset = 4
		} else {
			// Check if any remaining args contain help
			if containsHelpFlag(os.Args[3:]) {
				printSchedulingUsage("")
				os.Exit(0)
			}
			// No subcommand, second arg is resource type
			resourceType = os.Args[2]
			argsOffset = 3
		}
	} else {
		// Other commands (labels, annotations, owner)
		if cmdType != "labels" && cmdType != "annotations" && cmdType != "owner" {
			fmt.Fprintf(os.Stderr, "Error: command type must be 'labels', 'annotations', 'owner', 'scheduling', or 'completion', got '%s'\n", cmdType)
			printUsage()
			os.Exit(1)
		}

		// Check for help: kubectl getinfo <command> --help
		if len(os.Args) < 3 || isHelpFlag(os.Args[2]) {
			printCommandUsage(cmdType)
			os.Exit(0)
		}

		// Check if any remaining args contain help
		if containsHelpFlag(os.Args[3:]) {
			printCommandUsage(cmdType)
			os.Exit(0)
		}

		resourceType = os.Args[2]
		argsOffset = 3
	}

	// Parse flags
	var namespace string
	var allNamespaces bool
	var selector string
	var outputFormat string
	var colorOutput bool

	fs := flag.NewFlagSet("getinfo", flag.ExitOnError)
	fs.StringVar(&namespace, "n", "", "namespace")
	fs.StringVar(&namespace, "namespace", "", "namespace")
	fs.BoolVar(&allNamespaces, "A", false, "all-namespaces")
	fs.BoolVar(&allNamespaces, "all-namespaces", false, "all-namespaces")
	fs.StringVar(&selector, "l", "", "selector")
	fs.StringVar(&selector, "selector", "", "selector")
	fs.StringVar(&outputFormat, "o", "json", "output format (json, yaml, table)")
	fs.StringVar(&outputFormat, "output", "json", "output format (json, yaml, table)")
	fs.BoolVar(&colorOutput, "c", false, "colorize JSON output")
	fs.BoolVar(&colorOutput, "color", false, "colorize JSON output")

	// Parse remaining arguments (resource names and flags)
	args := os.Args[argsOffset:]
	// Preprocess args to expand short flags like -oyaml to -o yaml
	args = preprocessArgs(args)
	fs.Parse(args)

	// Get resource names (non-flag arguments after parsing)
	resourceNames := fs.Args()

	// Get kubeconfig
	config, err := getKubeconfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	// Get GVR (GroupVersionResource) for the resource type
	gvr, namespaced, err := getGVR(resourceType, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Determine namespace
	if allNamespaces {
		namespace = ""
	} else if namespace == "" && namespaced {
		// Try to get namespace from kubeconfig context
		namespace = getCurrentNamespace()
	}

	// Parse label selector
	var labelSelector labels.Selector
	if selector != "" {
		labelSelector, err = labels.Parse(selector)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing selector: %v\n", err)
			os.Exit(1)
		}
	}

	// Get resources
	items, err := getResources(dynamicClient, gvr, namespaced, namespace, resourceNames, labelSelector)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting resources: %v\n", err)
		os.Exit(1)
	}

	// Extract labels, annotations, or ownerReferences
	output := Output{Items: []OutputItem{}}
	for _, item := range items {
		outputItem := OutputItem{
			Name: item.GetName(),
		}

		if namespaced {
			outputItem.Namespace = item.GetNamespace()
		}

		switch cmdType {
		case "labels":
			labels := item.GetLabels()
			outputItem.Labels = &labels
		case "annotations":
			annotations := item.GetAnnotations()
			outputItem.Annotations = &annotations
		case "owner":
			ownerRefs := extractOwnerReferences(item)
			outputItem.OwnerReferences = ownerRefs
			// Don't fill labels and annotations when the command is owner
		case "scheduling":
			if subCommand == "" {
				// Show all scheduling info
				schedulingInfo := extractSchedulingInfo(item)
				outputItem.Scheduling = schedulingInfo
			} else {
				// Show only the specific subcommand field
				extractSchedulingSubcommand(item, &outputItem, subCommand)
			}
		}

		output.Items = append(output.Items, outputItem)
	}

	// Output in requested format
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		if colorOutput {
			coloredOutput := colorizeJSON(string(jsonOutput))
			fmt.Print(coloredOutput)
		} else {
			fmt.Println(string(jsonOutput))
		}
	case "yaml":
		yamlOutput, err := yaml.Marshal(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling YAML: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(yamlOutput))
	case "table":
		printTable(output, cmdType, subCommand, namespaced)
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported output format '%s'. Supported formats: json, yaml, table\n", outputFormat)
		os.Exit(1)
	}
}

func init() {
	// Add scheme for proper API discovery
	_ = scheme.AddToScheme(scheme.Scheme)
}
