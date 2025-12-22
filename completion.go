package main

import (
	"fmt"
	"os"
)

// printCompletionUsage prints usage for the completion command
func printCompletionUsage() {
	fmt.Fprintf(os.Stdout, `Usage: kubectl getinfo completion <shell>

Generate shell completion scripts for kubectl-getinfo.

Supported shells:
  bash    Generate bash completion script
  zsh     Generate zsh completion script
  fish    Generate fish completion script

Examples:
  # Bash (add to ~/.bashrc)
  source <(kubectl-getinfo completion bash)

  # Zsh (add to ~/.zshrc)
  source <(kubectl-getinfo completion zsh)

  # Fish (add to ~/.config/fish/config.fish)
  kubectl-getinfo completion fish | source

  # Or save to a file:
  kubectl-getinfo completion bash > /etc/bash_completion.d/kubectl-getinfo
`)
}

// generateBashCompletion generates bash completion script
func generateBashCompletion() {
	fmt.Print(`# bash completion for kubectl-getinfo

_kubectl_getinfo_completions() {
    local cur prev words cword
    _init_completion || return

    local commands="labels annotations owner scheduling completion"
    local scheduling_subcommands="tolerations affinity nodeselector resources topology priority runtime"
    local resource_types="pods po deployments deploy services svc nodes no configmaps cm secrets sec statefulsets sts daemonsets ds replicasets rs ingresses ing jobs cronjobs cj persistentvolumes pv persistentvolumeclaims pvc namespaces ns serviceaccounts sa endpoints ep events ev networkpolicies netpol"
    local output_formats="json yaml table"

    # Count non-flag arguments
    local args=()
    local i
    for ((i=1; i < cword; i++)); do
        if [[ "${words[i]}" != -* ]]; then
            args+=("${words[i]}")
        fi
    done

    # First argument: command
    if [[ ${#args[@]} -eq 0 ]]; then
        if [[ "$cur" == -* ]]; then
            COMPREPLY=($(compgen -W "-h --help" -- "$cur"))
        else
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
        fi
        return
    fi

    local cmd="${args[0]}"

    # Handle completion command
    if [[ "$cmd" == "completion" ]]; then
        if [[ ${#args[@]} -eq 1 ]]; then
            COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
        fi
        return
    fi

    # Handle scheduling command with subcommands
    if [[ "$cmd" == "scheduling" ]]; then
        if [[ ${#args[@]} -eq 1 ]]; then
            # Could be subcommand or resource type
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "-n --namespace -A --all-namespaces -l --selector -o --output -c --color -h --help" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "$scheduling_subcommands $resource_types" -- "$cur"))
            fi
            return
        fi

        # Check if second arg is a subcommand
        local second_arg="${args[1]}"
        local is_subcommand=0
        for sub in $scheduling_subcommands; do
            if [[ "$second_arg" == "$sub" ]]; then
                is_subcommand=1
                break
            fi
        done

        if [[ $is_subcommand -eq 1 && ${#args[@]} -eq 2 ]]; then
            # After subcommand, suggest resource types
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "-n --namespace -A --all-namespaces -l --selector -o --output -c --color -h --help" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "$resource_types" -- "$cur"))
            fi
            return
        fi
    fi

    # For other commands (labels, annotations, owner) or after resource type
    if [[ ${#args[@]} -eq 1 ]]; then
        # After command, suggest resource types
        if [[ "$cur" == -* ]]; then
            COMPREPLY=($(compgen -W "-n --namespace -A --all-namespaces -l --selector -o --output -c --color -h --help" -- "$cur"))
        else
            COMPREPLY=($(compgen -W "$resource_types" -- "$cur"))
        fi
        return
    fi

    # Handle flags
    if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "-n --namespace -A --all-namespaces -l --selector -o --output -c --color -h --help" -- "$cur"))
        return
    fi

    # Handle flag values
    case "$prev" in
        -o|--output)
            COMPREPLY=($(compgen -W "$output_formats" -- "$cur"))
            return
            ;;
        -n|--namespace)
            # Try to get namespaces from kubectl
            local namespaces
            if namespaces=$(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null); then
                COMPREPLY=($(compgen -W "$namespaces" -- "$cur"))
            fi
            return
            ;;
    esac
}

complete -F _kubectl_getinfo_completions kubectl-getinfo

# Also register for "kubectl getinfo" if using as plugin
if [[ -n "$BASH_VERSION" ]]; then
    # For kubectl plugin usage, we rely on kubectl's plugin completion
    :
fi
`)
}

// generateZshCompletion generates zsh completion script
func generateZshCompletion() {
	fmt.Print(`#compdef kubectl-getinfo

# zsh completion for kubectl-getinfo

_kubectl_getinfo() {
    local -a commands scheduling_subcommands resource_types output_formats
    
    commands=(
        'labels:List labels of resources'
        'annotations:List annotations of resources'
        'owner:List ownerReferences of resources'
        'scheduling:List scheduling-related fields'
        'completion:Generate shell completion scripts'
    )
    
    scheduling_subcommands=(
        'tolerations:List only tolerations'
        'affinity:List only affinity rules'
        'nodeselector:List only nodeSelector'
        'resources:List only resource requests/limits'
        'topology:List only topologySpreadConstraints'
        'priority:List only priority-related fields'
        'runtime:List only runtime-related fields'
    )
    
    resource_types=(
        'pods' 'po'
        'deployments' 'deploy'
        'services' 'svc'
        'nodes' 'no'
        'configmaps' 'cm'
        'secrets' 'sec'
        'statefulsets' 'sts'
        'daemonsets' 'ds'
        'replicasets' 'rs'
        'ingresses' 'ing'
        'jobs'
        'cronjobs' 'cj'
        'persistentvolumes' 'pv'
        'persistentvolumeclaims' 'pvc'
        'namespaces' 'ns'
        'serviceaccounts' 'sa'
        'endpoints' 'ep'
        'events' 'ev'
        'networkpolicies' 'netpol'
    )
    
    output_formats=('json' 'yaml' 'table')

    _arguments -C \
        '1: :->command' \
        '2: :->second' \
        '3: :->third' \
        '*: :->args' \
        && return

    case $state in
        command)
            _describe -t commands 'kubectl-getinfo commands' commands
            ;;
        second)
            case $words[2] in
                completion)
                    _values 'shell' bash zsh fish
                    ;;
                scheduling)
                    _alternative \
                        'subcommands:scheduling subcommand:_describe -t scheduling-subcommands "scheduling subcommands" scheduling_subcommands' \
                        'resources:resource type:compadd -a resource_types'
                    ;;
                labels|annotations|owner)
                    _describe -t resources 'resource types' resource_types
                    ;;
            esac
            ;;
        third)
            if [[ $words[2] == "scheduling" ]]; then
                # Check if second word is a subcommand
                case $words[3] in
                    tolerations|affinity|nodeselector|resources|topology|priority|runtime)
                        _describe -t resources 'resource types' resource_types
                        ;;
                    *)
                        _arguments \
                            '-n[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                            '--namespace[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                            '-A[All namespaces]' \
                            '--all-namespaces[All namespaces]' \
                            '-l[Label selector]:selector:' \
                            '--selector[Label selector]:selector:' \
                            '-o[Output format]:format:(json yaml table)' \
                            '--output[Output format]:format:(json yaml table)' \
                            '-c[Colorize output]' \
                            '--color[Colorize output]' \
                            '*:resource name:'
                        ;;
                esac
            else
                _arguments \
                    '-n[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                    '--namespace[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                    '-A[All namespaces]' \
                    '--all-namespaces[All namespaces]' \
                    '-l[Label selector]:selector:' \
                    '--selector[Label selector]:selector:' \
                    '-o[Output format]:format:(json yaml table)' \
                    '--output[Output format]:format:(json yaml table)' \
                    '-c[Colorize output]' \
                    '--color[Colorize output]' \
                    '*:resource name:'
            fi
            ;;
        args)
            _arguments \
                '-n[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                '--namespace[Namespace]:namespace:_kubectl_getinfo_namespaces' \
                '-A[All namespaces]' \
                '--all-namespaces[All namespaces]' \
                '-l[Label selector]:selector:' \
                '--selector[Label selector]:selector:' \
                '-o[Output format]:format:(json yaml table)' \
                '--output[Output format]:format:(json yaml table)' \
                '-c[Colorize output]' \
                '--color[Colorize output]' \
                '*:resource name:'
            ;;
    esac
}

_kubectl_getinfo_namespaces() {
    local namespaces
    if namespaces=($(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)); then
        _describe -t namespaces 'namespaces' namespaces
    fi
}

compdef _kubectl_getinfo kubectl-getinfo
`)
}

// generateFishCompletion generates fish completion script
func generateFishCompletion() {
	fmt.Print(`# fish completion for kubectl-getinfo

# Disable file completion by default
complete -c kubectl-getinfo -f

# Commands
complete -c kubectl-getinfo -n "__fish_use_subcommand" -a "labels" -d "List labels of resources"
complete -c kubectl-getinfo -n "__fish_use_subcommand" -a "annotations" -d "List annotations of resources"
complete -c kubectl-getinfo -n "__fish_use_subcommand" -a "owner" -d "List ownerReferences of resources"
complete -c kubectl-getinfo -n "__fish_use_subcommand" -a "scheduling" -d "List scheduling-related fields"
complete -c kubectl-getinfo -n "__fish_use_subcommand" -a "completion" -d "Generate shell completion scripts"

# Completion subcommand
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"

# Scheduling subcommands
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "tolerations" -d "List only tolerations"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "affinity" -d "List only affinity rules"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "nodeselector" -d "List only nodeSelector"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "resources" -d "List only resource requests/limits"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "topology" -d "List only topologySpreadConstraints"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "priority" -d "List only priority-related fields"
complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling; and not __fish_seen_subcommand_from tolerations affinity nodeselector resources topology priority runtime" -a "runtime" -d "List only runtime-related fields"

# Resource types (for all commands)
set -l resource_types pods po deployments deploy services svc nodes no configmaps cm secrets sec statefulsets sts daemonsets ds replicasets rs ingresses ing jobs cronjobs cj persistentvolumes pv persistentvolumeclaims pvc namespaces ns serviceaccounts sa endpoints ep events ev networkpolicies netpol

for cmd in labels annotations owner
    complete -c kubectl-getinfo -n "__fish_seen_subcommand_from $cmd" -a "$resource_types"
end

complete -c kubectl-getinfo -n "__fish_seen_subcommand_from scheduling" -a "$resource_types"

# Flags (for all commands except completion)
complete -c kubectl-getinfo -n "not __fish_seen_subcommand_from completion" -s n -l namespace -d "Specify namespace" -x -a "(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | string split ' ')"
complete -c kubectl-getinfo -n "not __fish_seen_subcommand_from completion" -s A -l all-namespaces -d "All namespaces"
complete -c kubectl-getinfo -n "not __fish_seen_subcommand_from completion" -s l -l selector -d "Label selector"
complete -c kubectl-getinfo -n "not __fish_seen_subcommand_from completion" -s o -l output -d "Output format" -x -a "json yaml table"
complete -c kubectl-getinfo -n "not __fish_seen_subcommand_from completion" -s c -l color -d "Colorize JSON output"
complete -c kubectl-getinfo -s h -l help -d "Show help"
`)
}

// handleCompletion handles the completion command
func handleCompletion(args []string) {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printCompletionUsage()
		os.Exit(0)
	}

	shell := args[0]
	switch shell {
	case "bash":
		generateBashCompletion()
	case "zsh":
		generateZshCompletion()
	case "fish":
		generateFishCompletion()
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported shell '%s'. Supported shells: bash, zsh, fish\n", shell)
		os.Exit(1)
	}
}
