package models

type TransformedService struct {
	Host              string
	ServiceSubsetList []ServiceSubset
}

type ServiceSubset struct {
	Labels              map[string]string
	SubsetHash          string
	VirtualEnvironments []string
}
