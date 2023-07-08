package models

// What we need:
// Something to store user data in
// We may also need a model to display stats between Rhaast and Shadow Assassin

type UserData struct {
	Region   string `json:"region"`
	Username string `json:"username"`
}

type RhaastData struct {
}

type ShadowAssassinData struct {
}
