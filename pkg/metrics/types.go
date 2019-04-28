package metrics

type EndPointType string

var EndPoints = endPointTypes{
	ScannersList:          EndPointType("ScannersList"),
	AgentsList:            EndPointType("AgentsList"),
	AgentGroups:       EndPointType("AgentGroups"),
	VulnsExportStart:  EndPointType("VulnsExportStart"),
	VulnsExportStatus: EndPointType("VulnsExportStatus"),
	VulnsExportGet:    EndPointType("VulnsExportGet"),
	VulnsExportQuery:  EndPointType("VulnsExportQuery"),
}

type endPointTypes struct {
	ScannersList          EndPointType
	AgentsList            EndPointType
	AgentGroups       EndPointType
	VulnsExportStart  EndPointType
	VulnsExportStatus EndPointType
	VulnsExportGet    EndPointType
	VulnsExportQuery  EndPointType
}

func (c EndPointType) String() string {
	return "api." + string(c)
}

var Methods = methodTypes{
	Service: serviceTypes{
		Get:    serviceMethodType("Get"),
		Update: serviceMethodType("Update"),
		Add:    serviceMethodType("Add"),
		Delete: serviceMethodType("Delete"),
	},
	DB: dbTypes{
		Delete: dbMethodType("Delete"),
		Update: dbMethodType("Update"),
		Read:   dbMethodType("Read"),
		Insert: dbMethodType("Insert"),
	},
	Cache: cacheTypes{
		Hit:        cacheMethodType("Hit"),
		Miss:       cacheMethodType("Miss"),
		Invalidate: cacheMethodType("Invalidate"),
		Store:      cacheMethodType("Store"),
	},
	Transport: transportTypes{
		Put:    transportMethodType("Put"),
		Delete: transportMethodType("Delete"),
		Post:   transportMethodType("Post"),
		Get:    transportMethodType("Get"),
		Head:   transportMethodType("Head"),
	},
}

type methodTypes struct {
	Service   serviceTypes
	DB        dbTypes
	Cache     cacheTypes
	Transport transportTypes
}

type serviceMethodType string
type serviceTypes struct {
	Get    serviceMethodType
	Update serviceMethodType
	Add    serviceMethodType
	Delete serviceMethodType
}

func (c serviceMethodType) String() string {
	return "service." + string(c)
}

type dbMethodType string
type dbTypes struct {
	Read   dbMethodType
	Update dbMethodType
	Insert dbMethodType
	Delete dbMethodType
}

func (c dbMethodType) String() string {
	return "db." + string(c)
}

type cacheMethodType string
type cacheTypes struct {
	Hit        cacheMethodType
	Miss       cacheMethodType
	Store      cacheMethodType
	Invalidate cacheMethodType
}

func (c cacheMethodType) String() string {
	return "cache." + string(c)
}

type transportMethodType string
type transportTypes struct {
	Get    transportMethodType
	Put    transportMethodType
	Post   transportMethodType
	Delete transportMethodType
	Head   transportMethodType
}

func (c transportMethodType) String() string {
	return "http." + string(c)
}
