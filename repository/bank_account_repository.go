package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	Get(bankAccounts *[]model.BankAccount, userID uint)
	Find(bankAccount *model.BankAccount, id uint)
	Create(bankAccount *model.BankAccount)
	Delete(bankAccount *model.BankAccount)
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

func (r *bankAccountRepository) Find(bankAccount *model.BankAccount, id uint) {
	r.db.Debug().Where("id = ?", id).Find(bankAccount)
}

func (r *bankAccountRepository) Create(bankAccount *model.BankAccount) {
	r.db.Debug().Save(bankAccount)
}

func (r *bankAccountRepository) Delete(bankAccount *model.BankAccount) {
	r.db.Debug().Delete(bankAccount)
}
