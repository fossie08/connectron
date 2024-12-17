package types

type GameWindow struct {
	GridWidth   int
	GridHeight  int
	LineLength  int
	PlayerCount int
	IncludeAI   bool
	Alliances   bool
}