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

// ScanHistoryItem are all the instances for this scan
type ScanHistoryItem struct {
	HistoryID        string
	UUID             string
	ScanType         string
	Status           string
	LastModifiedDate string
	CreationDate     string
}

// ScanHistoryDetail is a particular scan instance
type ScanHistoryDetail struct {
	Scan                *Scan
	History             []ScanHistoryItem
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
	CompliancePlugin    map[string]Plugin
	VulnPlugin          map[string]Plugin
	Host                map[string]HostScanSummary
	HostAssetMap        map[string]string
}

// HostScanSummary struct
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

// HostScanDetail struct
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

// Plugin struct
type Plugin struct {
	PluginID   string
	Name       string
	FamilyName string
	FamilyID   string
	Count      string
	Severity   string
	Detail     PluginDetail
}

// PluginDetail struct
type PluginDetail struct {
	FunctionName          string
	PluginPublicationDate string
	PatchPublicationDate  string
	Attribute             map[string]PluginDetailAttribute
}

// PluginDetailAttribute struct
type PluginDetailAttribute struct {
	Name  string
	Value string
}

// PluginFamily struct
type PluginFamily struct {
	ID    string
	Name  string
	Count string
}

// PluginFamilyDetail struct
type PluginFamilyDetail struct {
	ID       string
	Name     string
	PluginID []string
}

// Asset struct
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

// AssetVuln struct
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

// AssetVulnDetail struct
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

// AssetVulnOutput struct
type AssetVulnOutput struct {
	PluginOutput string
	States       []struct {
		Name   string
		Result []AssetVulnResult
	}
}

// AssetVulnResult struct
type AssetVulnResult struct {
	ApplicationProtocol string
	TransportProtocol   string
	Port                string
	Severity            string
	Assets              []interface{}
}

// AssetSevSummary struct
type AssetSevSummary struct {
	Total    string
	Severity []AssetSev
}

// AssetSev struct
type AssetSev struct {
	Count string
	Level string
	Name  string
}

// AssetInterface struct
type AssetInterface struct {
	Name       string
	IPV4       []string
	IPV6       []string
	FQDN       []string
	MACAddress []string
}

// AssetSource struct
type AssetSource struct {
	FirstSeenAt string
	LastSeenAt  string
	Name        string
}

// AssetTagDetail struct
type AssetTagDetail struct {
	UUID         string
	CategoryName string
	Value        string
	AddedBy      string
	AddedAt      string
	Source       string
}

// TagValue struct
type TagValue struct {
	UUID                string
	CategoryUUID        string
	CategoryName        string
	CategoryDescription string
	Value               string
	Description         string
}

//TagCategory struct
type TagCategory struct {
	ContainerUUID string
	UUID          string
	ModelName     string
	Name          string
	Description   string
}

// Scanner struct
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

// ScannerLicense struct
type ScannerLicense struct {
	Type         string
	IPS          string
	Agents       string
	Scanners     string
	AgentsUsed   string
	ScannersUsed string
}

// ScannerAgent struct
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

// AgentGroup struct
type AgentGroup struct {
	ID          string
	Name        string
	UUID        string
	AgentsCount string
	Agents      []ScannerAgent
}

// VulnExportStatus struct
type VulnExportStatus struct {
	Status          string
	Chunks          []string
	ChunksFailed    []string
	ChunksCancelled []string
}

// AssetExportStatus struct
type AssetExportStatus struct {
	Status          string
	Chunks          []string
	ChunksFailed    []string
	ChunksCancelled []string
}

// VulnExportChunk struct
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

// VulnExportChunkAsset struct
type VulnExportChunkAsset struct {
	DeviceType               string
	FQDN                     string
	HostName                 string
	UUID                     string
	IPV4                     string
	LastAuthenticatedResults string
	NetBIOSWorkgroup         []string
	OperatingSystem          []string
	Tracked                  bool
}

// VulnExportChunkPlugin struct
type VulnExportChunkPlugin struct {
	ID               string
	Name             string
	Description      string
	Family           string
	FamilyID         string
	HasPatch         bool
	ModificationDate string
	PublicationDate  string
	RiskFactor       string
	Solution         string
	Synopsis         string
	Type             string
	Version          string
}

// VulnExportChunkPort struct
type VulnExportChunkPort struct {
	Port     string
	Protocol string
	Service  string
}

// VulnExportChunkScan struct
type VulnExportChunkScan struct {
	CompletedAt  string
	ScheduleUUID string
	StartedAt    string
	UUID         string
}

// ScansExportStart is outputed at successful scans export
type ScansExportStart struct {
	FileUUID  string
	TempToken string
	Format    string
}

// ScansExportStatus is outputed at successful scans export
type ScansExportStatus struct {
	Status   string
	FileUUID string
}

// ScansExportGet is outputed at successful scans export
type ScansExportGet struct {
	ScanID       string
	HistoryID    string
	ScheduleUUID string

	Policy struct {
		Name        string
		Comments    string
		Preferences struct {
			Server  map[string]string
			Plugins []PolicyPreferencePlugin
		}
		FamilyStatus map[string]string
		Plugins      []PolicyPlugin
	}
	Report struct {
		Name  string
		Hosts []ReportHost
	}
	SourceFile struct {
		FileUUID       string
		CachedFileName string
	}
}

// PolicyPreferencePlugin used with export-scans nessus xml
type PolicyPreferencePlugin struct {
	PluginName       string
	PluginID         string
	FullName         string
	PreferenceName   string
	PreferenceType   string
	PreferenceValues string
	SelectedValue    string
}

// PolicyPlugin used with export-scans nessus xml
type PolicyPlugin struct {
	PluginID   string
	PluginName string
	Family     string
	Status     string
}

// ReportHost used with export-scans nessus xml
type ReportHost struct {
	Name       string
	HostTag    map[string]string
	ReportItem []ReportItemType
}

// ReportItemType used with export-scans nessus xml
type ReportItemType struct {
	BID                        []string
	CanvasPackage              string
	CVE                        []string
	CVSS3BaseScore             string
	CVSS3TemporalScore         string
	CVSS3TemporalVector        string
	CVSS3Vector                string
	CVSSBaseScore              string
	CVSSTemporalScore          string
	CVSSTemporalVector         string
	CVSSVector                 string
	Description                string
	ExploitAvailable           string
	ExploitedByMalware         string
	ExploitFrameworkCanvas     string
	ExploitFrameworkCore       string
	ExploitFrameworkMetasploit string
	InTheNews                  string
	MetasploitName             string
	PatchPublicationDate       string
	PluginFamily               string
	PluginID                   string
	PluginModificationDate     string
	PluginName                 string
	PluginOutput               string
	PluginPublicationDate      string
	PluginType                 string
	Port                       string
	Protocol                   string
	RiskFactor                 string
	ScriptVersion              string
	SeeAlso                    string
	Severity                   string
	Solution                   string
	SvcName                    string
	Synopsis                   string
	UnsupportedByVendor        string
	VulnPublicationDate        string
	XRef                       []string
}
