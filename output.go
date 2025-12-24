package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

// colorizeJSON adds ANSI color codes to JSON output (similar to jq)
func colorizeJSON(jsonStr string) string {
	const (
		reset      = "\033[0m"
		keyColor   = "\033[1;34m" // bold blue for keys
		strColor   = "\033[32m"   // green for strings
		numColor   = "\033[33m"   // yellow for numbers
		boolColor  = "\033[1;33m" // bold yellow for booleans
		nullColor  = "\033[90m"   // gray for null
		punctColor = "\033[37m"   // white for punctuation
	)

	result := jsonStr

	// Colorize punctuation first ({, }, [, ])
	punctRegex := regexp.MustCompile(`([{}\[\]])`)
	result = punctRegex.ReplaceAllStringFunc(result, func(match string) string {
		return punctColor + match + reset
	})

	// Colorize keys (pattern: "key":)
	keyRegex := regexp.MustCompile(`"([^"]+)":`)
	result = keyRegex.ReplaceAllStringFunc(result, func(match string) string {
		return keyColor + match + reset
	})

	// Colorize strings (values in quotes that are not keys)
	// We need to avoid colorizing keys again, so we do this after
	strRegex := regexp.MustCompile(`:\s*"([^"]*)"`)
	result = strRegex.ReplaceAllStringFunc(result, func(match string) string {
		// Preserve the ":" and spaces, colorize only the string
		if strings.HasPrefix(match, ": ") {
			return ": " + strColor + `"` + strings.TrimPrefix(strings.TrimSuffix(match[2:], `"`), `"`) + `"` + reset
		} else if strings.HasPrefix(match, ":") {
			return ":" + strColor + match[1:] + reset
		}
		return match
	})

	// Colorize numbers (integers and decimals)
	numRegex := regexp.MustCompile(`:\s*(-?\d+\.?\d*)`)
	result = numRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) == 2 {
			return parts[0] + ":" + numColor + strings.TrimSpace(parts[1]) + reset
		}
		return match
	})

	// Colorize booleans
	boolRegex := regexp.MustCompile(`:\s*(true|false)`)
	result = boolRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) == 2 {
			return parts[0] + ":" + boolColor + strings.TrimSpace(parts[1]) + reset
		}
		return match
	})

	// Colorize null
	nullRegex := regexp.MustCompile(`:\s*(null)`)
	result = nullRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) == 2 {
			return parts[0] + ":" + nullColor + strings.TrimSpace(parts[1]) + reset
		}
		return match
	})

	return result
}

// printTable outputs the data in table format
func printTable(output Output, cmdType string, subCommand string, namespaced bool) {
	if len(output.Items) == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	if namespaced {
		fmt.Fprintf(w, "NAME\tNAMESPACE\t")
	} else {
		fmt.Fprintf(w, "NAME\t")
	}

	// Determine column header based on cmdType
	if cmdType == "labels" {
		fmt.Fprintf(w, "LABELS\n")
	} else if cmdType == "annotations" {
		fmt.Fprintf(w, "ANNOTATIONS\n")
	} else if cmdType == "owner" {
		if namespaced {
			fmt.Fprintf(w, "OWNER NAMESPACE\tOWNER KIND\tOWNER NAME\n")
		} else {
			fmt.Fprintf(w, "OWNER KIND\tOWNER NAME\n")
		}
	} else if cmdType == "scheduling" {
		if subCommand == "" {
			// Show summary of all fields
			fmt.Fprintf(w, "NODESELECTOR\tAFFINITY\tTOLERATIONS\tRESOURCES\n")
		} else {
			// Show only the specific field
			switch subCommand {
			case "tolerations":
				fmt.Fprintf(w, "TOLERATIONS\n")
			case "affinity":
				fmt.Fprintf(w, "AFFINITY\n")
			case "nodeselector":
				fmt.Fprintf(w, "NODESELECTOR\n")
			case "resources":
				fmt.Fprintf(w, "RESOURCES\n")
			case "topology":
				fmt.Fprintf(w, "TOPOLOGY SPREAD CONSTRAINTS\n")
			case "priority":
				fmt.Fprintf(w, "PRIORITY\n")
			case "runtime":
				fmt.Fprintf(w, "RUNTIME\n")
			}
		}
	}

	// Print separator
	if namespaced {
		fmt.Fprintf(w, "----\t---------\t")
	} else {
		fmt.Fprintf(w, "----\t")
	}
	if cmdType == "owner" {
		if namespaced {
			fmt.Fprintf(w, "---------------\t----------\t----------\n")
		} else {
			fmt.Fprintf(w, "----------\t----------\n")
		}
	} else if cmdType == "scheduling" {
		if subCommand == "" {
			fmt.Fprintf(w, "-----------\t--------\t-----------\t---------\n")
		} else {
			fmt.Fprintf(w, "--------\n")
		}
	} else {
		fmt.Fprintf(w, "--------\n")
	}

	// Print items
	for _, item := range output.Items {
		if cmdType == "owner" {
			// Handle ownerReferences
			if len(item.OwnerReferences) == 0 {
				if namespaced {
					fmt.Fprintf(w, "%s\t%s\t<none>\t<none>\t<none>\n", item.Name, item.Namespace)
				} else {
					fmt.Fprintf(w, "%s\t<none>\t<none>\n", item.Name)
				}
			} else {
				for i, ownerRef := range item.OwnerReferences {
					if i == 0 {
						// First owner reference - show resource name
						if namespaced {
							fmt.Fprintf(w, "%s\t%s\t", item.Name, item.Namespace)
						} else {
							fmt.Fprintf(w, "%s\t", item.Name)
						}
					} else {
						// Additional owner references - show empty name/namespace
						if namespaced {
							fmt.Fprintf(w, "\t\t")
						} else {
							fmt.Fprintf(w, "\t")
						}
					}

					if namespaced {
						ownerNamespace := ownerRef.Namespace
						if ownerNamespace == "" {
							ownerNamespace = "<none>"
						}
						fmt.Fprintf(w, "%s\t%s\t%s\n", ownerNamespace, ownerRef.Kind, ownerRef.Name)
					} else {
						fmt.Fprintf(w, "%s\t%s\n", ownerRef.Kind, ownerRef.Name)
					}
				}
			}
		} else {
			// Handle labels or annotations
			if namespaced {
				fmt.Fprintf(w, "%s\t%s\t", item.Name, item.Namespace)
			} else {
				fmt.Fprintf(w, "%s\t", item.Name)
			}

			// Format labels or annotations as key=value pairs
			var pairs []string
			if cmdType == "labels" && item.Labels != nil {
				// Sort keys for consistent output
				keys := make([]string, 0, len(*item.Labels))
				for k := range *item.Labels {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, (*item.Labels)[k]))
				}
			} else if cmdType == "annotations" && item.Annotations != nil {
				// Sort keys for consistent output
				keys := make([]string, 0, len(*item.Annotations))
				for k := range *item.Annotations {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, (*item.Annotations)[k]))
				}
			}

			if len(pairs) > 0 {
				fmt.Fprintf(w, "%s\n", strings.Join(pairs, ","))
			} else {
				fmt.Fprintf(w, "<none>\n")
			}
		}

		if cmdType == "scheduling" {
			// Handle scheduling
			if namespaced {
				fmt.Fprintf(w, "%s\t%s\t", item.Name, item.Namespace)
			} else {
				fmt.Fprintf(w, "%s\t", item.Name)
			}

			if subCommand == "" {
				// Show summary
				var nodeSelectorStr, affinityStr, tolerationsStr, resourcesStr string

				if item.Scheduling != nil {
					if len(item.Scheduling.NodeSelector) > 0 {
						var pairs []string
						for k, v := range item.Scheduling.NodeSelector {
							pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
						}
						sort.Strings(pairs)
						nodeSelectorStr = strings.Join(pairs, ",")
					} else {
						nodeSelectorStr = "<none>"
					}

					if len(item.Scheduling.Affinity) > 0 {
						affinityStr = "present"
					} else {
						affinityStr = "<none>"
					}

					if len(item.Scheduling.Tolerations) > 0 {
						tolerationsStr = fmt.Sprintf("%d item(s)", len(item.Scheduling.Tolerations))
					} else {
						tolerationsStr = "<none>"
					}

					if item.Scheduling.ResourceRequests != nil || item.Scheduling.ResourceLimits != nil {
						resourcesStr = "present"
					} else {
						resourcesStr = "<none>"
					}
				} else {
					nodeSelectorStr = "<none>"
					affinityStr = "<none>"
					tolerationsStr = "<none>"
					resourcesStr = "<none>"
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", nodeSelectorStr, affinityStr, tolerationsStr, resourcesStr)
			} else {
				// Show specific field
				var valueStr string
				switch subCommand {
				case "tolerations":
					if len(item.Tolerations) > 0 {
						valueStr = fmt.Sprintf("%d toleration(s)", len(item.Tolerations))
					} else {
						valueStr = "<none>"
					}
				case "affinity":
					if len(item.Affinity) > 0 {
						valueStr = "present"
					} else {
						valueStr = "<none>"
					}
				case "nodeselector":
					if len(item.NodeSelector) > 0 {
						var pairs []string
						for k, v := range item.NodeSelector {
							pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
						}
						sort.Strings(pairs)
						valueStr = strings.Join(pairs, ",")
					} else {
						valueStr = "<none>"
					}
				case "resources":
					if len(item.Resources) > 0 {
						valueStr = fmt.Sprintf("%d container(s)", len(item.Resources))
					} else {
						valueStr = "<none>"
					}
				case "topology":
					if len(item.TopologySpreadConstraints) > 0 {
						valueStr = fmt.Sprintf("%d constraint(s)", len(item.TopologySpreadConstraints))
					} else {
						valueStr = "<none>"
					}
				case "priority":
					if len(item.Priority) > 0 {
						valueStr = "present"
					} else {
						valueStr = "<none>"
					}
				case "runtime":
					if len(item.Runtime) > 0 {
						valueStr = "present"
					} else {
						valueStr = "<none>"
					}
				}
				fmt.Fprintf(w, "%s\n", valueStr)
			}
		}
	}
}

