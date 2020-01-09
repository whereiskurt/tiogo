package tenable

import (
	"encoding/json"
	"fmt"
	"time"
)

//ScannerList is from the Tenable.io documentation
type ScannerList struct {
	Scanners []Scanner
}

//Scanner is from the Tenable.io documentation
type Scanner struct {
	ID               json.Number `json:"id"`
	UUID             string      `json:"uuid"`
	Name             string      `json:"name"`
	Type             string      `json:"type"`
	Status           string      `json:"status"`
	ScanCount        json.Number `json:"scan_count"`
	EngineVersion    string      `json:"engine_version"`
	Platform         string      `json:"platform"`
	LoadedPluginSet  string      `json:"loaded_plugin_set"`
	RegistrationCode string      `json:"registration_code"`
	Owner            string      `json:"owner"`
	Key              string      `json:"key"`
	Addresses        []string    `json:"ip_addresses"`
	License          struct {
		Type         string      `json:"type"`
		IPS          json.Number `json:"ips"`
		Agents       json.Number `json:"agents"`
		Scanners     json.Number `json:"scanners"`
		AgentsUsed   json.Number `json:"agents_used"`
		ScannersUsed json.Number `json:"scanners_used"`
	}
}

// ScannerAgent is from the Tenable.io documentation - GET /scanners/{scanner_id}/agents
// Pagination is for groups of 5000 maximum limit
type ScannerAgent struct {
	Agents []struct {
		ID          json.Number `json:"id"`
		UUID        string      `json:"uuid"`
		Name        string      `json:"name"`
		Distro      string      `json:"distro"`
		IP          string      `json:"ip"`
		LastScanned json.Number `json:"last_scanned"`
		Platform    string      `json:"platform"`
		LinkedOn    json.Number `json:"linked_on"`
		LastConnect json.Number `json:"last_connect"`
		Feed        string      `json:"plugin_feed_id"`
		CoreBuild   string      `json:"core_build"`
		CoreVersion string      `json:"core_version"`
		Status      string      `json:"status"`
		Groups      []struct {
			ID   json.Number `json:"id"`
			Name string      `json:"name"`
		}
	}
	Pagination Pagination
}

// ScannerAgentGroups struct
type ScannerAgentGroups struct {
	Groups []ScannerAgentGroup
}

// ScannerAgentGroup struct
type ScannerAgentGroup struct {
	ID           json.Number `json:"id"`
	UUID         string      `json:"uuid"`
	Name         string      `json:"name"`
	AgentCount   json.Number `json:"agents_count"`
	LastModified json.Number `json:"last_modification_date"`
	Created      json.Number `json:"creation_date"`
}

// Pagination struct
type Pagination struct {
	ScanDetailHistory
	Total  json.Number `json:"total"`
	Offset json.Number `json:"offset"`
	Limit  json.Number `json:"limit"`
	Sort   []struct {
		Name  string `json:"name"`
		Order string `json:"order"`
	}
}

// ScansList struct // https://cloud.tenable.com/api#/resources/scans/
type ScansList struct {
	Folders []struct {
		ID json.Number `json:"id"`
	}
	Scans     []ScanListItem `json:"scans"`
	Timestamp json.Number    `json:"timestamp"`
}

// ScanListItem is returned for each scan
type ScanListItem struct {
	ID               json.Number `json:"id"`
	UUID             string      `json:"uuid"`
	ScheduleUUID     string      `json:"schedule_uuid"`
	Name             string      `json:"name"`
	Status           string      `json:"status"`
	Type             string      `json:"type"` //eg. agent,remote
	Owner            string      `json:"owner"`
	UserPermissions  json.Number `json:"user_permissions"`
	Permissions      json.Number `json:"permissions"`
	Enabled          bool        `json:"enabled"`
	Legacy           bool        `json:"legacy"`
	Read             bool        `json:"read"`
	Shared           bool        `json:"shared"`
	Control          bool        `json:"control"`
	RRules           string      `json:"rrules"`
	Timezone         string      `json:"timezone"`
	StartTime        string      `json:"startTime"`
	CreationDate     json.Number `json:"creation_date"`
	LastModifiedDate json.Number `json:"last_modification_date"`
}

// ScanDetail struct https://cloud.tenable.com/api#/resources/scans/{scanId}
type ScanDetail struct {
	Info            ScanDetailInfo
	Hosts           []ScanDetailHosts
	Vulnerabilities []ScanDetailVulnerabilities
	History         []ScanDetailHistory
}

// ScanDetailInfo struct
type ScanDetailInfo struct {
	ID           json.Number   `json:"object_id"`
	UUID         string        `json:"uuid"`
	ScheduleUUID string        `json:"schedule_uuid"`
	Owner        string        `json:"owner"`
	Start        json.Number   `json:"scan_start"` // Last Started Time
	End          json.Number   `json:"scan_end"`   // Last End Time - can be empty for running/never finsihed
	Timestamp    json.Number   `json:"timestamp"`  // Last time the info was updated
	ScannerStart json.Number   `json:"scanner_start"`
	ScannerEnd   json.Number   `json:"scanner_end"`
	ScannerName  string        `json:"scanner_name"`
	AgentCount   json.Number   `json:"agent_count"`
	AgentTarget  []AgentTarget `json:"agent_targets"`
	ScanType     string        `json:"scan_type"`
	HostCount    json.Number   `json:"hostcount"`
	Targets      string        `json:"targets"`
	PolicyName   string        `json:"policy"`
}

// AgentTarget struct
type AgentTarget struct {
	ID   json.Number
	UUID string
	Name string
}

// ScanDetailHosts struct
type ScanDetailHosts struct {
	ID               json.Number `json:"host_id"`
	AssetID          json.Number `json:"asset_id"`
	Index            json.Number `json:"host_index"`
	HostnameOrIP     string      `json:"hostname"` // the documentation is bad on this! It's actually IP address or NAME
	SeverityTotal    json.Number `json:"severity"`
	SeverityCritical json.Number `json:"critical"`
	SeverityHigh     json.Number `json:"high"`
	SeverityMedium   json.Number `json:"medium"`
	SeverityLow      json.Number `json:"low"`
	SeverityInfo     json.Number `json:"info"`
	Progress         string      `json:"progress"`
	Score            json.Number `json:"score"`
	ProgressCurrent  json.Number `json:"scanprogresscurrent"`
	ProgressTotal    json.Number `json:"scanprogresstotal"`
	ChecksConsidered json.Number `json:"numchecksconsidered"`
	ChecksTotal      json.Number `json:"totalchecksconsidered"`
}

// ScanDetailVulnerabilities struct
type ScanDetailVulnerabilities struct {
	ID       json.Number `json:"vuln_index"`
	PluginID json.Number `json:"plugin_id"`
	Name     string      `json:"plugin_name"`
	HostName string      `json:"hostname"`
	Family   string      `json:"plugin_family"`
	Count    json.Number `json:"count"`
	Severity json.Number `json:"severity"`
}

// ScanDetailHistory struct
type ScanDetailHistory struct {
	HistoryID        json.Number `json:"history_id"`
	UUID             string      `json:"uuid"`
	ScanType         string      `json:"type"`
	Status           string      `json:"status"`
	LastModifiedDate json.Number `json:"last_modification_date"`
	CreationDate     json.Number `json:"creation_date"`
}

// HostDetail struct // http://eagain.net/articles/go-dynamic-json/
// https://cloud.tenable.com/api#/resources/scans/{id}/host/{host_id}
type HostDetail struct {
	Info            HostDetailInfo
	Vulnerabilities []HostDetailVulnerabilities
}

// HostDetailInfo struct
type HostDetailInfo struct {
	HostStart       json.Number `json:"host_start"`       // becoming a number
	HostEnd         json.Number `json:"host_end"`         // becoming a number
	OperatingSystem []string    `json:"operating-system"` // becoming an array
	MACAddress      string      `json:"mac-address"`
	FQDN            string      `json:"host-fqdn"`
	NetBIOS         string      `json:"netbios-name"`
	HostIP          string      `json:"host-ip"`
}

// MarshalJSON for host detail infor // NOTE: This is needed for Marshal'ing back un-modified HotDetail object.
//      I think there is some type confusion where it treats the json.Numbers
//      as Date Strings... not sure what's up or why this is necessary entirely.
func (hdi HostDetailInfo) MarshalJSON() ([]byte, error) {
	var TimeFormatNoTZ = "Mon Jan _2 15:04:05 2006"

	type Alias HostDetailInfo

	stm, _ := time.Parse(TimeFormatNoTZ, fmt.Sprintf("%s", hdi.HostStart))
	etm, _ := time.Parse(TimeFormatNoTZ, fmt.Sprintf("%s", hdi.HostEnd))

	start := json.Number(fmt.Sprintf("%d", stm.Unix()))
	end := json.Number(fmt.Sprintf("%d", etm.Unix()))

	return json.Marshal(&struct {
		HostStart json.Number `json:"host_start"`
		HostEnd   json.Number `json:"host_end"`
		Alias
	}{
		HostStart: json.Number(start),
		HostEnd:   json.Number(end),
		Alias:     (Alias)(hdi),
	})
}

// HostDetailLegacyV2 struct from Tenable.io
type HostDetailLegacyV2 struct {
	Info            HostDetailInfoLegacyV2
	Vulnerabilities []HostDetailVulnerabilities
}

// HostDetailInfoLegacyV2 struct from Tenable.io
type HostDetailInfoLegacyV2 struct {
	HostStart       string `json:"host_start"`
	HostEnd         string `json:"host_end"`
	MACAddress      string `json:"mac-address"`
	FQDN            string `json:"host-fqdn"`
	NetBIOS         string `json:"netbios-name"`
	OperatingSystem string `json:"operating-system"`
	HostIP          string `json:"host-ip"`
}

//HostDetailVulnerabilities struct from Tenable.io
type HostDetailVulnerabilities struct {
	HostID       json.Number `json:"host_id"`
	HostName     string      `json:"hostname"`
	PluginID     json.Number `json:"plugin_id"`
	PluginName   string      `json:"plugin_name"`
	PluginFamily string      `json:"plugin_family"`
	Count        json.Number `json:"count"`
	Severity     json.Number `json:"severity"`
}

//PluginFamilies struct from Tenable.io
type PluginFamilies struct {
	Families []struct {
		ID    json.Number `json:"id"`
		Name  string      `json:"name"`
		Count json.Number `json:"count"`
	}
}

//FamilyPlugins struct from Tenable.io
type FamilyPlugins struct {
	ID      json.Number `json:"id"`
	Name    string      `json:"name"`
	Plugins []Plugin
}

// Plugin struct // https://cloud.tenable.com/api#/resources/plugins/plugin/{pluginId}
// NOTE: A cache record would basically never goes stale.
//Plugin struct from Tenable.io
type Plugin struct {
	ID         json.Number `json:"id"`
	Name       string      `json:"name"`
	FamilyName string      `json:"family_name"`
	Attributes []struct {
		Name  string `json:"attribute_name"`
		Value string `json:"attribute_value"`
	}
}

//TagCategories struct from Tenable.io
type TagCategories struct {
	Categories []TagCategory
}

//TagCategory struct from Tenable.io
type TagCategory struct {
	ContainerUUID string `json:"container_uuid"`
	UUID          string `json:"uuid"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
	ModelName     string `json:"model_name"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

//TagValues struct from Tenable.io
type TagValues struct {
	Values []TagValue
}

//TagValue struct from Tenable.io
type TagValue struct {
	ContainerUUID       string `json:"container_uuid"`
	UUID                string `json:"uuid"`
	CreatedAt           string `json:"created_at"`
	CreatedBy           string `json:"created_by"`
	UpdatedAt           string `json:"updated_at"`
	UpdatedBy           string `json:"updated_by"`
	ModelName           string `json:"model_name"`
	Value               string `json:"value"`
	Description         string `json:"description"`
	Type                string `json:"type"`
	CategoryUUID        string `json:"category_uuid"`
	CategoryName        string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
}

//AssetHost struct from Tenable.io
type AssetHost struct {
	Assets []struct {
		HostID     json.Number `json:"id"`
		UUID       string      `json:"uuid"`
		FQDN       []string    `json:"fqdn"`
		IPV4       []string    `json:"ipv4"`
		IPV6       []string    `json:"ipv6"`
		Severities []struct {
			Count json.Number `json:"count"`
			Level json.Number `json:"level"`
			Name  string      `json:"name"`
		}
	}
}

//AssetSearchResults struct from Tenable.io
type AssetSearchResults struct {
	Assets []AssetInfo
	Total  json.Number `json:"total"`
}

//Asset struct from Tenable.io
type Asset struct {
	Info AssetInfo
}

//AssetInfo struct from Tenable.io // https://cloud.tenable.com/api#/resources/workbenches/asset-info
type AssetInfo struct {
	ID                      string   `json:"id"`
	TimeEnd                 string   `json:"time_end"`
	UUID                    string   `json:"uuid"`
	OperatingSystem         []string `json:"operating_system"`
	HasAgent                bool     `json:"has_agent"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
	FirstSeenAt             string   `json:"first_seen"`
	LastSeenAt              string   `json:"last_seen"`
	LastAuthenticatedScanAt string   `json:"last_authenticated_scan_date"`
	LastLicensedScanAt      string   `json:"last_licensed_scan_date"`
	IPV4                    []string `json:"ipv4"`
	IPV6                    []string `json:"ipv6"`
	FQDN                    []string `json:"fqdn"`
	MACAddress              []string `json:"mac_address"`
	NetBIOS                 []string `json:"netbios_name"`
	SystemType              []string `json:"system_type"`
	TenableUUID             []string `json:"tenable_uuid"` // NOTE: This is Agent UUID!!
	HostName                []string `json:"hostname"`
	AgentName               []string `json:"agent_name"`
	BIOSUUID                []string `json:"bios_uuid"`
	AWSEC2InstanceID        []string `json:"aws_ec2_instance_id"`
	AWSEC2InstanceAMIID     []string `json:"aws_ec2_instance_ami_id"`
	AWSOwnerID              []string `json:"aws_owner_id"`
	AWSAvailabilityZone     []string `json:"aws_availability_zone"`
	AWSRegion               []string `json:"aws_region"`
	AWSVPCID                []string `json:"aws_vpc_id"`
	AWSEC2InstanceGroupName []string `json:"aws_ec2_instance_group_name"`
	AWSEC2InstanceStateName []string `json:"aws_ec2_instance_state_name"`
	AWSEC2InstanceType      []string `json:"aws_ec2_instance_type"`
	AWSSubnetID             []string `json:"aws_subnet_id"`
	AWSEC2ProductCode       []string `json:"aws_ec2_product_code"`
	AWSEC2Name              []string `json:"aws_ec2_name"`
	AzureVMId               []string `json:"azure_vm_id"`
	AzureResourceID         []string `json:"azure_resource_id"`
	SSHFingerPrint          []string `json:"ssh_fingerprint"`
	McafeeEPOGUID           []string `json:"mcafee_epo_guid"`
	McafeeEPOAgentGUID      []string `json:"mcafee_epo_agent_guid"`
	QualysHostID            []string `json:"qualys_host_id"`
	QualysAssetID           []string `json:"qualys_asset_id"`
	ServiceNowSystemID      []string `json:"servicenow_sysid"`
	Counts                  struct {
		Vulnerabilities struct {
			Total      json.Number `json:"total"`
			Severities []struct {
				Count json.Number `json:"count"`
				Level json.Number `json:"level"`
				Name  string      `json:"name"`
			}
		}
		Audits struct {
			Total      json.Number `json:"total"`
			Severities []struct {
				Count json.Number `json:"count"`
				Level json.Number `json:"level"`
				Name  string      `json:"name"`
			}
		}
	}
	Interfaces []struct {
		Name       string   `json:"name"`
		IPV4       []string `json:"ipv4"`
		IPV6       []string `json:"ipv6"`
		FQDN       []string `json:"fqdn"`
		MACAddress []string `json:"mac_address"`
	}
	Sources []struct {
		FirstSeenAt string `json:"first_seen"`
		LastSeenAt  string `json:"last_seen"`
		Name        string `json:"name"`
	}
	Tags []struct {
		UUID         string `json:"tag_uuid"`
		CategoryName string `json:"tag_key"`
		Value        string `json:"tag_value"`
		AddedBy      string `json:"added_by"`
		AddedAt      string `json:"added_at"`
		Source       string `json:"source"`
	}
}

//AssetVuln struct from Tenable.io // GET /workbenches/assets/{asset_id}/vulnerabilities
type AssetVuln struct {
	Vulnerabilities []struct {
		PluginID     json.Number `json:"plugin_id"`
		PluginName   string      `json:"plugin_name"`
		PluginFamily string      `json:"plugin_family"`
		Count        json.Number `json:"count"`
		State        string      `json:"vulnerability_state"`
		Severity     json.Number `json:"severity"`
	}
}

//AssetVulnInfo struct from Tenable.io // GET /workbenches/assets/{asset_id}/vulnerabilities/{plugin_id}/info
type AssetVulnInfo struct {
	Info struct {
		Description   string      `json:"description"`
		Solution      string      `json:"solution"`
		Synopsis      string      `json:"synopsis"`
		Count         json.Number `json:"count"`
		Severity      json.Number `json:"severity"`
		Discovery     interface{} `json:"discovery"`
		PluginDetails interface{} `json:"plugin_details"`
		ReferenceInfo interface{} `json:"reference_information"`
		RiskInfo      interface{} `json:"risk_information"`
		SeeAlso       interface{} `json:"see_also"`
		VulnInfo      interface{} `json:"vulnerability_information"`
	}
}

//AssetVulnOutput struct from Tenable.io // GET /workbenches/assets/{asset_id}/vulnerabilities/{plugin_id}/outputs
type AssetVulnOutput struct {
	Outputs []struct {
		PluginOutput string `json:"plugin_output"`
		States       []struct {
			Name   string `json:"name"`
			Result []AssetVulnResult
		}
	}
}

//AssetVulnResult struct from Tenable.io
type AssetVulnResult struct {
	ApplicationProtocol string        `json:"application_protocol"`
	TransportProtocol   string        `json:"transport_protocol"`
	Port                json.Number   `json:"port"`
	Severity            json.Number   `json:"severity"`
	Assets              []interface{} `json:"assets"`
}

//AssetExportStart struct from Tenable.io
type AssetExportStart struct {
	UUID string `json:"export_uuid"`
}

//AssetExportStatus struct from Tenable.io
type AssetExportStatus struct {
	Status          string        `json:"status"`
	Chunks          []json.Number `json:"chunks_available"`
	ChunksFailed    []json.Number `json:"chunks_failed"`
	ChunksCancelled []json.Number `json:"chunks_cancelled"`
}

//AssetExportChunk struct from Tenable.io
type AssetExportChunk struct {
	UUID                    string   `json:"id"`
	HasAgent                bool     `json:"has_agent"`
	HasPluginResult         bool     `json:"has_plugin_results"`
	CreatedAt               string   `json:"created_at"`
	TerminatedAt            string   `json:"terminated_at"`
	TerminatedBy            string   `json:"terminated_by"`
	UpdatedAt               string   `json:"updated_at"`
	DeletedAt               string   `json:"deleted_at"`
	DeletedBy               string   `json:"deleted_by"`
	FirstSeenAt             string   `json:"first_seen"`
	LastSeenAt              string   `json:"last_seen"`
	FirstScanAt             string   `json:"first_scan_time"`
	LastScanAt              string   `json:"last_scan_time"`
	LastAuthenticatedScanAt string   `json:"last_authenticated_scan_date"`
	LastLicensedScanAt      string   `json:"last_licensed_scan_date"`
	AzureVMId               string   `json:"azure_vm_id"`
	AzureResourceID         string   `json:"azure_resource_id"`
	AWSEC2InstanceAMIId     string   `json:"aws_ec2_instance_ami_id"`
	AWSEC2InstanceID        string   `json:"aws_ec2_instance_id"`
	AWSOwnerID              string   `json:"aws_owner_id"`
	AgentUUID               string   `json:"agent_uuid"`
	BIOSUUID                string   `json:"bios_uuid"`
	AWSAvailabilityZone     string   `json:"aws_availability_zone"`
	AWSRegion               string   `json:"aws_region"`
	AWSVPCID                string   `json:"aws_vpc_id"`
	AWSEC2InstanceGroupName string   `json:"aws_ec2_instance_group_name"`
	AWSEC2InstanceStateName string   `json:"aws_ec2_instance_state_name"`
	AWSEC2InstanceType      string   `json:"aws_ec2_instance_type"`
	AWSSubnetID             string   `json:"aws_subnet_id"`
	AWSEC2ProductCode       string   `json:"aws_ec2_product_code"`
	AWSEC2Name              string   `json:"aws_ec2_name"`
	AgentNames              []string `json:"agent_names"`
	McafeeEPOGUID           []string `json:"mcafee_epo_guid"`
	McafeeEPOAgentGUID      []string `json:"mcafee_epo_agent_guid"`
	EnvironmentID           string   `json:"environment_uuid"`
	IPV4                    []string `json:"ipv4s"`
	IPV6                    []string `json:"ipv6s"`
	FQDN                    []string `json:"fqdns"`
	MACAddress              []string `json:"mac_addresses"`
	NetBIOS                 []string `json:"netbios_names"`
	OperatingSystem         []string `json:"operating_systems"`
	SystemType              []string `json:"system_types"`
	HostName                []string `json:"hostnames"`
	SSHFingerPrint          []string `json:"ssh_fingerprints"`
	QualysAssetID           []string `json:"qualys_asset_id"`
	QualysHostID            []string `json:"qualys_host_id"`
	ServiceNowSystemID      []string `json:"servicenow_sysid"`
	Sources                 []struct {
		FirstSeenAt string `json:"first_seen"`
		LastSeenAt  string `json:"last_seen"`
		Name        string `json:"name"`
	}
	Tags []struct {
		UUID         string `json:"tag_uuid"`
		CategoryName string `json:"tag_key"`
		Value        string `json:"tag_value"`
		AddedBy      string `json:"added_by"`
		AddedAt      string `json:"added_at"`
		Source       string `json:"source"`
	}
	Interfaces []struct {
		Name       string   `json:"name"`
		IPV4       []string `json:"ipv4"`
		IPV6       []string `json:"ipv6"`
		FQDN       []string `json:"fqdn"`
		MACAddress []string `json:"mac_address"`
	}
}

//VulnExportStart struct from Tenable.io
type VulnExportStart struct {
	UUID string `json:"export_uuid"`
}

//VulnExportStatus struct from Tenable.io
type VulnExportStatus struct {
	Status          string        `json:"status"`
	Chunks          []json.Number `json:"chunks_available"`
	ChunksFailed    []json.Number `json:"chunks_failed"`
	ChunksCancelled []json.Number `json:"chunks_cancelled"`
}

//VulnExportChunk struct from Tenable.io
type VulnExportChunk struct {
	Asset struct {
		DeviceType               string   `json:"device_type"`
		FQDN                     string   `json:"fqdn"`
		HostName                 string   `json:"hostname"`
		UUID                     string   `json:"uuid"`
		IPV4                     string   `json:"ipv4"`
		LastAuthenticatedResults string   `json:"last_unauthenticated_results"`
		NETBIOSWorkgroup         []string `json:"netbios_workgroup"`
		OperatingSystem          []string `json:"operating_system"`
		Tracked                  bool     `json:"tracked"`
	}
	Output string
	Plugin struct {
		Description      string      `json:"description"`
		Family           string      `json:"family"`
		FamilyID         json.Number `json:"family_id"`
		HasPatch         bool        `json:"has_patch"`
		PluginID         json.Number `json:"id"`
		Name             string      `json:"name"`
		ModificationDate string      `json:"modification_date"`
		PublicationDate  string      `json:"publication_date"`
		RiskFactor       string      `json:"risk_factor"`
		Solution         string      `json:"solution"`
		Synopsis         string      `json:"synopsis"`
		Type             string      `json:"type"`
		Version          string      `json:"version"`
	}
	Port struct {
		Port     json.Number `json:"port"`
		Protocol string      `json:"protocol"`
		Service  string      `json:"service"`
	}
	Scan struct {
		CompletedAt  string `json:"completed_at"`
		ScheduleUUID string `json:"schedule_uuid"`
		StartedAt    string `json:"started_at"`
		UUID         string `json:"uuid"`
	}
	Severity             string      `json:"severity"`
	SeverityID           json.Number `json:"severity_id"`
	SeverityDefaultID    json.Number `json:"severity_default_id"`
	SeverityModification string      `json:"severity_modification_type"`
	FirstFound           string      `json:"first_found"`
	LastFound            string      `json:"last_found"`
	State                string      `json:"state"`
}

// ExportFilter is shared and not the same for
//ExportFilter struct from Tenable.io
type ExportFilter struct {
	ExportRequest string      `json:"export-request"`
	Limit         json.Number `json:"chunk_size"`
	Filters       struct {
		Since        json.Number `json:"since"`
		LastAssessed json.Number `json:"last_assessed"`
	} `json:"filters"`
}

// ScansExportStart is outputed at successful scans export
// Note: FileUUID is a uuid other times it's an unquoted number! (ie. pdf)
//ScansExportStart struct from Tenable.io
type ScansExportStart struct {
	FileUUID  DownloadFileID `json:"file"`
	TempToken string         `json:"temp_token"`
}

// ScansExportStatus returns 'ready' when done
//ScansExportStatus struct from Tenable.io
type ScansExportStatus struct {
	Status string `json:"status"`
}

// ScansExportStartPost returns 'ready' when done
type ScansExportStartPost struct {
	Format   string `json:"format"`
	Chapters string `json:"chapters"`
}

// DownloadFileID allows us to create customer marshal/unmarshal code for this type
// that isn't always quoted/unquoted
type DownloadFileID string

// MarshalJSON will JSONinfy DownloadFileID
func (f DownloadFileID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

// UnmarshalJSON will allow us to decode JSON numbers as strings
func (f *DownloadFileID) UnmarshalJSON(data []byte) (err error) {
	if data[0] == '"' {
		var s string
		err = json.Unmarshal(data, &s)
		*f = DownloadFileID(s)
	} else {
		var v json.Number
		err = json.Unmarshal(data, &v)
		*f = DownloadFileID(v)
	}
	return err
}
