package metrics

// EndPointType creates a type for the map lookups
type EndPointType string

// EndPoints are each type of metric we will produce
var EndPoints = endPointTypes{
	ScannersList:       EndPointType("ScannersList"),
	AgentsList:         EndPointType("AgentsList"),
	AgentGroups:        EndPointType("AgentGroups"),
	AgentsGroup:        EndPointType("AgentsGroup"),
	AgentsUngroup:      EndPointType("AgentsUngroup"),
	VulnsExportStart:   EndPointType("VulnsExportStart"),
	VulnsExportStatus:  EndPointType("VulnsExportStatus"),
	VulnsExportGet:     EndPointType("VulnsExportGet"),
	VulnsExportQuery:   EndPointType("VulnsExportQuery"),
	AssetsExportStart:  EndPointType("AssetsExportStart"),
	AssetsExportStatus: EndPointType("AssetsExportStatus"),
	AssetsExportGet:    EndPointType("AssetsExportGet"),
	AssetsExportQuery:  EndPointType("AssetsExportQuery"),
	ScansList:          EndPointType("ScansList"),
}

type endPointTypes struct {
	ScannersList       EndPointType
	AgentsList         EndPointType
	AgentGroups        EndPointType
	AgentsGroup        EndPointType
	AgentsUngroup      EndPointType
	VulnsExportStart   EndPointType
	VulnsExportStatus  EndPointType
	VulnsExportGet     EndPointType
	VulnsExportQuery   EndPointType
	AssetsExportStart  EndPointType
	AssetsExportStatus EndPointType
	AssetsExportGet    EndPointType
	AssetsExportQuery  EndPointType
	ScansList          EndPointType
}

func (c EndPointType) String() string {
	return "api." + string(c)
}

// Methods types for metrics
var Methods = MethodTypes{
	Service: ServiceTypes{
		Get:    ServiceMethodType("Get"),
		Update: ServiceMethodType("Update"),
		Add:    ServiceMethodType("Add"),
		Delete: ServiceMethodType("Delete"),
	},
	DB: DbTypes{
		Delete: DbMethodType("Delete"),
		Update: DbMethodType("Update"),
		Read:   DbMethodType("Read"),
		Insert: DbMethodType("Insert"),
	},
	Cache: CacheTypes{
		Hit:        CacheMethodType("Hit"),
		Miss:       CacheMethodType("Miss"),
		Invalidate: CacheMethodType("Invalidate"),
		Store:      CacheMethodType("Store"),
	},
	Transport: TransportTypes{
		Put:    TransportMethodType("Put"),
		Delete: TransportMethodType("Delete"),
		Post:   TransportMethodType("Post"),
		Get:    TransportMethodType("Get"),
		Head:   TransportMethodType("Head"),
	},
}

// MethodTypes are all of the possible method types to metric on - Serveice,DB, Cache, Transport...
type MethodTypes struct {
	Service   ServiceTypes
	DB        DbTypes
	Cache     CacheTypes
	Transport TransportTypes
}

// ServiceMethodType wrapper from string for type safety
type ServiceMethodType string

// ServiceTypes we can have add/get/update/delete
type ServiceTypes struct {
	Get    ServiceMethodType
	Update ServiceMethodType
	Add    ServiceMethodType
	Delete ServiceMethodType
}

func (c ServiceMethodType) String() string {
	return "service." + string(c)
}

// DbMethodType wrapper from string for type safety
type DbMethodType string

// DbTypes we can have insert/read/update/delete
type DbTypes struct {
	Read   DbMethodType
	Update DbMethodType
	Insert DbMethodType
	Delete DbMethodType
}

func (c DbMethodType) String() string {
	return "db." + string(c)
}

// CacheMethodType wrapper from string for type safety
type CacheMethodType string

// CacheTypes we can have hit/miss/store/invalidate
type CacheTypes struct {
	Hit        CacheMethodType
	Miss       CacheMethodType
	Store      CacheMethodType
	Invalidate CacheMethodType
}

func (c CacheMethodType) String() string {
	return "cache." + string(c)
}

// TransportMethodType wrapper from string for type safety
type TransportMethodType string

// TransportTypes we can have get/put/post/delete/head
type TransportTypes struct {
	Get    TransportMethodType
	Put    TransportMethodType
	Post   TransportMethodType
	Delete TransportMethodType
	Head   TransportMethodType
}

func (c TransportMethodType) String() string {
	return "http." + string(c)
}
