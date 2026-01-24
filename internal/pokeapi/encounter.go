package pokeapi

type Encounter struct {
	Method          NamedAPIResource   `json:"method"`
	ConditionValues []NamedAPIResource `json:"condition_values"`
	Chance          int                `json:"chance"`
	MinLevel        int                `json:"min_level"`
	MaxLevel        int                `json:"max_level"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource          `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetails `json:"version_details"`
}

type EncounterVersionDetails struct {
	Version NamedAPIResource `json:"version"`
	Rate    int              `json:"rate"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource         `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

type VersionEncounterDetail struct {
	Version          NamedAPIResource `json:"version"`
	EncounterDetails []Encounter      `json:"encounter_details"`
	MaxChance        int              `json:"max_chance"`
}
