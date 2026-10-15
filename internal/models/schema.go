package models

// MonsterStats def D&D struct for Entity
type MonsterStats struct {
	ArmorClass          int               `json:"armor_class"`
	HitPoints           string            `json:"hit_points"`
	Speed               string            `json:"speed"`
	Attributes          map[string]int    `json:"attributes"`
	SavingThrows        []string          `json:"saving_throws"`
	Skills              []string          `json:"skills"`
	DamageResistances   []string          `json:"damage_resistances"`
	DamageImmunities    []string          `json:"damage_immunities"`
	ConditionImmunities []string          `json:"condition_immunities"`
	Senses              []string          `json:"senses"`
	Languages           []string          `json:"languages"`
	ChallengeRating     string            `json:"challenge_rating"`
	Actions             []Action          `json:"actions,omitempty"`
	LegendaryActions    []LegendaryAction `json:"legendary_actions,omitempty"`
}

type Action struct {
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Hit    string `json:"hit,omitempty"`
	Reach  string `json:"reach,omitempty"`
	Damage string `json:"damage,omitempty"`
	Effect string `json:"effect,omitempty"`
}

type LegendaryAction struct {
	Name   string `json:"name"`
	Effect string `json:"effect"`
}

// ItemProperties def dinamic attributes for Item JSON.
type ItemProperties struct {
	WeaponType         string          `json:"weapon_type,omitempty"`
	Subtype            string          `json:"subtype,omitempty"`
	Rarity             string          `json:"rarity,omitempty"`
	AttunementRequired bool            `json:"attunement_required,omitempty"`
	Damage             string          `json:"damage,omitempty"`
	DamageType         string          `json:"damage_type,omitempty"`
	Weight             string          `json:"weight,omitempty"`
	MagicBonus         string          `json:"magic_bonus,omitempty"`
	SpecialEffects     []SpecialEffect `json:"special_effects,omitempty"`
	Curse              string          `json:"curse,omitempty"`
}

type SpecialEffect struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
