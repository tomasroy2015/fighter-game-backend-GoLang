package domain

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Fighter interface {
	GetID() string
	GetPower() float64
}

type Knight struct {
	ID          string `gorm:"size:255;not null;unique" json:"id"`
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	Strength    uint64 `gorm:"not null" json:"strength"`
	WeaponPower uint64 `gorm:"not null" json:"weapon_power"`
}

func (p *Knight) Prepare() {
	p.ID = html.EscapeString(strings.TrimSpace(p.ID))
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Strength = 0
	p.WeaponPower = 0

}

func (p *Knight) Validate() error {
	if p.Strength <= 0 {
		return errors.New("Required Strength")
	}
	if p.WeaponPower <= 0 {
		return errors.New("Required WeaponPower")
	}
	return nil
}

func (u *Knight) SaveKnight(db *gorm.DB) (*Knight, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Knight{}, err
	}
	return u, nil
}

func (u *Knight) FindAllKnights(db *gorm.DB) (*[]Knight, error) {
	var err error
	knights := []Knight{}
	err = db.Debug().Model(&Knight{}).Limit(100).Find(&knights).Error
	if err != nil {
		return &[]Knight{}, err
	}
	return &knights, err
}

func (u *Knight) FindKnightByID(db *gorm.DB, uid uint64) (*Knight, error) {
	var err error
	err = db.Debug().Model(Knight{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Knight{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Knight{}, errors.New("Knight Not Found")
	}
	return u, err
}
