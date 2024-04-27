package service

type RankingSummary struct {
	ItemId       int
	UsingRate    float64
	TotalPreset  int
	TotalAccount int
	TotalEnchant int
	Enchants     map[string]float64
}

// userId -> jobId -> skillName -> itemPosition -> { "itemId": 3, enchants: { "0-0-0": 999 } }
type ItemPositionSummary = map[string]map[int]*ItemSummary
type UsingSkillSummary = map[string]ItemPositionSummary
type JobSummary = map[int]UsingSkillSummary
type AllSummary = map[string]JobSummary

type PresetSummary struct {
	SummaryClassSkillMap map[int]map[string]int
	TotalSelectedJobMap  map[int]int
	JobSummary           map[int]map[string]map[string][]RankingSummary
}

type PresetSummaryService interface {
	GenerateSummary() (*PresetSummary, error)
}
