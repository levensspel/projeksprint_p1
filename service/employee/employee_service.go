package user_service

import (
	"context"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
)

type EmployeeService interface {
	GetEmployees(input dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error)
	GetEmployeesWithJoin(input dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error)
}

type service struct {
	employeeRepo repositories.EmployeeRepository
	logger       logger.Logger
}

func NewEmployeeService(
	employeeRepo repositories.EmployeeRepository,
	logger logger.Logger,
) EmployeeService {
	return &service{
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *service) GetEmployees(input dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error) {
	param := &dto.GetEmployeesRequest{
		Limit:          input.Limit,
		Offset:         input.Offset,
		IdentityNumber: input.IdentityNumber,
		Name:           input.Name,
		Gender:         input.Gender,
		DepartmentID:   input.DepartmentID,
		ManagerID:      input.ManagerID,
	}

	employees, err := s.employeeRepo.GetEmployees(context.Background(), param)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, param)
		return []dto.GetEmployeeResponseItem{}, err
	}

	return employees, nil
}

func (s *service) GetEmployeesWithJoin(input dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error) {
	param := &dto.GetEmployeesRequest{
		Limit:          input.Limit,
		Offset:         input.Offset,
		IdentityNumber: input.IdentityNumber,
		Name:           input.Name,
		Gender:         input.Gender,
		DepartmentID:   input.DepartmentID,
		ManagerID:      input.ManagerID,
	}

	employees, err := s.employeeRepo.GetEmployeesWithJoin(context.Background(), param)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, param)
		return []dto.GetEmployeeResponseItem{}, err
	}

	return employees, nil
}
