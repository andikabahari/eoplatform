package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	Get(bankAccounts *[]model.BankAccount, userID uint)
	Create(bankAccount *model.BankAccount)
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) *bankAccountRepository {
	return &bankAccountRepository{db}
}

func (r *bankAccountRepository) Get(bankAccounts *[]model.BankAccount, userID uint) {
	r.db.Debug().Where("user_id = ?", userID).Find(bankAccounts)
}

func (r *bankAccountRepository) Create(bankAccount *model.BankAccount) {
	r.db.Debug().Save(bankAccount)
}
