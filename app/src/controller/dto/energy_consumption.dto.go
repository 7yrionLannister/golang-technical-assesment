package dto

type EnergyConsumptionDTO struct {
	MeterId uint      `json:"meter_id"`
	Address string    `json:"address"` // TODO consume external API
	Active  []float64 `json:"active"`
	// ReactiveInductive  []float64 `json:"reactive_inductive"`
	// ReactiveCapacitive []float64 `json:"reactive_capacitive"`
	// Exported           []float64 `json:"exported"`
}
