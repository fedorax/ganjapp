package models

import (
	"gorm.io/gorm"
)

// CannabisStrain represents a strain of cannabis, along with other data
type CannabisStrain struct {
	gorm.Model
	Breeder       string `gorm:"index:breeder;index:breeder_strain,unique"`
	Strain        string `gorm:"index:strain;index:breeder_strain,unique"`
	URL           string
	Genetics      string
	Environment   string
	FloweringTime string
	Gender        string
}

// FindStrainAndBreeder searches the database for strain and breeder matching (case insensitive) the passed values
func FindStrainAndBreeder(strain string, breeder string) (CannabisStrain, error) {
	var result CannabisStrain
	r := DB.Where("strain LIKE ?", strain).Where("breeder LIKE ?", breeder).First(&result)
	return result, r.Error
}

// FindStrainAndBreederLike searches the database for strain and breeder LIKE the passed values
func FindStrainAndBreederLike(strain string, breeder string) ([]CannabisStrain, error) {
	var results []CannabisStrain
	r := DB.Where("strain LIKE ?", "%"+strain+"%").Where("breeder LIKE ?", "%"+breeder+"%").Find(&results)
	return results, r.Error
}

// FindStrainLike searches the database for strain LIKE the passed value
func FindStrainLike(strain string) ([]CannabisStrain, error) {
	var results []CannabisStrain
	r := DB.Where("strain LIKE ?", "%"+strain+"%").Find(&results)
	return results, r.Error
}

// FindBreederLike searches the database for breeder LIKE the passed value
func FindBreederLike(breeder string) ([]CannabisStrain, error) {
	var results []CannabisStrain
	r := DB.Where("breeder LIKE ?", "%"+breeder+"%").Find(&results)
	return results, r.Error
}
