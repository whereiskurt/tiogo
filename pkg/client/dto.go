package client

import (
	"time"
)

// Scan DTO is an all string representation
type Scan struct {
	ScanID           string
	UUID             string
	ScheduleUUID     string // Considered the UUID for ScanDetails etc.
	Name             string
	Type             string
	Status           string
	Owner            string
	UserPermissions  string
	Enabled          string
	RRules           string
	Timezone         string
	StartTime        string
	CreationDate     string
	LastModifiedDate string
	Timestamp        string
}

// type ScanHistory struct {
// 	Scan    Scan
// 	History []ScanHistoryDetail
// }

type ScanHistoryDetail struct {
	Scan                *Scan
	HistoryID           string
	HistoryCount        string
	ScannerName         string
	PolicyName          string
	Targets             string
	AgentGroup          []AgentGroup
	Owner               string
	HistoryIndex        string
	Status              string
	CreationDate        string
	LastModifiedDate    string
	PluginCriticalCount string
	PluginHighCount     string
	PluginMediumCount   string
	PluginLowCount      string
	PluginInfoCount     string
	PluginTotalCount    string
	ScanStart           string
	ScanStartUnix       string
	ScanEnd             string
	ScanEndUnix         string
	Timestamp           string
	TimestampUnix       string
	ScanDuration        string
	HostCount           string
	AgentCount          string
	ScanType            string
	Host                map[string]HostScanSummary
	HostPlugin          map[string]Plugin
	HostAssetMap        map[string]string
}

type HostScanSummary struct {
	HostID              string
	AssetID             string
	HostnameOrIP        string
	ScanHistoryDetail   *ScanHistoryDetail
	HostDetail          HostScanDetail
	Asset               Asset
	PluginCriticalCount string
	PluginHighCount     string
	PluginMediumCount   string
	PluginLowCount      string
	PluginTotalCount    string
	Score               string
	Progress            string
	ScanProgressCurrent string
	ScanProgressTotal   string
	ChecksConsidered    string
	ChecksTotal         string
}

func (h *HostScanSummary) HasAsset() (hasAsset bool) {
	if h.Asset.UUID != "" {
		hasAsset = true
	}
	return
}

type HostScanDetail struct {
	IP               string
	FQDN             string
	NetBIOS          string
	MACAddresses     string
	OperatingSystems string
	ScanStart        string
	ScanStartUnix    string
	ScanEnd          string
	ScanEndUnix      string
	ScanDuration     string
	PluginMap        map[string]Plugin
}

type Plugin struct {
	PluginID   string
	Name       string
	FamilyName string
	FamilyID   string
	Count      string
	Severity   string
	Detail     PluginDetail
}
type PluginDetail struct {
	FunctionName          string
	PluginPublicationDate string
	PatchPublicationDate  string
	Attribute             map[string]PluginDetailAttribute
}
type PluginDetailAttribute struct {
	Name  string
	Value string
}
type PluginFamily struct {
	ID    string
	Name  string
	Count string
}
type PluginFamilyDetail struct {
	ID       string
	Name     string
	PluginID []string
}

type Asset struct {
	TimeEnd                 string
	ID                      string
	UUID                    string
	OperatingSystem         []string
	TenableUUID             []string
	IPV4                    []string
	IPV6                    []string
	FQDN                    []string
	MACAddress              []string
	NetBIOS                 []string
	SystemType              []string
	HostName                []string
	AgentName               []string
	BIOSUUID                []string
	HasAgent                bool
	CreatedAt               string
	UpdatedAt               string
	FirstSeenAt             string
	LastSeenAt              string
	LastAuthenticatedScanAt string
	LastLicensedScanAt      string
	AWSEC2InstanceID        []string
	AWSEC2InstanceAMIID     []string
	AWSOwnerID              []string
	AWSAvailabilityZone     []string
	AWSRegion               []string
	AWSVPCID                []string
	AWSEC2InstanceGroupName []string
	AWSEC2InstanceStateName []string
	AWSEC2InstanceType      []string
	AWSSubnetID             []string
	AWSEC2ProductCode       []string
	AWSEC2Name              []string
	AzureVMID               []string
	AzureResourceID         []string
	SSHFingerPrint          []string
	McafeeEPOGUID           []string
	McafeeEPOAgentGUID      []string
	QualysHostID            []string
	QualysAssetID           []string
	ServiceNowSystemID      []string
	Tags                    []AssetTagDetail
	VulnSevCount            AssetSevSummary
	AuditSevCount           AssetSevSummary
	Interface               []AssetInterface
	Source                  []AssetSource
	Vuln                    []AssetVuln // TODO: Convert to map!
}
type AssetVuln struct {
	PluginID     string
	PluginName   string
	PluginFamily string
	Count        string
	State        string
	Severity     string
	Detail       AssetVulnDetail
	Output       []AssetVulnOutput
}
type AssetVulnDetail struct {
	Description   string
	Solution      string
	Synopsis      string
	Count         string
	Severity      string
	Discovery     interface{}
	PluginDetails interface{}
	ReferenceInfo interface{}
	RiskInfo      interface{}
	SeeAlso       interface{}
	VulnInfo      interface{}
}
type AssetVulnOutput struct {
	PluginOutput string
	States       []struct {
		Name   string
		Result []AssetVulnResult
	}
}
type AssetVulnResult struct {
	ApplicationProtocol string
	TransportProtocol   string
	Port                string
	Severity            string
	Assets              []interface{}
}
type AssetSevSummary struct {
	Total    string
	Severity []AssetSev
}
type AssetSev struct {
	Count string
	Level string
	Name  string
}
type AssetInterface struct {
	Name       string
	IPV4       []string
	IPV6       []string
	FQDN       []string
	MACAddress []string
}
type AssetSource struct {
	FirstSeenAt string
	LastSeenAt  string
	Name        string
}
type AssetTagDetail struct {
	UUID         string
	CategoryName string
	Value        string
	AddedBy      string
	AddedAt      string
	Source       string
}

type TagValue struct {
	ContainerUUID       string
	UUID                string
	ModelName           string
	Value               string
	Description         string
	Type                string
	CategoryUUID        string
	CategoryName        string
	CategoryDescription string
}
type TagCategory struct {
	ContainerUUID string
	UUID          string
	ModelName     string
	Name          string
	Description   string
}

type Scanner struct {
	ID               string
	UUID             string
	Name             string
	Type             string
	Status           string
	ScanCount        string
	EngineVersion    string
	Platform         string
	LoadedPluginSet  string
	RegistrationCode string
	Owner            string
	Key              string
	IP               string
	License          ScannerLicense
	Agents           []ScannerAgent
}
type ScannerLicense struct {
	Type         string
	IPS          string
	Agents       string
	Scanners     string
	AgentsUsed   string
	ScannersUsed string
}
type ScannerAgent struct {
	ID          string
	UUID        string
	Name        string
	Distro      string
	IP          string
	LastScanned time.Time
	Platform    string
	LinkedOn    time.Time
	LastConnect time.Time
	Feed        string
	CoreBuild   string
	CoreVersion string
	Status      string
	Groups      map[string]AgentGroup
	Scanner     Scanner
}
type AgentGroup struct {
	ID          string
	Name        string
	UUID        string
	AgentsCount string
	Agents      []ScannerAgent
}

type VulnExportStatus struct {
	Status          string   `json:"status"`
	Chunks          []string `json:"chunks_available"`
	ChunksFailed    []string `json:"chunks_failed"`
	ChunksCancelled []string `json:"chunks_cancelled"`
}

type AssetExportStatus struct {
	Status          string   `json:"status"`
	Chunks          []string `json:"chunks_available"`
	ChunksFailed    []string `json:"chunks_failed"`
	ChunksCancelled []string `json:"chunks_cancelled"`
}

type VulnExportChunk struct {
	Asset                VulnExportChunkAsset
	Output               string
	Plugin               VulnExportChunkPlugin
	Network              VulnExportChunkPort
	Scan                 VulnExportChunkScan
	Severity             string
	SeverityID           string
	SeverityDefaultID    string
	SeverityModification string
	FirstFound           string
	LastFound            string
	State                string
}
type VulnExportChunkAsset struct {
	DeviceType               string   `json:"device_type"`
	FQDN                     string   `json:"fqdn"`
	HostName                 string   `json:"hostname"`
	UUID                     string   `json:"uuid"`
	IPV4                     string   `json:"ipv4"`
	LastAuthenticatedResults string   `json:"last_unauthenticated_results"`
	NetBIOSWorkgroup         []string `json:"netbios_workgroup"`
	OperatingSystem          []string `json:"operating_system"`
	Tracked                  bool     `json:"tracked"`
}
type VulnExportChunkPlugin struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"decsription"`
	Family           string `json:"family"`
	FamilyID         string `json:"family_id"`
	HasPatch         bool   `json:"has_patch"`
	ModificationDate string `json:"modification_date"`
	PublicationDate  string `json:"publication_date"`
	RiskFactor       string `json:"risk_factor"`
	Solution         string `json:"solution"`
	Synopsis         string `json:"synopsis"`
	Type             string `json:"type"`
	Version          string `json:"version"`
}
type VulnExportChunkPort struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
}
type VulnExportChunkScan struct {
	CompletedAt  string `json:"completed_at"`
	ScheduleUUID string `json:"schedule_uuid"`
	StartedAt    string `json:"started_at"`
	UUID         string `json:"uuid"`
}
