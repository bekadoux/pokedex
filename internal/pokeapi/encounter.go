package pokeapi

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionDetails struct {
	EncounterDetails []EncounterDetails `json:"encounter_details"`
	Rate             int                `json:"rate"`
	MaxChance        int                `json:"max_chance"`
	Version          Version            `json:"version"`
}

type EncounterDetails struct {
	Chance          int             `json:"chance"`
	ConditionValues []any           `json:"condition_values"`
	MaxLevel        int             `json:"max_level"`
	Method          EncounterMethod `json:"method"`
	MinLevel        int             `json:"min_level"`
}

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterMethodRate struct {
	EncounterMethod EncounterMethod  `json:"encounter_method"`
	VersionDetails  []VersionDetails `json:"version_details"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}
