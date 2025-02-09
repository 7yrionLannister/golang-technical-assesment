package view

type EnergyConsumptionDTO struct {
	MeterId          uint
	Address          string // TODO consume external API
	TotalConsumption float64
}
