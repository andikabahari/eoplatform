package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	Create(bankAccount *model.BankAccount)
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) *bankAccountRepository {
	return &bankAccountRepository{db}
}

func (r *bankAccountRepository) Create(bankAccount *model.BankAccount) {
	r.db.Debug().Save(bankAccount)
}
