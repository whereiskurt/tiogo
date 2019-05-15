package metrics

type EndPointType string

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
}

func (c EndPointType) String() string {
	return "api." + string(c)
}

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

type MethodTypes struct {
	Service   ServiceTypes
	DB        DbTypes
	Cache     CacheTypes
	Transport TransportTypes
}

type ServiceMethodType string
type ServiceTypes struct {
	Get    ServiceMethodType
	Update ServiceMethodType
	Add    ServiceMethodType
	Delete ServiceMethodType
}

func (c ServiceMethodType) String() string {
	return "service." + string(c)
}

type DbMethodType string
type DbTypes struct {
	Read   DbMethodType
	Update DbMethodType
	Insert DbMethodType
	Delete DbMethodType
}

func (c DbMethodType) String() string {
	return "db." + string(c)
}

type CacheMethodType string
type CacheTypes struct {
	Hit        CacheMethodType
	Miss       CacheMethodType
	Store      CacheMethodType
	Invalidate CacheMethodType
}

func (c CacheMethodType) String() string {
	return "cache." + string(c)
}

type TransportMethodType string
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
