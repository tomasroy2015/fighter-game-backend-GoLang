package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"gitlab.com/zenport.io/go-assignment/api/models"
)

var knights = []models.Knight{
	models.Knight{
		Name:        "Steven victor",
		Strength:    200,
		WeaponPower: 10000,
	},
	models.Knight{
		Name:        "William Kane",
		Strength:    100,
		WeaponPower: 1000,
	},
	models.Knight{
		Name:        "Arnold carlton",
		Strength:    20,
		WeaponPower: 200,
	},
	models.Knight{
		Name:        "Saadman Lee",
		Strength:    70,
		WeaponPower: 700,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Knight{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Knight{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range knights {
		err = db.Debug().Model(&models.Knight{}).Create(&knights[i]).Error
		if err != nil {
			log.Fatalf("cannot seed knights table: %v", err)
		}

	}
}
