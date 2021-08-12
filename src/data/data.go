package data

type Log struct {
	CreatedOn string `json:"created_on"`
	Component string `json:"component"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type MyLog struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	InHits []struct {
		Source Source `json:"_source"`
	} `json:"hits"`
}

type Source struct {
	LogDate    []string `json:"log_date"`
	LogMessage []string `json:"logmessage"`
	Fields     Fields   `json:"fields"`
	LogLevel   []string `json:"log_level"`
}

type Fields struct {
	LogType string `json:"log_type"`
}

type ComponentError struct {
	Nova       bool `json:"nova"`
	Heat       bool `json:"heat"`
	Cinder     bool `json:"cinder"`
	Neutron    bool `json:"neutron"`
	Keystone   bool `json:"keystone"`
	Swift      bool `json:"swift"`
	Agent      bool `json:"agent"`
	Management bool `json:"management"`
}

type TokenRequest struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Identity Identity `json:"identity"`
	Scope    *Scope   `json:"scope,omitempty"`
}

type Identity struct {
	Methods  []string `json:"methods"`
	Password Password `json:"password"`
}

type Password struct {
	User User `json:"user"`
}

type User struct {
	Name     string  `json:"name,omitempty"`
	Domain   *Domain `json:"domain,omitempty"`
	Id       string  `json:"id,omitempty"`
	Password string  `json:"password"`
}

type Scope struct {
	System  *System `json:"system,omitempty"`
	Domain  *Domain `json:"domain,omitempty"`
	Project Project `json:"project,omitempty"`
}

type System struct {
	All bool `json:"all"`
}

type Domain struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Project struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Login struct {
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	ProjectId string `json:"project_id"`
}

type Metrics struct {
	OpenStackMetrics  Statistics `json:"openstack_metrics"`
	CloudStackMetrics Statistics `json:"cloudstack_metrics"`
}

type Statistics struct {
	VCPU        int `json:"vcpu"`
	VCPUUsed    int `json:"vcpu_used"`
	Memory      int `json:"memory"`
	MemoryUsed  int `json:"memory_used"`
	Storage     int `json:"storage"`
	StorageUsed int `json:"storage_used"`
}

type Hypervisors struct {
	Statistics OpenStackStatistics `json:"hypervisor_statistics"`
}

type OpenStackStatistics struct {
	VCPUsUsed          int `json:"vcpus_used"`
	LocalGBUsed        int `json:"local_gb_used"`
	MemoryMB           int `json:"memory_mb"`
	CurrentWorkload    int `json:"current_workload"`
	VCPUs              int `json:"vcpus"`
	RunningVMs         int `json:"running_vms"`
	FreeDiskGB         int `json:"free_disk_gb"`
	DiskAvailableLeast int `json:"disk_available_least"`
	LocalGB            int `json:"local_gb"`
	FreeRamMB          int `json:"free_ram_mb"`
	MemoryMBUsed       int `json:"memory_mb_used"`
}

type CloudStackMetrics struct {
	Response ListCapacityResponse `json:"listcapacityresponse"`
}

type ListCapacityResponse struct {
	Count    int        `json:"count"`
	Capacity []Capacity `json:"capacity"`
}

type Capacity struct {
	Type              int    `json:"type"`
	Name              string `json:"name"`
	ZoneId            string `json:"zoneid"`
	ZoneName          string `json:"zonename"`
	CapacityAllocated int    `json:"capacityalocated"`
	CapacityUsed      int    `json:"capacityused"`
	CapacityTotal     int    `json:"capacitytotal"`
	PercentUsed       string `json:"percentusd"`
}
