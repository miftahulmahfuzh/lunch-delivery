package testutils

import (
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
)

// Mock functions that return proper model types

func MockCompany(id int, name string) models.Company {
	return models.Company{
		ID:      id,
		Name:    name,
		Address: "123 Test St",
		Contact: "test@company.com",
	}
}

func MockEmployee(id, companyID int, name, email string) models.Employee {
	return models.Employee{
		ID:           id,
		CompanyID:    companyID,
		Name:         name,
		Email:        email,
		PasswordHash: "hashed_password",
		Active:       true,
	}
}

func MockMenuItem(id int, name string, price int) models.MenuItem {
	return models.MenuItem{
		ID:     id,
		Name:   name,
		Price:  price,
		Active: true,
	}
}