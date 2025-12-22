package main

import (
	"fmt"
	"os"
)

// printUsage prints the general usage information
func printUsage() {
	fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo <command> [subcommand] <resource-type> [resource-name...] [flags]

Commands:
  labels       List labels of resources
  annotations  List annotations of resources
  owner        List ownerReferences of resources
  scheduling   List scheduling-related fields (nodeSelector, affinity, tolerations, etc.)
  completion   Generate shell completion scripts (bash, zsh, fish)

Scheduling Subcommands (optional):
  tolerations       List only tolerations
  affinity          List only affinity rules
  nodeselector      List only nodeSelector
  resources         List only resource requests/limits
  topology          List only topologySpreadConstraints
  priority          List only priority-related fields
  runtime           List only runtime-related fields (runtimeClassName, hostNetwork, etc.)

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help

Examples:
  kubectl getinfo labels pods pod1 pod2
  kubectl getinfo annotations nodes -l env=prod
  kubectl getinfo owner pods
  kubectl getinfo scheduling pods
  kubectl getinfo scheduling tolerations pods
  kubectl getinfo scheduling affinity pods -n kube-system
  kubectl getinfo labels deployments -n kube-system -o yaml
  kubectl getinfo labels pods -o table
  kubectl getinfo labels pods -o json -c

Use "kubectl getinfo <command> --help" for more information about a command.
`)
}

// printCommandUsage prints usage information for a specific command
func printCommandUsage(cmdType string) {
	switch cmdType {
	case "labels":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo labels <resource-type> [resource-name...] [flags]

List labels of Kubernetes resources.

Examples:
  kubectl getinfo labels pods                          # List labels of all pods in current namespace
  kubectl getinfo labels pods pod1 pod2                # List labels of specific pods
  kubectl getinfo labels pods -A                       # List labels of all pods in all namespaces
  kubectl getinfo labels pods -n kube-system           # List labels of pods in kube-system namespace
  kubectl getinfo labels deployments -l app=nginx     # List labels of deployments with label app=nginx
  kubectl getinfo labels pods -o yaml                  # Output in YAML format
  kubectl getinfo labels pods -o table                 # Output in table format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "annotations":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo annotations <resource-type> [resource-name...] [flags]

List annotations of Kubernetes resources.

Examples:
  kubectl getinfo annotations pods                     # List annotations of all pods in current namespace
  kubectl getinfo annotations pods pod1 pod2          # List annotations of specific pods
  kubectl getinfo annotations pods -A                  # List annotations of all pods in all namespaces
  kubectl getinfo annotations services -n default     # List annotations of services in default namespace
  kubectl getinfo annotations deployments -o yaml     # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "owner":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo owner <resource-type> [resource-name...] [flags]

List ownerReferences of Kubernetes resources. Shows the parent resources that own/control the queried resources.

Examples:
  kubectl getinfo owner pods                           # List owner references of all pods
  kubectl getinfo owner pods pod1                      # List owner references of specific pod
  kubectl getinfo owner pods -A                        # List owner references of all pods in all namespaces
  kubectl getinfo owner replicasets -n kube-system    # List owner references of replicasets
  kubectl getinfo owner pods -o yaml                   # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	}
}

// printSchedulingUsage prints usage information for the scheduling command
func printSchedulingUsage(subCommand string) {
	if subCommand == "" {
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling [subcommand] <resource-type> [resource-name...] [flags]

List scheduling-related fields of Kubernetes resources (pods, deployments, statefulsets, etc.).

When no subcommand is specified, shows all scheduling information including:
  - nodeSelector
  - nodeName
  - affinity (nodeAffinity, podAffinity, podAntiAffinity)
  - tolerations
  - topologySpreadConstraints
  - resource requests and limits
  - schedulerName, priorityClassName, priority
  - runtimeClassName, hostNetwork, hostPID, hostIPC

Subcommands:
  tolerations       List only tolerations
  affinity          List only affinity rules
  nodeselector      List only nodeSelector
  resources         List only resource requests/limits
  topology          List only topologySpreadConstraints
  priority          List only priority-related fields (priorityClassName, priority, preemptionPolicy)
  runtime           List only runtime-related fields (runtimeClassName, hostNetwork, hostPID, hostIPC)

Examples:
  kubectl getinfo scheduling pods                      # List all scheduling info of pods
  kubectl getinfo scheduling pods -A                   # List scheduling info of all pods in all namespaces
  kubectl getinfo scheduling deployments -n prod      # List scheduling info of deployments in prod namespace
  kubectl getinfo scheduling tolerations pods          # List only tolerations
  kubectl getinfo scheduling affinity pods             # List only affinity rules
  kubectl getinfo scheduling resources pods            # List only resource requests/limits
  kubectl getinfo scheduling pods -o yaml              # Output in YAML format
  kubectl getinfo scheduling pods -o table             # Output in table format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help

Use "kubectl getinfo scheduling <subcommand> --help" for more information about a subcommand.
`)
		return
	}

	// Subcommand-specific help
	switch subCommand {
	case "tolerations":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling tolerations <resource-type> [resource-name...] [flags]

List tolerations of Kubernetes resources. Tolerations allow pods to be scheduled on nodes with matching taints.

Examples:
  kubectl getinfo scheduling tolerations pods                    # List tolerations of all pods
  kubectl getinfo scheduling tolerations pods -A                 # List tolerations of all pods in all namespaces
  kubectl getinfo scheduling tolerations deployments -n prod    # List tolerations of deployments in prod
  kubectl getinfo scheduling tolerations pods -o yaml            # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "affinity":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling affinity <resource-type> [resource-name...] [flags]

List affinity rules of Kubernetes resources. Includes nodeAffinity, podAffinity, and podAntiAffinity.

Examples:
  kubectl getinfo scheduling affinity pods                       # List affinity rules of all pods
  kubectl getinfo scheduling affinity pods -A                    # List affinity of all pods in all namespaces
  kubectl getinfo scheduling affinity deployments -n prod       # List affinity of deployments in prod
  kubectl getinfo scheduling affinity pods -o yaml               # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "nodeselector":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling nodeselector <resource-type> [resource-name...] [flags]

List nodeSelector of Kubernetes resources. NodeSelector is the simplest way to constrain pods to nodes with specific labels.

Examples:
  kubectl getinfo scheduling nodeselector pods                   # List nodeSelector of all pods
  kubectl getinfo scheduling nodeselector pods -A                # List nodeSelector of all pods in all namespaces
  kubectl getinfo scheduling nodeselector deployments -n prod   # List nodeSelector of deployments in prod
  kubectl getinfo scheduling nodeselector pods -o yaml           # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "resources":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling resources <resource-type> [resource-name...] [flags]

List resource requests and limits of Kubernetes resources. Shows CPU and memory requests/limits for containers.

Examples:
  kubectl getinfo scheduling resources pods                      # List resources of all pods
  kubectl getinfo scheduling resources pods -A                   # List resources of all pods in all namespaces
  kubectl getinfo scheduling resources deployments -n prod      # List resources of deployments in prod
  kubectl getinfo scheduling resources pods -o yaml              # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "topology":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling topology <resource-type> [resource-name...] [flags]

List topologySpreadConstraints of Kubernetes resources. These constraints control how pods are spread across topology domains.

Examples:
  kubectl getinfo scheduling topology pods                       # List topology constraints of all pods
  kubectl getinfo scheduling topology pods -A                    # List topology constraints of all pods
  kubectl getinfo scheduling topology deployments -n prod       # List topology constraints of deployments
  kubectl getinfo scheduling topology pods -o yaml               # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "priority":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling priority <resource-type> [resource-name...] [flags]

List priority-related fields of Kubernetes resources. Includes priorityClassName, priority value, and preemptionPolicy.

Examples:
  kubectl getinfo scheduling priority pods                       # List priority info of all pods
  kubectl getinfo scheduling priority pods -A                    # List priority info of all pods
  kubectl getinfo scheduling priority deployments -n prod       # List priority info of deployments
  kubectl getinfo scheduling priority pods -o yaml               # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	case "runtime":
		fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo scheduling runtime <resource-type> [resource-name...] [flags]

List runtime-related fields of Kubernetes resources. Includes runtimeClassName, hostNetwork, hostPID, and hostIPC.

Examples:
  kubectl getinfo scheduling runtime pods                        # List runtime info of all pods
  kubectl getinfo scheduling runtime pods -A                     # List runtime info of all pods
  kubectl getinfo scheduling runtime deployments -n prod        # List runtime info of deployments
  kubectl getinfo scheduling runtime pods -o yaml                # Output in YAML format

Flags:
  -n, --namespace <namespace>      Specify namespace
  -A, --all-namespaces             All namespaces
  -l, --selector <selector>        Label selector (e.g., -l app=nginx)
  -o, --output <format>            Output format (json, yaml, table). Default: json
  -c, --color                      Colorize JSON output
  -h, --help                       Show help
`)
	}
}

