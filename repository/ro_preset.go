package repository

import "time"

type PresetModel struct {
	Class                 int            `bson:"class" json:"class"`
	Level                 int            `bson:"level" json:"level"`
	JobLevel              int            `bson:"jobLevel" json:"jobLevel"`
	Str                   int            `bson:"str" json:"str"`
	JobStr                int            `bson:"jobStr" json:"jobStr"`
	Agi                   int            `bson:"agi" json:"agi"`
	JobAgi                int            `bson:"jobAgi" json:"jobAgi"`
	Vit                   int            `bson:"vit" json:"vit"`
	JobVit                int            `bson:"jobVit" json:"jobVit"`
	Int                   int            `bson:"int" json:"int"`
	JobInt                int            `bson:"jobInt" json:"jobInt"`
	Dex                   int            `bson:"dex" json:"dex"`
	JobDex                int            `bson:"jobDex" json:"jobDex"`
	Luk                   int            `bson:"luk" json:"luk"`
	JobLuk                int            `bson:"jobLuk" json:"jobLuk"`
	SelectedAtkSkill      string         `bson:"selectedAtkSkill" json:"selectedAtkSkill"`
	RawOptionTxts         []interface{}  `bson:"rawOptionTxts" json:"rawOptionTxts"`
	WeaponRefine          int            `bson:"weaponRefine" json:"weaponRefine"`
	HeadUpperRefine       int            `bson:"headUpperRefine" json:"headUpperRefine"`
	Armor                 int            `bson:"armor" json:"armor"`
	ArmorRefine           int            `bson:"armorRefine" json:"armorRefine"`
	ShieldRefine          int            `bson:"shieldRefine" json:"shieldRefine"`
	Garment               int            `bson:"garment" json:"garment"`
	GarmentRefine         int            `bson:"garmentRefine" json:"garmentRefine"`
	Boot                  int            `bson:"boot" json:"boot"`
	BootRefine            int            `bson:"bootRefine" json:"bootRefine"`
	AccRight              int            `bson:"accRight" json:"accRight"`
	CostumeEnchantUpper   int            `bson:"costumeEnchantUpper" json:"costumeEnchantUpper"`
	CostumeEnchantMiddle  int            `bson:"costumeEnchantMiddle" json:"costumeEnchantMiddle"`
	CostumeEnchantLower   int            `bson:"costumeEnchantLower" json:"costumeEnchantLower"`
	CostumeEnchantGarment int            `bson:"costumeEnchantGarment" json:"costumeEnchantGarment"`
	ShadowWeapon          int            `bson:"shadowWeapon" json:"shadowWeapon"`
	ShadowWeaponRefine    int            `bson:"shadowWeaponRefine" json:"shadowWeaponRefine"`
	ShadowArmor           int            `bson:"shadowArmor" json:"shadowArmor"`
	ShadowArmorRefine     int            `bson:"shadowArmorRefine" json:"shadowArmorRefine"`
	ShadowShield          int            `bson:"shadowShield" json:"shadowShield"`
	ShadowShieldRefine    int            `bson:"shadowShieldRefine" json:"shadowShieldRefine"`
	ShadowBoot            int            `bson:"shadowBoot" json:"shadowBoot"`
	ShadowBootRefine      int            `bson:"shadowBootRefine" json:"shadowBootRefine"`
	ShadowEarring         int            `bson:"shadowEarring" json:"shadowEarring"`
	ShadowEarringRefine   int            `bson:"shadowEarringRefine" json:"shadowEarringRefine"`
	ShadowPendant         int            `bson:"shadowPendant" json:"shadowPendant"`
	ShadowPendantRefine   int            `bson:"shadowPendantRefine" json:"shadowPendantRefine"`
	SkillBuffMap          map[string]int `bson:"SkillBuffMap" json:"SkillBuffMap"`
	SkillBuffs            []int          `bson:"skillBuffs" json:"skillBuffs"`
	ActiveSkillMap        map[string]int `bson:"ActiveSkillMap" json:"ActiveSkillMap"`
	ActiveSkills          []int          `bson:"activeSkills" json:"activeSkills"`
	PassiveSkillMap       map[string]int `bson:"PassiveSkillMap" json:"PassiveSkillMap"`
	PassiveSkills         []int          `bson:"passiveSkills" json:"passiveSkills"`
	ConsumableMap         map[string]int `bson:"ConsumableMap" json:"ConsumableMap"`
	Consumables           []interface{}  `bson:"consumables" json:"consumables"`
	Consumable2Map        map[string]int `bson:"Consumable2Map" json:"Consumable2Map"`
	Consumables2          []interface{}  `bson:"consumables2" json:"consumables2"`
	AspdPotionMap         map[string]int `bson:"AspdPotionMap" json:"AspdPotionMap"`
	AspdPotions           []interface{}  `bson:"aspdPotions" json:"aspdPotions"`
}

type RoPreset struct {
	Id          string      `bson:"id" json:"id"`
	UserId      string      `bson:"user_id" json:"userId"`
	Name        string      `bson:"name" json:"name"`
	Label       string      `bson:"label" json:"label"`
	Model       PresetModel `bson:"model" json:"model"`
	ClassId     int         `bson:"class_id" json:"classId"`
	CreatedAt   time.Time   `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updatedAt"`
	PublishName string      `bson:"publish_name" json:"publishName"`
	IsPublished bool        `bson:"is_published" json:"isPublished"`
	PublishedAt time.Time   `bson:"published_at" json:"publishedAt"`
}

type CreatePresetInput struct {
	UserId string      `bson:"user_id" json:"userId"`
	Label  string      `bson:"label" json:"label"`
	Model  PresetModel `bson:"model" json:"model"`
}

type UpdatePresetInput struct {
	ClassId     int          `bson:"class_id,omitempty" json:"classId"`
	UserId      string       `bson:"user_id,omitempty" json:"userId"`
	Label       string       `bson:"label,omitempty" json:"label"`
	Model       *PresetModel `bson:"model,omitempty" json:"model"`
	UpdatedAt   time.Time    `bson:"updated_at,omitempty" json:"updatedAt"`
	PublishName string       `bson:"publish_name,omitempty" json:"publishName"`
	IsPublished bool         `bson:"is_published,omitempty" json:"isPublished"`
	PublishedAt time.Time    `bson:"published_at,omitempty" json:"publishedAt"`
}

type UnPublishPresetInput struct {
	IsPublished bool      `bson:"is_published"`
	PublishedAt time.Time `bson:"published_at"`
}

type UpdateTagsInput struct {
	Id        string    `bson:"id,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

type BulkCreatePresetInput struct {
	UserId   string `bson:"user_id" json:"userId"`
	BulkData []struct {
		Label string      `bson:"label" json:"label"`
		Model PresetModel `bson:"model" json:"model"`
	} `json:"bulkData"`
}

type PartialSearchRoPresetInput struct {
	Id           *string `bson:"id,omitempty"`
	UserId       *string `bson:"user_id,omitempty"`
	ClassId      *int    `bson:"class_id,omitempty"`
	Label        *string `bson:"label,omitempty"`
	Skip         *int
	Take         *int
	InCludeModel bool
}

type IdSearchInput struct {
	Id string `bson:"id"`
}

type PartialSearchRoPresetResult struct {
	Items []RoPreset
	Total int64
}

type FindPresetByIdInput struct {
	Id           string
	InCludeModel bool
}

type RoPresetRepository interface {
	FindPresetById(FindPresetByIdInput) (*RoPreset, error)
	FindPresetByIds([]string) ([]RoPreset, error)
	PartialSearchPresets(PartialSearchRoPresetInput) (*PartialSearchRoPresetResult, error)
	CreatePreset(CreatePresetInput) (*RoPreset, error)
	CreatePresets(BulkCreatePresetInput) ([]RoPreset, error)
	UpdatePreset(id string, i UpdatePresetInput) error
	UnpublishedPreset(id string) error
	DeletePresetById(string) (*int, error)
}
