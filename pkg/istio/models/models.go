package models

//TransformedService is
type TransformedService struct {
	Host              string
	ServiceSubsetList []ServiceSubset
}

//ServiceSubset is
type ServiceSubset struct {
	Labels              map[string]string
	SubsetHash          string
	VirtualEnvironments []string
}
