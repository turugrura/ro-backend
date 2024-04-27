package service

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"ro-backend/repository"
	"slices"
	"sort"
	"strings"
)

func NewSummaryPresetService(pRepo repository.RoPresetRepository) PresetSummaryService {
	return summaryPresetService{pRepo: pRepo}
}

type summaryPresetService struct {
	pRepo repository.RoPresetRepository
}

type EnchantSummary struct {
	ItemId int
	Total  int
}

type ItemSummary struct {
	Total    int
	Enchants map[string]int
}

func (s summaryPresetService) GenerateSummary() (*PresetSummary, error) {
	skip := int(0)
	take := int(500)
	res, err := s.pRepo.PartialSearchPresets(repository.PartialSearchRoPresetInput{
		InCludeModel: true,
		Skip:         &skip,
		Take:         &take,
	})

	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return &PresetSummary{}, nil
	}

	userDataMap := AllSummary{}
	presetSummaryMap := map[int]map[string]int{}
	setSummary(&userDataMap, res.Items, &presetSummaryMap)

	total := int64(res.Total)
	round := int(math.Ceil(float64(total) / float64(take)))
	fmt.Printf("total %v\n", res.Total)
	for i := 1; i < round; i++ {
		skip = i * take
		res, err = s.pRepo.PartialSearchPresets(repository.PartialSearchRoPresetInput{
			InCludeModel: true,
			Skip:         &skip,
			Take:         &take,
		})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Round %v passed\n", i)

		setSummary(&userDataMap, res.Items, &presetSummaryMap)
	}

	type UsingItemFrequency struct {
		TotalPreset      int
		TotalAccount     int
		TotalEnchant     int
		ItemUsingRate    float64
		EnchantUsingRate map[string]float64
	}
	allUserSummary := map[int]map[string]map[string]map[int]*UsingItemFrequency{}
	summaryClassSkillMap := map[int]map[string]int{}
	totalSelectedJobMap := map[int]int{}

	for _, jobMap := range userDataMap {
		isAlreadyCalc := map[int]int{}
		totalSelectedJobMap[0] += 1

		for jobId, skillEquipmentMap := range jobMap {
			if isAlreadyCalc[jobId] == 0 {
				isAlreadyCalc[jobId] = 1
				totalSelectedJobMap[jobId] += 1
			}

			if allUserSummary[jobId] == nil {
				allUserSummary[jobId] = map[string]map[string]map[int]*UsingItemFrequency{}
			}

			if summaryClassSkillMap[jobId] == nil {
				summaryClassSkillMap[jobId] = map[string]int{}
			}

			for skillName, equipmentMap := range skillEquipmentMap {
				summaryClassSkillMap[jobId][skillName] += 1

				if allUserSummary[jobId][skillName] == nil {
					allUserSummary[jobId][skillName] = map[string]map[int]*UsingItemFrequency{}
				}

				for itemPosition, usingItemMap := range equipmentMap {
					if allUserSummary[jobId][skillName][itemPosition] == nil {
						allUserSummary[jobId][skillName][itemPosition] = map[int]*UsingItemFrequency{}
					}

					totalItem := 0
					totalEnchantMap := map[int]int{}
					for itemId, itemSummary := range usingItemMap {
						totalItem += itemSummary.Total
						for _, totalEnchant := range itemSummary.Enchants {
							totalEnchantMap[itemId] += totalEnchant
						}
					}

					for itemId, itemSummary := range usingItemMap {
						if allUserSummary[jobId][skillName][itemPosition][itemId] == nil {
							allUserSummary[jobId][skillName][itemPosition][itemId] = &UsingItemFrequency{
								EnchantUsingRate: map[string]float64{},
							}
						}
						allUserSummary[jobId][skillName][itemPosition][itemId].ItemUsingRate += float64(itemSummary.Total) / float64(totalItem)
						allUserSummary[jobId][skillName][itemPosition][itemId].TotalEnchant += totalEnchantMap[itemId]
						allUserSummary[jobId][skillName][itemPosition][itemId].TotalPreset += itemSummary.Total
						allUserSummary[jobId][skillName][itemPosition][itemId].TotalAccount += 1

						for enchantStr, totalEnchant := range itemSummary.Enchants {
							allUserSummary[jobId][skillName][itemPosition][itemId].EnchantUsingRate[enchantStr] += float64(totalEnchant) / float64(totalEnchantMap[itemId])
						}
					}
				}
			}
		}
	}

	totalRanking := 10
	jobSummary := map[int]map[string]map[string][]RankingSummary{}
	for jobId, skillEquipmentMap := range allUserSummary {
		if jobSummary[jobId] == nil {
			jobSummary[jobId] = map[string]map[string][]RankingSummary{}
		}

		for skillName, equipmentMap := range skillEquipmentMap {
			if jobSummary[jobId][skillName] == nil {
				jobSummary[jobId][skillName] = map[string][]RankingSummary{}
			}

			for itemPosition, usingItemMap := range equipmentMap {
				rankings := []RankingSummary{}
				for itemId, usingFrequency := range usingItemMap {
					rankings = append(rankings, RankingSummary{
						TotalPreset:  usingFrequency.TotalPreset,
						TotalAccount: usingFrequency.TotalAccount,
						TotalEnchant: usingFrequency.TotalEnchant,
						ItemId:       itemId,
						UsingRate:    usingFrequency.ItemUsingRate,
						Enchants:     usingFrequency.EnchantUsingRate,
					})
				}
				slices.SortFunc(rankings, func(a, b RankingSummary) int {
					return cmp.Compare(a.UsingRate, b.UsingRate) * -1
				})
				if len(rankings) > totalRanking {
					etc := rankings[totalRanking:]
					etcRate := float64(0)
					for _, v := range etc {
						etcRate += v.UsingRate
					}

					rankings = rankings[:totalRanking]
				}

				jobSummary[jobId][skillName][itemPosition] = rankings
			}
		}
	}

	writeJsonFile(presetSummaryMap, "x_presetSummaryMap.json")
	writeJsonFile(summaryClassSkillMap, "x_summaryClassSkillMap.json")
	writeJsonFile(totalSelectedJobMap, "x_totalSelectedJobMap.json")
	writeJsonFile(jobSummary, "x.json")

	return &PresetSummary{
		SummaryClassSkillMap: summaryClassSkillMap,
		TotalSelectedJobMap:  totalSelectedJobMap,
		JobSummary:           jobSummary,
	}, nil
}

func setSummary(summary *AllSummary, presets []repository.RoPreset, presetSummary *map[int]map[string]int) {
	for _, preset := range presets {
		skillName := getSkillName(preset.Model.SelectedAtkSkill)

		if ((*presetSummary)[preset.ClassId]) == nil {
			(*presetSummary)[preset.ClassId] = map[string]int{}
		}
		(*presetSummary)[preset.ClassId][skillName] += 1

		if (*summary)[preset.UserId] == nil {
			(*summary)[preset.UserId] = JobSummary{}
			(*summary)[preset.UserId][preset.ClassId] = UsingSkillSummary{}
			(*summary)[preset.UserId][preset.ClassId][skillName] = map[string]map[int]*ItemSummary{}
			initJobSummary(summary, skillName, preset)
		} else if (*summary)[preset.UserId][preset.ClassId] == nil {
			(*summary)[preset.UserId][preset.ClassId] = UsingSkillSummary{}
			(*summary)[preset.UserId][preset.ClassId][skillName] = map[string]map[int]*ItemSummary{}
			initJobSummary(summary, skillName, preset)
		} else if (*summary)[preset.UserId][preset.ClassId][skillName] == nil {
			(*summary)[preset.UserId][preset.ClassId][skillName] = map[string]map[int]*ItemSummary{}
			initJobSummary(summary, skillName, preset)
		}

		itemSummary := (*summary)[preset.UserId][preset.ClassId][skillName]

		if preset.Model.Weapon != 0 {
			if itemSummary["Weapon"][preset.Model.Weapon] == nil {
				itemSummary["Weapon"][preset.Model.Weapon] = &ItemSummary{}
			}
			itemSummary["Weapon"][preset.Model.Weapon].Total += 1
			if itemSummary["Weapon"][preset.Model.Weapon].Enchants == nil {
				itemSummary["Weapon"][preset.Model.Weapon].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.WeaponEnchant1, preset.Model.WeaponEnchant2, preset.Model.WeaponEnchant3)
			itemSummary["Weapon"][preset.Model.Weapon].Enchants[enchant] += 1

			if preset.Model.WeaponCard1 != 0 {
				if itemSummary["WeaponCard"][preset.Model.WeaponCard1] == nil {
					itemSummary["WeaponCard"][preset.Model.WeaponCard1] = &ItemSummary{}
				}
				itemSummary["WeaponCard"][preset.Model.WeaponCard1].Total += 1
			}
			if preset.Model.WeaponCard2 != 0 {
				if itemSummary["WeaponCard"][preset.Model.WeaponCard2] == nil {
					itemSummary["WeaponCard"][preset.Model.WeaponCard2] = &ItemSummary{}
				}
				itemSummary["WeaponCard"][preset.Model.WeaponCard2].Total += 1
			}
			if preset.Model.WeaponCard3 != 0 {
				if itemSummary["WeaponCard"][preset.Model.WeaponCard3] == nil {
					itemSummary["WeaponCard"][preset.Model.WeaponCard3] = &ItemSummary{}
				}
				itemSummary["WeaponCard"][preset.Model.WeaponCard3].Total += 1
			}
			if preset.Model.WeaponCard4 != 0 {
				if itemSummary["WeaponCard"][preset.Model.WeaponCard4] == nil {
					itemSummary["WeaponCard"][preset.Model.WeaponCard4] = &ItemSummary{}
				}
				itemSummary["WeaponCard"][preset.Model.WeaponCard4].Total += 1
			}
		}

		if preset.Model.LeftWeapon != 0 {
			if itemSummary["LeftWeapon"][preset.Model.LeftWeapon] == nil {
				itemSummary["LeftWeapon"][preset.Model.LeftWeapon] = &ItemSummary{}
			}

			itemSummary["LeftWeapon"][preset.Model.LeftWeapon].Total += 1
			if itemSummary["LeftWeapon"][preset.Model.LeftWeapon].Enchants == nil {
				itemSummary["LeftWeapon"][preset.Model.LeftWeapon].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.LeftWeaponEnchant1, preset.Model.LeftWeaponEnchant2, preset.Model.LeftWeaponEnchant3)
			itemSummary["LeftWeapon"][preset.Model.LeftWeapon].Enchants[enchant] += 1

			if preset.Model.LeftWeaponCard1 != 0 {
				if itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard1] == nil {
					itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard1] = &ItemSummary{}
				}
				itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard1].Total += 1
			}
			if preset.Model.LeftWeaponCard2 != 0 {
				if itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard2] == nil {
					itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard2] = &ItemSummary{}
				}
				itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard2].Total += 1
			}
			if preset.Model.LeftWeaponCard3 != 0 {
				if itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard3] == nil {
					itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard3] = &ItemSummary{}
				}
				itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard3].Total += 1
			}
			if preset.Model.LeftWeaponCard4 != 0 {
				if itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard4] == nil {
					itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard4] = &ItemSummary{}
				}
				itemSummary["LeftWeaponCard"][preset.Model.LeftWeaponCard4].Total += 1
			}
		}

		if preset.Model.Shield != 0 {
			if itemSummary["Shield"][preset.Model.Shield] == nil {
				itemSummary["Shield"][preset.Model.Shield] = &ItemSummary{}
			}

			itemSummary["Shield"][preset.Model.Shield].Total += 1
			if itemSummary["Shield"][preset.Model.Shield].Enchants == nil {
				itemSummary["Shield"][preset.Model.Shield].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.ShieldEnchant1, preset.Model.ShieldEnchant2, preset.Model.ShieldEnchant3)
			itemSummary["Shield"][preset.Model.Shield].Enchants[enchant] += 1

			if preset.Model.ShieldCard != 0 {
				if itemSummary["ShieldCard"][preset.Model.ShieldCard] == nil {
					itemSummary["ShieldCard"][preset.Model.ShieldCard] = &ItemSummary{}
				}
				itemSummary["ShieldCard"][preset.Model.ShieldCard].Total += 1
			}
		}

		if preset.Model.HeadUpper != 0 {
			if itemSummary["HeadUpper"][preset.Model.HeadUpper] == nil {
				itemSummary["HeadUpper"][preset.Model.HeadUpper] = &ItemSummary{}
			}
			itemSummary["HeadUpper"][preset.Model.HeadUpper].Total += 1
			if itemSummary["HeadUpper"][preset.Model.HeadUpper].Enchants == nil {
				itemSummary["HeadUpper"][preset.Model.HeadUpper].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.HeadUpperEnchant1, preset.Model.HeadUpperEnchant2, preset.Model.HeadUpperEnchant3)
			itemSummary["HeadUpper"][preset.Model.HeadUpper].Enchants[enchant] += 1

			if preset.Model.HeadUpperCard != 0 {
				if itemSummary["HeadUpperCard"][preset.Model.HeadUpperCard] == nil {
					itemSummary["HeadUpperCard"][preset.Model.HeadUpperCard] = &ItemSummary{}
				}
				itemSummary["HeadUpperCard"][preset.Model.HeadUpperCard].Total += 1
			}
		}

		if preset.Model.HeadMiddle != 0 {
			if itemSummary["HeadMiddle"][preset.Model.HeadMiddle] == nil {
				itemSummary["HeadMiddle"][preset.Model.HeadMiddle] = &ItemSummary{}
			}
			itemSummary["HeadMiddle"][preset.Model.HeadMiddle].Total += 1
			if itemSummary["HeadMiddle"][preset.Model.HeadMiddle].Enchants == nil {
				itemSummary["HeadMiddle"][preset.Model.HeadMiddle].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.HeadMiddleEnchant1, preset.Model.HeadMiddleEnchant2, preset.Model.HeadMiddleEnchant3)
			itemSummary["HeadMiddle"][preset.Model.HeadMiddle].Enchants[enchant] += 1

			if preset.Model.HeadMiddleCard != 0 {
				if itemSummary["HeadMiddleCard"][preset.Model.HeadMiddleCard] == nil {
					itemSummary["HeadMiddleCard"][preset.Model.HeadMiddleCard] = &ItemSummary{}
				}
				itemSummary["HeadMiddleCard"][preset.Model.HeadMiddleCard].Total += 1
			}
		}

		if preset.Model.HeadLower != 0 {
			if itemSummary["HeadLower"][preset.Model.HeadLower] == nil {
				itemSummary["HeadLower"][preset.Model.HeadLower] = &ItemSummary{}
			}
			itemSummary["HeadLower"][preset.Model.HeadLower].Total += 1
			if itemSummary["HeadLower"][preset.Model.HeadLower].Enchants == nil {
				itemSummary["HeadLower"][preset.Model.HeadLower].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.HeadLowerEnchant1, preset.Model.HeadLowerEnchant2, preset.Model.HeadLowerEnchant3)
			itemSummary["HeadLower"][preset.Model.HeadLower].Enchants[enchant] += 1
		}

		if preset.Model.Armor != 0 {
			if itemSummary["Armor"][preset.Model.Armor] == nil {
				itemSummary["Armor"][preset.Model.Armor] = &ItemSummary{}
			}
			itemSummary["Armor"][preset.Model.Armor].Total += 1
			if itemSummary["Armor"][preset.Model.Armor].Enchants == nil {
				itemSummary["Armor"][preset.Model.Armor].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.ArmorEnchant1, preset.Model.ArmorEnchant2, preset.Model.ArmorEnchant3)
			itemSummary["Armor"][preset.Model.Armor].Enchants[enchant] += 1

			if preset.Model.ArmorCard != 0 {
				if itemSummary["ArmorCard"][preset.Model.ArmorCard] == nil {
					itemSummary["ArmorCard"][preset.Model.ArmorCard] = &ItemSummary{}
				}
				itemSummary["ArmorCard"][preset.Model.ArmorCard].Total += 1
			}
		}

		if preset.Model.Garment != 0 {
			if itemSummary["Garment"][preset.Model.Garment] == nil {
				itemSummary["Garment"][preset.Model.Garment] = &ItemSummary{}
			}
			itemSummary["Garment"][preset.Model.Garment].Total += 1
			if itemSummary["Garment"][preset.Model.Garment].Enchants == nil {
				itemSummary["Garment"][preset.Model.Garment].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.GarmentEnchant1, preset.Model.GarmentEnchant2, preset.Model.GarmentEnchant3)
			itemSummary["Garment"][preset.Model.Garment].Enchants[enchant] += 1

			if preset.Model.GarmentCard != 0 {
				if itemSummary["GarmentCard"][preset.Model.GarmentCard] == nil {
					itemSummary["GarmentCard"][preset.Model.GarmentCard] = &ItemSummary{}
				}
				itemSummary["GarmentCard"][preset.Model.GarmentCard].Total += 1
			}
		}

		if preset.Model.Boot != 0 {
			if itemSummary["Boot"][preset.Model.Boot] == nil {
				itemSummary["Boot"][preset.Model.Boot] = &ItemSummary{}
			}
			itemSummary["Boot"][preset.Model.Boot].Total += 1
			if itemSummary["Boot"][preset.Model.Boot].Enchants == nil {
				itemSummary["Boot"][preset.Model.Boot].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.BootEnchant1, preset.Model.BootEnchant2, preset.Model.BootEnchant3)
			itemSummary["Boot"][preset.Model.Boot].Enchants[enchant] += 1

			if preset.Model.BootCard != 0 {
				if itemSummary["BootCard"][preset.Model.BootCard] == nil {
					itemSummary["BootCard"][preset.Model.BootCard] = &ItemSummary{}
				}
				itemSummary["BootCard"][preset.Model.BootCard].Total += 1
			}
		}

		if preset.Model.AccLeft != 0 {
			if itemSummary["AccLeft"][preset.Model.AccLeft] == nil {
				itemSummary["AccLeft"][preset.Model.AccLeft] = &ItemSummary{}
			}
			itemSummary["AccLeft"][preset.Model.AccLeft].Total += 1
			if itemSummary["AccLeft"][preset.Model.AccLeft].Enchants == nil {
				itemSummary["AccLeft"][preset.Model.AccLeft].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.AccLeftEnchant1, preset.Model.AccLeftEnchant2, preset.Model.AccLeftEnchant3)
			itemSummary["AccLeft"][preset.Model.AccLeft].Enchants[enchant] += 1

			if preset.Model.AccLeftCard != 0 {
				if itemSummary["AccLeftCard"][preset.Model.AccLeftCard] == nil {
					itemSummary["AccLeftCard"][preset.Model.AccLeftCard] = &ItemSummary{}
				}
				itemSummary["AccLeftCard"][preset.Model.AccLeftCard].Total += 1
			}
		}

		if preset.Model.AccRight != 0 {
			if itemSummary["AccRight"][preset.Model.AccRight] == nil {
				itemSummary["AccRight"][preset.Model.AccRight] = &ItemSummary{}
			}
			itemSummary["AccRight"][preset.Model.AccRight].Total += 1
			if itemSummary["AccRight"][preset.Model.AccRight].Enchants == nil {
				itemSummary["AccRight"][preset.Model.AccRight].Enchants = map[string]int{}
			}

			enchant := buildEnchantStr(preset.Model.AccRightEnchant1, preset.Model.AccRightEnchant2, preset.Model.AccRightEnchant3)
			itemSummary["AccRight"][preset.Model.AccRight].Enchants[enchant] += 1

			if preset.Model.AccRightCard != 0 {
				if itemSummary["AccRightCard"][preset.Model.AccRightCard] == nil {
					itemSummary["AccRightCard"][preset.Model.AccRightCard] = &ItemSummary{}
				}
				itemSummary["AccRightCard"][preset.Model.AccRightCard].Total += 1
			}
		}

		if preset.Model.CostumeEnchantUpper != 0 {
			if itemSummary["CostumeEnchantUpper"][preset.Model.CostumeEnchantUpper] == nil {
				itemSummary["CostumeEnchantUpper"][preset.Model.CostumeEnchantUpper] = &ItemSummary{}
			}
			itemSummary["CostumeEnchantUpper"][preset.Model.CostumeEnchantUpper].Total += 1
		}
		if preset.Model.CostumeEnchantMiddle != 0 {
			if itemSummary["CostumeEnchantMiddle"][preset.Model.CostumeEnchantMiddle] == nil {
				itemSummary["CostumeEnchantMiddle"][preset.Model.CostumeEnchantMiddle] = &ItemSummary{}
			}
			itemSummary["CostumeEnchantMiddle"][preset.Model.CostumeEnchantMiddle].Total += 1
		}
		if preset.Model.CostumeEnchantLower != 0 {
			if itemSummary["CostumeEnchantLower"][preset.Model.CostumeEnchantLower] == nil {
				itemSummary["CostumeEnchantLower"][preset.Model.CostumeEnchantLower] = &ItemSummary{}
			}
			itemSummary["CostumeEnchantLower"][preset.Model.CostumeEnchantLower].Total += 1
		}
		if preset.Model.CostumeEnchantGarment != 0 {
			if itemSummary["CostumeEnchantGarment"][preset.Model.CostumeEnchantGarment] == nil {
				itemSummary["CostumeEnchantGarment"][preset.Model.CostumeEnchantGarment] = &ItemSummary{}
			}
			itemSummary["CostumeEnchantGarment"][preset.Model.CostumeEnchantGarment].Total += 1
		}

		if preset.Model.ShadowWeapon != 0 {
			if itemSummary["ShadowWeapon"][preset.Model.ShadowWeapon] == nil {
				itemSummary["ShadowWeapon"][preset.Model.ShadowWeapon] = &ItemSummary{}
			}
			itemSummary["ShadowWeapon"][preset.Model.ShadowWeapon].Total += 1
		}
		if preset.Model.ShadowArmor != 0 {
			if itemSummary["ShadowArmor"][preset.Model.ShadowArmor] == nil {
				itemSummary["ShadowArmor"][preset.Model.ShadowArmor] = &ItemSummary{}
			}
			itemSummary["ShadowArmor"][preset.Model.ShadowArmor].Total += 1
		}
		if preset.Model.ShadowShield != 0 {
			if itemSummary["ShadowShield"][preset.Model.ShadowShield] == nil {
				itemSummary["ShadowShield"][preset.Model.ShadowShield] = &ItemSummary{}
			}
			itemSummary["ShadowShield"][preset.Model.ShadowShield].Total += 1
		}
		if preset.Model.ShadowBoot != 0 {
			if itemSummary["ShadowBoot"][preset.Model.ShadowBoot] == nil {
				itemSummary["ShadowBoot"][preset.Model.ShadowBoot] = &ItemSummary{}
			}
			itemSummary["ShadowBoot"][preset.Model.ShadowBoot].Total += 1
		}
		if preset.Model.ShadowEarring != 0 {
			if itemSummary["ShadowEarring"][preset.Model.ShadowEarring] == nil {
				itemSummary["ShadowEarring"][preset.Model.ShadowEarring] = &ItemSummary{}
			}
			itemSummary["ShadowEarring"][preset.Model.ShadowEarring].Total += 1
		}
		if preset.Model.ShadowPendant != 0 {
			if itemSummary["ShadowPendant"][preset.Model.ShadowPendant] == nil {
				itemSummary["ShadowPendant"][preset.Model.ShadowPendant] = &ItemSummary{}
			}
			itemSummary["ShadowPendant"][preset.Model.ShadowPendant].Total += 1
		}
	}
}

func initJobSummary(summary *AllSummary, skillName string, preset repository.RoPreset) {
	itemSummary := (*summary)[preset.UserId][preset.ClassId][skillName]
	itemSummary["Weapon"] = map[int]*ItemSummary{}
	itemSummary["WeaponCard"] = map[int]*ItemSummary{}
	itemSummary["LeftWeapon"] = map[int]*ItemSummary{}
	itemSummary["LeftWeaponCard"] = map[int]*ItemSummary{}
	itemSummary["Shield"] = map[int]*ItemSummary{}
	itemSummary["ShieldCard"] = map[int]*ItemSummary{}
	itemSummary["HeadUpper"] = map[int]*ItemSummary{}
	itemSummary["HeadUpperCard"] = map[int]*ItemSummary{}
	itemSummary["HeadMiddle"] = map[int]*ItemSummary{}
	itemSummary["HeadMiddleCard"] = map[int]*ItemSummary{}
	itemSummary["HeadLower"] = map[int]*ItemSummary{}
	itemSummary["Armor"] = map[int]*ItemSummary{}
	itemSummary["ArmorCard"] = map[int]*ItemSummary{}
	itemSummary["Garment"] = map[int]*ItemSummary{}
	itemSummary["GarmentCard"] = map[int]*ItemSummary{}
	itemSummary["Boot"] = map[int]*ItemSummary{}
	itemSummary["BootCard"] = map[int]*ItemSummary{}
	itemSummary["AccLeft"] = map[int]*ItemSummary{}
	itemSummary["AccLeftCard"] = map[int]*ItemSummary{}
	itemSummary["AccRight"] = map[int]*ItemSummary{}
	itemSummary["AccRightCard"] = map[int]*ItemSummary{}

	itemSummary["CostumeEnchantUpper"] = map[int]*ItemSummary{}
	itemSummary["CostumeEnchantMiddle"] = map[int]*ItemSummary{}
	itemSummary["CostumeEnchantLower"] = map[int]*ItemSummary{}
	itemSummary["CostumeEnchantGarment"] = map[int]*ItemSummary{}

	itemSummary["ShadowWeapon"] = map[int]*ItemSummary{}
	itemSummary["ShadowArmor"] = map[int]*ItemSummary{}
	itemSummary["ShadowShield"] = map[int]*ItemSummary{}
	itemSummary["ShadowBoot"] = map[int]*ItemSummary{}
	itemSummary["ShadowEarring"] = map[int]*ItemSummary{}
	itemSummary["ShadowPendant"] = map[int]*ItemSummary{}
}

func buildEnchantStr(e1, e2, e3 int) string {
	arr := []int{e1, e2, e3}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	return fmt.Sprintf("%v-%v-%v", arr[0], arr[1], arr[2])
}

var mapSkillName = map[string]string{
	"Ignition Break3":             "Ignition Break",
	"Ignition Break7":             "Ignition Break",
	"Cross Slash(Cross Wound)":    "Cross Slash",
	"Silvervine Stem Spear Earth": "Silvervine Stem Spear",
	"Tiger Cannon Combo":          "Tiger Cannon",
	"Adoramus Ancilla":            "Adoramus",
}

func getSkillName(rawSkillName string) string {
	skillName := strings.Split(rawSkillName, "==")[0]
	skillName = strings.ReplaceAll(skillName, "[Improved 2nd] ", "")
	skillName = strings.ReplaceAll(skillName, "[Improved 2rd] ", "")
	skillName = strings.ReplaceAll(skillName, "[Improved 1st] ", "")
	skillName = strings.ReplaceAll(skillName, "[Improved 1nd] ", "")
	skillName = strings.ReplaceAll(skillName, "[Improved] ", "")

	mapped := mapSkillName[skillName]
	if mapped != "" {
		return mapped
	}

	return skillName
}

func writeJsonFile(data interface{}, filename string) {
	file2, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile(filename, file2, fs.ModePerm)
}
