package models

// What we need:
// Something to store user data in
// We may also need a model to display stats between Rhaast and Shadow Assassin

type UserData struct {
	Region   string `json:"region"`
	Username string `json:"username"`
}

type KaynData struct {
	Transformation      string   `json:"transformation"`
	Level               int      `json:"level"`
	Deaths              int      `json:"deaths"`
	FirstBlood          bool     `json:"firstBlood"`
	Gold                int      `json:"gold"`
	Victory             bool     `json:"victory"`
	Items               []string `json:"items"`
	Kills               int      `json:"kills"`
	Pentakills          int      `json:"pentakills"`
	LargestKillingSpree int      `json:"largestKillingSpree"`
	ObjectivesStolen    int      `json:"objectivesStolen"`
	TimePlayed          int      `json:"timePlayed"`
	TotalDamage         int      `json:"totalDamage"`
	TotalDamageTaken    int      `json:"totalDamageTaken"`
	DoubleKills         int      `json:"doubleKills"`
	TripleKills         int      `json:"tripleKills"`
	QuadraKills         int      `json:"quadraKills"`
	TrueDamage          int      `json:"trueDamage"`
	VisionScore         int      `json:"visionScore"`
	BaronKills          int      `json:"baronKills"`
	DragonKills         int      `json:"dragonKills"`
	Assists             int      `json:"assists"`
	Ratio               int      `json:"ratio"`
}
