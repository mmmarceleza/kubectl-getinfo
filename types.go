package main

// OwnerReference represents a reference to an owner of a Kubernetes resource
type OwnerReference struct {
	Namespace string `json:"namespace,omitempty"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
}

// ContainerResources represents resource requests and limits for a single container
type ContainerResources struct {
	Name     string                 `json:"name" yaml:"name"`
	Requests map[string]interface{} `json:"requests,omitempty" yaml:"requests,omitempty"`
	Limits   map[string]interface{} `json:"limits,omitempty" yaml:"limits,omitempty"`
}

// SchedulingInfo contains scheduling-related fields from a pod spec
type SchedulingInfo struct {
	NodeSelector              map[string]string      `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	NodeName                  string                 `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	Affinity                  map[string]interface{} `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Tolerations               []interface{}          `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	TopologySpreadConstraints []interface{}          `json:"topologySpreadConstraints,omitempty" yaml:"topologySpreadConstraints,omitempty"`
	ResourceRequests          map[string]interface{} `json:"resourceRequests,omitempty" yaml:"resourceRequests,omitempty"`
	ResourceLimits            map[string]interface{} `json:"resourceLimits,omitempty" yaml:"resourceLimits,omitempty"`
	SchedulerName             string                 `json:"schedulerName,omitempty" yaml:"schedulerName,omitempty"`
	PriorityClassName         string                 `json:"priorityClassName,omitempty" yaml:"priorityClassName,omitempty"`
	Priority                  *int32                 `json:"priority,omitempty" yaml:"priority,omitempty"`
	PreemptionPolicy          string                 `json:"preemptionPolicy,omitempty" yaml:"preemptionPolicy,omitempty"`
	RuntimeClassName          string                 `json:"runtimeClassName,omitempty" yaml:"runtimeClassName,omitempty"`
	HostNetwork               bool                   `json:"hostNetwork,omitempty" yaml:"hostNetwork,omitempty"`
	HostPID                   bool                   `json:"hostPID,omitempty" yaml:"hostPID,omitempty"`
	HostIPC                   bool                   `json:"hostIPC,omitempty" yaml:"hostIPC,omitempty"`
}

// OutputItem represents a single resource in the output
type OutputItem struct {
	Name            string             `json:"name"`
	Namespace       string             `json:"namespace,omitempty"`
	Labels          *map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     *map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	OwnerReferences []OwnerReference   `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Scheduling      *SchedulingInfo    `json:"scheduling,omitempty" yaml:"scheduling,omitempty"`
	// Specific fields for scheduling subcommands
	Tolerations               []interface{}          `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Affinity                  map[string]interface{} `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	NodeSelector              map[string]string      `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Resources                 []ContainerResources   `json:"resources,omitempty" yaml:"resources,omitempty"`
	TopologySpreadConstraints []interface{}          `json:"topologySpreadConstraints,omitempty" yaml:"topologySpreadConstraints,omitempty"`
	Priority                  map[string]interface{} `json:"priority,omitempty" yaml:"priority,omitempty"`
	Runtime                   map[string]interface{} `json:"runtime,omitempty" yaml:"runtime,omitempty"`
}

// Output represents the complete output structure
type Output struct {
	Items []OutputItem `json:"items"`
}

