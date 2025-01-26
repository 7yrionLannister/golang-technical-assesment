package dto

type PeriodicConsumptionDTO struct {
	Period    []string                `json:"period"`
	DataGraph []*EnergyConsumptionDTO `json:"data_graph"`
}
