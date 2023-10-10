package main

import "encoding/xml"

type CraftDataDocument struct {
	XMLName  xml.Name    `xml:"ArrayOfCraftData"`
	Data     []CraftData `xml:"CraftData"`
	FileName string      `xml:"-"`
}

type CraftData struct {
	XMLName              xml.Name    `xml:"CraftData"`
	Machine              string      `xml:"-"`
	Key                  string      `xml:"Key"`
	Name                 string      `xml:"Name"`
	Category             string      `xml:"Category"`
	Tier                 int         `xml:"Tier"`
	CraftedKey           string      `xml:"CraftedKey"`
	CraftedAmount        int         `xml:"CraftedAmount"`
	Costs                []CraftCost `xml:"Costs>CraftCost"`
	Description          string      `xml:"Description"`
	Hint                 string      `xml:"Hint"`
	ResearchCost         int         `xml:"ResearchCost"`
	ScanRequirements     []string    `xml:"ScanRequirements>Scan"`
	MoreScanRequirements []string    `xml:"ScanRequirements>scan"`
	ResearchRequirements []string    `xml:"ResearchRequirements>Research"`
	CanCraftAnywhere     bool        `xml:"CanCraftAnywhere"`
	MasterRecipe         string      `xml:"MasterRecipe"`
	RequiredModule       string      `xml:"RequiredModule"`
	Research             string      `xml:"Research"`
}

type CraftCost struct {
	XMLName xml.Name `xml:"CraftCost"`
	Key     string   `xml:"Key"`
	Name    string   `xml:"Name"`
	Amount  int      `xml:"Amount"`
}

type GenericAutoCrafterDataEntry struct {
	XMLName                xml.Name `xml:"GenericAutoCrafterDataEntry"`
	FriendlyName           string   `xml:"FriendlyName"`
	SpawnObject            string   `xml:"SpawnObject"`
	Value                  string   `xml:"Value"`
	OffloadTarget          string   `xml:"OffloadTarget"`
	CraftingString         string   `xml:"CraftingString"`
	PowerUsePerSecond      float32  `xml:"PowerUsePerSecond"`
	PowerTransferPerSecond int      `xml:"PowerTransferPerSecond"`
	MaxPowerStorage        int      `xml:"MaxPowerStorage"`
	WorkingAnimation       string   `xml:"WorkingAnimation"`
	UnloadingAnimation     string   `xml:"UnloadingAnimation"`
	LoadingAnimation       string   `xml:"LoadingAnimation"`
	CraftTime              int      `xml:"CraftTime"`
	OptionalIngredients    bool     `xml:"OptionalIngredients"`
	AnimScalar             float32  `xml:"AnimScalar"`
	Recipe                 Recipe   `xml:"Recipe"`
}

type Recipe struct {
	XMLName       xml.Name    `xml:"Recipe"`
	Key           string      `xml:"Key"`
	CraftedKey    string      `xml:"CraftedKey"`
	CraftedName   string      `xml:"CraftedName"`
	CraftedAmount int         `xml:"CraftedAmount"`
	Costs         []CraftCost `xml:"Costs>CraftCost"`
	Description   string      `xml:"Description"`
}

type ItemEntryDocument struct {
	XMLName     xml.Name    `xml:"ArrayOfItemEntry"`
	ItemEntries []ItemEntry `xml:"ItemEntry"`
}

type ItemEntry struct {
	XMLName              xml.Name `xml:"ItemEntry"`
	ItemID               int      `xml:"ItemId"`
	Key                  string   `xml:"Key"`
	Name                 string   `xml:"Name"`
	Plural               string   `xml:"Plural"`
	Type                 string   `xml:"Type"`
	Object               string   `xml:"Object"`
	Sprite               string   `xml:"Sprite"`
	Category             string   `xml:"Category"`
	ScanRequirements     []string `xml:"ScanRequirements>Scan"`
	ResearchRequirements []string `xml:"ResearchRequirements>Research"`
	ItemAction           string   `xml:"ItemAction"`
	Hidden               bool     `xml:"Hidden"`
	DecomposeValue       int      `xml:"DecomposeValue"`
	MaxDurability        int      `xml:"MaxDurability"`
	ActionParameter      string   `xml:"ActionParameter"`
}

type RecipeSetDocument struct {
	XMLName    xml.Name    `xml:"ArrayOfRecipeSet"`
	RecipeSets []RecipeSet `xml:"RecipeSet"`
}

type RecipeSet struct {
	XMLName    xml.Name `xml:"RecipeSet"`
	Id         string   `xml:"Id"`
	Name       string   `xml:"Name"`
	FileName   string   `xml:"FileName"`
	Icon       string   `xml:"Icon"`
	MachineKey string   `xml:"MachineKey"`
}

type ResearchDataEntryDocument struct {
	XMLName             xml.Name            `xml:"ArrayOfResearchDataEntry"`
	ResearchDataEntries []ResearchDataEntry `xml:"ResearchDataEntry"`
}

type ResearchDataEntry struct {
	XMLName                  xml.Name      `xml:"ResearchDataEntry"`
	Key                      string        `xml:"Key"`
	Name                     string        `xml:"Name"`
	IconName                 string        `xml:"IconName"`
	ResearchCost             int           `xml:"ResearchCost"`
	PreDescription           string        `xml:"PreDescription"`
	PostDescription          string        `xml:"PostDescription"`
	ScanRequirements         []string      `xml:"ScanRequirements>Scan"`
	ResearchRequirements     []string      `xml:"ResearchRequirements>Research"`
	MoreResearchRequirements []string      `xml:"ResearchRequirements>Key"`
	ProjectItemRequirements  []Requirement `xml:"Requirement"`
}

type Requirement struct {
	XMLName xml.Name `xml:"Requirement"`
	Key     string   `xml:"Key"`
	Amount  int      `xml:"Amount"`
}

type TerrainDataEntryDocument struct {
	XMLName            xml.Name           `xml:"ArrayOfTerrainDataEntry"`
	TerrainDataEntries []TerrainDataEntry `xml:"TerrainDataEntry"`
}

type TerrainDataEntry struct {
	XMLName             xml.Name     `xml:"TerrainDataEntry"`
	CubeType            int          `xml:"CubeType"`
	Key                 string       `xml:"Key"`
	Name                string       `xml:"Name"`
	Description         string       `xml:"Description"`
	LayerType           string       `xml:"LayerType"`
	TopTexture          int          `xml:"TopTexture"`
	SideTexture         int          `xml:"SideTexture"`
	BottomTexture       int          `xml:"BottomTexture"`
	GUITexture          int          `xml:"GUITexture"`
	IconName            string       `xml:"IconName"`
	Hidden              bool         `xml:"Hidden"`
	IsSolid             bool         `xml:"isSolid"`
	IsTransparent       bool         `xml:"isTransparent"`
	IsHollow            bool         `xml:"isHollow"`
	IsGlass             bool         `xml:"isGlass"`
	IsPassable          bool         `xml:"isPassable"`
	IsReinforced        bool         `xml:"isReinforced"`
	IsGarbage           bool         `xml:"isGarbage"`
	IsColorised         bool         `xml:"isColorised"`
	IsMultiBlockMachine bool         `xml:"isMultiBlockMachine"`
	IsPaintable         bool         `xml:"isPaintable"`
	IsHidden            bool         `xml:"isHidden"`
	PickReplacement     string       `xml:"PickReplacement"`
	HasObject           bool         `xml:"hasObject"`
	HasEntity           bool         `xml:"hasEntity"`
	AudioWalkType       string       `xml:"AudioWalkType"`
	AudioBuildType      string       `xml:"AudioBuildType"`
	AudioDestroyType    string       `xml:"AudioDestroyType"`
	Tags                []string     `xml:"tags>tag"`
	Category            string       `xml:"Category"`
	Hardness            int          `xml:"Hardness"`
	DefaultValue        int          `xml:"DefaultValue"`
	Value               int          `xml:"Value"`
	Values              []ValueEntry `xml:"Values>ValueEntry"`
	Stages              []Stage      `xml:"Stages>Stage"`
	MaxStack            int          `xml:"MaxStack"`
	OreType             string       `xml:"OreType"`
	Fuel                Fuel         `xml:"Fuel"`
	ResearchValue       int          `xml:"ResearchValue"`
	ModStartValue       int          `xml:"ModStartValue"`
	TypeDroppedOnDig    string       `xml:"TypeDroppedOnDig"`
	Plural              string       `xml:"Plural"`
}

type ValueEntry struct {
	XMLName          xml.Name `xml:"ValueEntry"`
	Value            int      `xml:"Value"`
	Key              string   `xml:"Key"`
	Name             string   `xml:"Name"`
	IconName         string   `xml:"IconName"`
	Description      string   `xml:"Description"`
	ResearchValue    int      `xml:"ResearchValue"`
	TypeDroppedOnDig string   `xml:"TypeDroppedOnDig"`
	Hardness         int      `xml:"Hardness"`
}

type Stage struct {
	XMLName       xml.Name `xml:"Stage"`
	RangeMinimum  int      `xml:"RangeMinimum"`
	RangeMaximum  int      `xml:"RangeMaximum"`
	TopTexture    int      `xml:"TopTexture"`
	SideTexture   int      `xml:"SideTexture"`
	BottomTexture int      `xml:"BottomTexture"`
	GUITexture    int      `xml:"GUITexture"`
	IconName      string   `xml:"IconName"`
}

type Fuel struct {
	XMLName          xml.Name      `xml:"Fuel"`
	Energy           int           `xml:"Energy"`
	ParticleEffect   string        `xml:"ParticleEffect"`
	ParticleLifetime float32       `xml:"ParticleLifetime"`
	ParticleColor    ParticleColor `xml:"ParticleColor"`
}

type ParticleColor struct {
	XMLName xml.Name `xml:"ParticleColor"`
	R       int      `xml:"R"`
	G       int      `xml:"G"`
	B       int      `xml:"B"`
	A       int      `xml:"A"`
}
