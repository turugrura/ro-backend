package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	Id        string      `bson:"_id,omitempty" json:"id"`
	UserId    string      `bson:"user_id" json:"userId"`
	Label     string      `bson:"label" json:"label"`
	Model     PresetModel `bson:"model" json:"model"`
	CreatedAt string      `bson:"created_at"`
	UpdatedAt string      `bson:"updated_at"`
}

type CreatePresetInput struct {
	UserId string      `bson:"user_id" json:"userId"`
	Label  string      `bson:"label" json:"label"`
	Model  PresetModel `bson:"model" json:"model"`
}

type UpdatePresetInput struct {
	Id     string      `bson:"_id,omitempty" json:"id"`
	UserId string      `bson:"user_id" json:"userId"`
	Label  string      `bson:"label" json:"label"`
	Model  PresetModel `bson:"model" json:"model"`
}

type BulkCreatePresetInput struct {
	UserId   string `bson:"user_id" json:"userId"`
	BulkData []struct {
		Label string      `bson:"label" json:"label"`
		Model PresetModel `bson:"model" json:"model"`
	} `json:"bulkData"`
}

type PartialSearchRoPreset struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId string             `bson:"user_id,omitempty" json:"userId"`
	Label  string             `bson:"label,omitempty" json:"label"`
}

type FindPreset struct {
	Id    string `bson:"_id,omitempty" json:"id"`
	Label string `bson:"label" json:"label"`
}

type RoPresetRepository interface {
	FindPresetById(string) (*RoPreset, error)
	FindPresetsByUserId(string) (*[]FindPreset, error)
	CreatePreset(CreatePresetInput) (*RoPreset, error)
	UpdatePreset(UpdatePresetInput) (*RoPreset, error)
	CreatePresets(BulkCreatePresetInput) (*[]RoPreset, error)
	DeletePresetById(string) (*int, error)
}
