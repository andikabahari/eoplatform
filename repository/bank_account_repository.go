package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	Get(bankAccounts *[]model.BankAccount, userID uint)
	Find(bankAccount *model.BankAccount, id uint)
	FindByUserID(bankAccount *model.BankAccount, userID uint)
	Create(bankAccount *model.BankAccount)
	Update(bankAccount *model.BankAccount, req *request.UpdateBankAccountRequest)
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

func (r *bankAccountRepository) FindByUserID(bankAccount *model.BankAccount, userID uint) {
	r.db.Debug().Where("user_id = ?", userID).Find(bankAccount)
}

func (r *bankAccountRepository) Find(bankAccount *model.BankAccount, id uint) {
	r.db.Debug().Where("id = ?", id).Find(bankAccount)
}

func (r *bankAccountRepository) Create(bankAccount *model.BankAccount) {
	r.db.Debug().Omit("User").Save(bankAccount)
}

func (r *bankAccountRepository) Update(bankAccount *model.BankAccount, req *request.UpdateBankAccountRequest) {
	bankAccount.Bank = req.Bank
	bankAccount.VANumber = req.VANumber

	r.db.Debug().Omit("User").Save(bankAccount)
}

func (r *bankAccountRepository) Delete(bankAccount *model.BankAccount) {
	r.db.Debug().Delete(bankAccount)
}
