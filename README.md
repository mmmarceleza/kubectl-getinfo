# kubectl-getinfo

A `kubectl` plugin to list labels and annotations from Kubernetes objects in multiple formats (JSON, YAML, table).

## Installation

1. Build the plugin:
```bash
go build -o kubectl-getinfo
```

2. Put the binary on your PATH:
```bash
sudo mv kubectl-getinfo /usr/local/bin/
```

Or add it to your local PATH:
```bash
export PATH=$PATH:$(pwd)
```

## Shell Aliases (Optional)

For faster command execution, you can use the provided shell aliases file (`.getinfo_aliases`):

```bash
# Add to your .bashrc or .zshrc
[ -f /path/to/.getinfo_aliases ] && source /path/to/.getinfo_aliases
```

This provides 735+ aliases like:

| Alias | Command |
|-------|---------|
| `kgilp` | `kubectl getinfo labels pods` |
| `kgiap` | `kubectl getinfo annotations pods` |
| `kgiop` | `kubectl getinfo owner pods` |
| `kgisp` | `kubectl getinfo scheduling pods` |
| `kgistp` | `kubectl getinfo scheduling tolerations pods` |
| `kgilpA` | `kubectl getinfo labels pods -A` |
| `kgilpoyaml` | `kubectl getinfo labels pods -o yaml` |
| `kgilpAojson` | `kubectl getinfo labels pods -A -o json` |

**Alias pattern:**
- `kgi` = kubectl getinfo
- `l/a/o/s` = labels/annotations/owner/scheduling
- `t/af/ns/r` = tolerations/affinity/nodeselector/resources (scheduling subcommands)
- `p/d/svc/no/sts/ds` = pods/deploy/services/nodes/statefulsets/daemonsets
- `A/n/l` = -A (all namespaces) / -n (namespace) / -l (selector)
- `oyaml/ojson/otable` = output format

## Short Names Support

The plugin supports Kubernetes short names (just like `kubectl`):

```bash
kubectl getinfo labels po      # pods
kubectl getinfo labels no      # nodes
kubectl getinfo labels deploy  # deployments
kubectl getinfo labels svc     # services
kubectl getinfo labels cm      # configmaps
kubectl getinfo labels sec     # secrets
kubectl getinfo labels sts     # statefulsets
kubectl getinfo labels ds      # daemonsets
kubectl getinfo labels rs      # replicasets
kubectl getinfo labels ing     # ingresses
kubectl getinfo labels pv      # persistentvolumes
kubectl getinfo labels pvc     # persistentvolumeclaims
```

All short names are resolved dynamically via Kubernetes API Discovery, including CRDs with custom short names.

## Usage

### General Syntax

```bash
kubectl getinfo <type> <resource-type> [resource-name...] [flags]
```

Where:
- `<type>` can be `labels`, `annotations`, `owner`, or `scheduling`
- `[subcommand]` is optional and only used with `scheduling` (tolerations, affinity, nodeselector, resources, topology, priority, runtime)
- `<resource-type>` is the resource type (pods, nodes, deployments, etc.)
- `[resource-name...]` are optional names of specific resources
- `[flags]` are optional flags

**Note:** The plugin supports all Kubernetes resource types, including CRDs (Custom Resource Definitions). If the resource is not present in the internal map, the plugin uses Kubernetes discovery API to find it automatically.

### Supported Flags

- `-n, --namespace <namespace>` - Specify namespace
- `-A, --all-namespaces` - All namespaces
- `-l, --selector <selector>` - Filter by label selector (e.g., `-l app=nginx`)
- `-o, --output <format>` - Output format: `json` (default), `yaml`, or `table`
- `-c, --color` - Colorize JSON output (JSON format only)
- `-h, --help` - Show help (context-aware)

### Examples

#### Labels

```bash
# Labels for specific pods
kubectl getinfo labels pods pod1 pod2

# Labels for all pods in the current namespace
kubectl getinfo labels pods

# Labels for all pods in a specific namespace
kubectl getinfo labels pods -n kube-system

# Labels for pods filtered by a label selector
kubectl getinfo labels pods -l app=nginx

# Labels for nodes
kubectl getinfo labels nodes node1 node2

# Labels for all nodes
kubectl getinfo labels nodes

# Labels for nodes filtered by a label selector
kubectl getinfo labels nodes -l node-role.kubernetes.io/worker=
```

#### Annotations

```bash
# Annotations for specific pods
kubectl getinfo annotations pods pod1 pod2

# Annotations for all pods
kubectl getinfo annotations pods

# Annotations for deployments
kubectl getinfo annotations deployments -n default

# Annotations for nodes
kubectl getinfo annotations nodes

# OwnerReferences for pods
kubectl getinfo owner pods
kubectl getinfo owner pods -n kube-system

# Scheduling - all scheduling-related fields
kubectl getinfo scheduling pods
kubectl getinfo scheduling pods -n kube-system

# Scheduling - tolerations only
kubectl getinfo scheduling tolerations pods

# Scheduling - affinity only
kubectl getinfo scheduling affinity pods

# Scheduling - nodeSelector only
kubectl getinfo scheduling nodeselector pods

# Scheduling - resources only (requests/limits)
kubectl getinfo scheduling resources pods

# Labels from CRDs (Custom Resource Definitions)
kubectl getinfo labels servicemonitor
kubectl getinfo annotations prometheus -n monitoring

# Labels from any custom resource
kubectl getinfo labels mycustomresource

# Using different output formats
kubectl getinfo labels pods -o yaml
kubectl getinfo labels pods -o table
kubectl getinfo labels pods -o json  # default

# JSON with colors (similar to jq)
kubectl getinfo labels pods -o json -c
```

## Output Formats

The plugin supports three output formats, controlled by the `-o` or `--output` flag:

### JSON (default)

```bash
kubectl getinfo labels pods -o json
```

```json
{
  "items": [
    {
      "name": "pod-name",
      "namespace": "default",
      "labels": {
        "app": "nginx",
        "version": "1.0"
      }
    }
  ]
}
```

### YAML

```bash
kubectl getinfo labels pods -o yaml
```

```yaml
items:
  - name: pod-name
    namespace: default
    labels:
      app: nginx
      version: "1.0"
```

### Table

```bash
kubectl getinfo labels pods -o table
```

```
NAME        NAMESPACE    LABELS
pod-name    default      app=nginx,version=1.0
```

For annotations, the format is similar, but uses the `annotations` field instead of `labels`.

For ownerReferences:

```bash
kubectl getinfo owner pods -o table
```

```
NAME        NAMESPACE    OWNER NAMESPACE    OWNER KIND    OWNER NAME
pod-name    default      default            ReplicaSet    rs-name
```

#### OwnerReferences

```bash
# OwnerReferences for pods
kubectl getinfo owner pods -o json
```

```json
{
  "items": [
    {
      "name": "pod-name",
      "namespace": "default",
      "ownerReferences": [
        {
          "namespace": "default",
          "kind": "ReplicaSet",
          "name": "rs-name"
        }
      ]
    }
  ]
}
```

**Note:** If an object has no `ownerReferences`, the field is returned as an empty array `[]`.

#### Scheduling

The `scheduling` command lists all scheduling-related fields in pods that can affect the Kubernetes scheduler:

```bash
# All scheduling fields
kubectl getinfo scheduling pods -o json
```

```json
{
  "items": [
    {
      "name": "my-pod",
      "namespace": "default",
      "scheduling": {
        "nodeSelector": {
          "disktype": "ssd"
        },
        "affinity": {
          "nodeAffinity": { ... }
        },
        "tolerations": [ ... ],
        "resourceRequests": {
          "cpu": "2",
          "memory": "4Gi"
        },
        "resourceLimits": {
          "cpu": "4",
          "memory": "8Gi"
        },
        "topologySpreadConstraints": [ ... ],
        "schedulerName": "default-scheduler",
        "priorityClassName": "high-priority"
      }
    }
  ]
}
```

**Available subcommands:**
- `tolerations` - Lists tolerations only
- `affinity` - Lists affinity rules only
- `nodeselector` - Lists nodeSelector only
- `resources` - Lists resource requests/limits only
- `topology` - Lists topologySpreadConstraints only
- `priority` - Lists priority-related fields only
- `runtime` - Lists runtime-related fields only (runtimeClassName, hostNetwork, etc.)

**Example with subcommand:**
```bash
kubectl getinfo scheduling tolerations pods -o json
```

```json
{
  "items": [
    {
      "name": "my-pod",
      "namespace": "default",
      "tolerations": [
        {
          "key": "dedicated",
          "operator": "Equal",
          "value": "database",
          "effect": "NoSchedule"
        }
      ]
    }
  ]
}
```

**Note:** The `scheduling` command works with Pods and resources that have a Pod template (Deployments, StatefulSets, DaemonSets, Jobs, CronJobs, etc.). For template resources, fields are extracted from `spec.template.spec`.

### Colors in JSON

When using `-c` or `--color` with JSON output, the output is colorized using ANSI codes (similar to `jq`):

- **Keys**: bold blue
- **Strings**: green
- **Numbers**: yellow
- **Booleans**: bold yellow
- **Null**: gray
- **Punctuation** ({, }, [, ], :, ,): white

**Note**: Colors are only available in JSON output. YAML and table outputs do not support colors.

## Requirements

- `kubectl` configured and connected to a Kubernetes cluster
- Go 1.21 or newer (build only)


