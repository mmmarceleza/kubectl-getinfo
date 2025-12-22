package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// getPodSpecPath returns the path to the pod spec based on the resource kind
func getPodSpecPath(item unstructured.Unstructured) []string {
	kind := item.GetKind()

	// Para Pods, o spec está diretamente em spec
	if kind == "Pod" {
		return []string{"spec"}
	}

	// Para recursos com template (Deployments, StatefulSets, etc.)
	// O spec do pod está em spec.template.spec
	templateKinds := []string{
		"Deployment", "StatefulSet", "DaemonSet", "ReplicaSet",
		"Job", "CronJob",
	}

	for _, tk := range templateKinds {
		if kind == tk {
			return []string{"spec", "template", "spec"}
		}
	}

	// Default: tentar spec diretamente
	return []string{"spec"}
}

// extractOwnerReferences extracts owner references from a resource
func extractOwnerReferences(item unstructured.Unstructured) []OwnerReference {
	ownerRefs := []OwnerReference{}

	// Get ownerReferences from metadata
	metadata, found, err := unstructured.NestedSlice(item.Object, "metadata", "ownerReferences")
	if !found || err != nil {
		return ownerRefs
	}

	for _, ref := range metadata {
		refMap, ok := ref.(map[string]interface{})
		if !ok {
			continue
		}

		ownerRef := OwnerReference{}

		// Extract kind
		if kind, ok := refMap["kind"].(string); ok {
			ownerRef.Kind = kind
		}

		// Extract name
		if name, ok := refMap["name"].(string); ok {
			ownerRef.Name = name
		}

		// Extract namespace (may not be present in all cases)
		// If ownerReference doesn't have namespace, use the namespace of the current object
		if namespace, ok := refMap["namespace"].(string); ok && namespace != "" {
			ownerRef.Namespace = namespace
		} else {
			// Use the namespace of the object containing the ownerReference
			ownerRef.Namespace = item.GetNamespace()
		}

		ownerRefs = append(ownerRefs, ownerRef)
	}

	return ownerRefs
}

// extractSchedulingInfo extracts all scheduling-related information from a resource
func extractSchedulingInfo(item unstructured.Unstructured) *SchedulingInfo {
	specPath := getPodSpecPath(item)
	scheduling := &SchedulingInfo{}

	// NodeSelector
	if nodeSelector, found, _ := unstructured.NestedStringMap(item.Object, append(specPath, "nodeSelector")...); found && len(nodeSelector) > 0 {
		scheduling.NodeSelector = nodeSelector
	}

	// NodeName
	if nodeName, found, _ := unstructured.NestedString(item.Object, append(specPath, "nodeName")...); found && nodeName != "" {
		scheduling.NodeName = nodeName
	}

	// Affinity
	if affinity, found, _ := unstructured.NestedMap(item.Object, append(specPath, "affinity")...); found && len(affinity) > 0 {
		scheduling.Affinity = affinity
	}

	// Tolerations
	if tolerations, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "tolerations")...); found && len(tolerations) > 0 {
		scheduling.Tolerations = tolerations
	}

	// TopologySpreadConstraints
	if topology, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "topologySpreadConstraints")...); found && len(topology) > 0 {
		scheduling.TopologySpreadConstraints = topology
	}

	// Resource Requests and Limits (from containers)
	if containers, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "containers")...); found {
		requests := make(map[string]interface{})
		limits := make(map[string]interface{})

		for _, container := range containers {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				continue
			}

			if resources, ok := containerMap["resources"].(map[string]interface{}); ok {
				if req, ok := resources["requests"].(map[string]interface{}); ok {
					for k, v := range req {
						if existing, exists := requests[k]; exists {
							// Sum resources if multiple containers
							if existingStr, ok := existing.(string); ok {
								if vStr, ok := v.(string); ok {
									// Simple string comparison for now
									requests[k] = existingStr + "," + vStr
								}
							}
						} else {
							requests[k] = v
						}
					}
				}
				if lim, ok := resources["limits"].(map[string]interface{}); ok {
					for k, v := range lim {
						if existing, exists := limits[k]; exists {
							if existingStr, ok := existing.(string); ok {
								if vStr, ok := v.(string); ok {
									limits[k] = existingStr + "," + vStr
								}
							}
						} else {
							limits[k] = v
						}
					}
				}
			}
		}

		if len(requests) > 0 {
			scheduling.ResourceRequests = requests
		}
		if len(limits) > 0 {
			scheduling.ResourceLimits = limits
		}
	}

	// SchedulerName
	if schedulerName, found, _ := unstructured.NestedString(item.Object, append(specPath, "schedulerName")...); found && schedulerName != "" {
		scheduling.SchedulerName = schedulerName
	}

	// PriorityClassName
	if priorityClassName, found, _ := unstructured.NestedString(item.Object, append(specPath, "priorityClassName")...); found && priorityClassName != "" {
		scheduling.PriorityClassName = priorityClassName
	}

	// Priority
	if priority, found, _ := unstructured.NestedInt64(item.Object, append(specPath, "priority")...); found {
		priorityInt32 := int32(priority)
		scheduling.Priority = &priorityInt32
	}

	// PreemptionPolicy
	if preemptionPolicy, found, _ := unstructured.NestedString(item.Object, append(specPath, "preemptionPolicy")...); found && preemptionPolicy != "" {
		scheduling.PreemptionPolicy = preemptionPolicy
	}

	// RuntimeClassName
	if runtimeClassName, found, _ := unstructured.NestedString(item.Object, append(specPath, "runtimeClassName")...); found && runtimeClassName != "" {
		scheduling.RuntimeClassName = runtimeClassName
	}

	// HostNetwork
	if hostNetwork, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostNetwork")...); found {
		scheduling.HostNetwork = hostNetwork
	}

	// HostPID
	if hostPID, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostPID")...); found {
		scheduling.HostPID = hostPID
	}

	// HostIPC
	if hostIPC, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostIPC")...); found {
		scheduling.HostIPC = hostIPC
	}

	// Return nil if no scheduling info found
	if scheduling.NodeSelector == nil && scheduling.NodeName == "" && scheduling.Affinity == nil &&
		len(scheduling.Tolerations) == 0 && len(scheduling.TopologySpreadConstraints) == 0 &&
		scheduling.ResourceRequests == nil && scheduling.ResourceLimits == nil &&
		scheduling.SchedulerName == "" && scheduling.PriorityClassName == "" && scheduling.Priority == nil &&
		scheduling.PreemptionPolicy == "" && scheduling.RuntimeClassName == "" &&
		!scheduling.HostNetwork && !scheduling.HostPID && !scheduling.HostIPC {
		return nil
	}

	return scheduling
}

// extractSchedulingSubcommand extracts a specific scheduling field based on subcommand
func extractSchedulingSubcommand(item unstructured.Unstructured, outputItem *OutputItem, subCommand string) {
	specPath := getPodSpecPath(item)

	switch subCommand {
	case "tolerations":
		if tolerations, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "tolerations")...); found && len(tolerations) > 0 {
			outputItem.Tolerations = tolerations
		}
	case "affinity":
		if affinity, found, _ := unstructured.NestedMap(item.Object, append(specPath, "affinity")...); found && len(affinity) > 0 {
			outputItem.Affinity = affinity
		}
	case "nodeselector":
		if nodeSelector, found, _ := unstructured.NestedStringMap(item.Object, append(specPath, "nodeSelector")...); found && len(nodeSelector) > 0 {
			outputItem.NodeSelector = nodeSelector
		}
	case "resources":
		resources := make(map[string]interface{})
		if containers, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "containers")...); found {
			requests := make(map[string]interface{})
			limits := make(map[string]interface{})

			for _, container := range containers {
				containerMap, ok := container.(map[string]interface{})
				if !ok {
					continue
				}

				if res, ok := containerMap["resources"].(map[string]interface{}); ok {
					if req, ok := res["requests"].(map[string]interface{}); ok {
						for k, v := range req {
							requests[k] = v
						}
					}
					if lim, ok := res["limits"].(map[string]interface{}); ok {
						for k, v := range lim {
							limits[k] = v
						}
					}
				}
			}

			if len(requests) > 0 {
				resources["requests"] = requests
			}
			if len(limits) > 0 {
				resources["limits"] = limits
			}
		}
		if len(resources) > 0 {
			outputItem.Resources = resources
		}
	case "topology":
		if topology, found, _ := unstructured.NestedSlice(item.Object, append(specPath, "topologySpreadConstraints")...); found && len(topology) > 0 {
			outputItem.TopologySpreadConstraints = topology
		}
	case "priority":
		priority := make(map[string]interface{})
		if priorityClassName, found, _ := unstructured.NestedString(item.Object, append(specPath, "priorityClassName")...); found && priorityClassName != "" {
			priority["priorityClassName"] = priorityClassName
		}
		if prio, found, _ := unstructured.NestedInt64(item.Object, append(specPath, "priority")...); found {
			priority["priority"] = prio
		}
		if preemptionPolicy, found, _ := unstructured.NestedString(item.Object, append(specPath, "preemptionPolicy")...); found && preemptionPolicy != "" {
			priority["preemptionPolicy"] = preemptionPolicy
		}
		if len(priority) > 0 {
			outputItem.Priority = priority
		}
	case "runtime":
		runtime := make(map[string]interface{})
		if runtimeClassName, found, _ := unstructured.NestedString(item.Object, append(specPath, "runtimeClassName")...); found && runtimeClassName != "" {
			runtime["runtimeClassName"] = runtimeClassName
		}
		if hostNetwork, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostNetwork")...); found {
			runtime["hostNetwork"] = hostNetwork
		}
		if hostPID, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostPID")...); found {
			runtime["hostPID"] = hostPID
		}
		if hostIPC, found, _ := unstructured.NestedBool(item.Object, append(specPath, "hostIPC")...); found {
			runtime["hostIPC"] = hostIPC
		}
		if len(runtime) > 0 {
			outputItem.Runtime = runtime
		}
	}
}

