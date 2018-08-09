package resources

const (
	Prefix                  = "remesh-"
	GatewaySuffix           = "-gateway"
	VsSuffix                = "-vs"
	DsSuffix                = "-ds"
	AutoGeneratedLabelName  = "remesh.bevyx.com/auto-generated"
	AutoGeneratedLabelValue = "true"
)

var AutoGeneratedLabels = map[string]string{
	AutoGeneratedLabelName: AutoGeneratedLabelValue,
}

func GetSubsetName(host string, subsetHash string) string {
	return host + "-" + subsetHash
}

const (
	HostAll = "*"
	Mesh    = "mesh"
)
