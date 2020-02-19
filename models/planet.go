package models

type Planet struct {
	ID      string `bson:"_id" json:"id"`
	Name    string `bson:"name" json:"name"`
	Climate string `bson:"climate" json:"climate"`
	Terrain string `bson:"terrain" json:"terrain"`
	Films   int    `bson:"films" json:"films"`
}

type PlanetDocument struct {
	Name    string
	Climate string
	Terrain string
	Films   int
}
