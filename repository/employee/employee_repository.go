package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/dto"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetEmployees(ctx context.Context, input *dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error) {
	//	The constructed query if all the parameters are sent:
	// `SELECT
	// 		identity_number,
	// 		name,
	// 		gender,
	// 		image_uri,
	// 		department_id
	// FROM employees
	// WHERE
	//  	manager_id = $1
	// 		AND LOWER(identity_number) ILIKE $2 || '%'
	// 		AND name ILIKE '%' || $3 || '%'"
	// 		AND gender = $4
	// 		AND department_id = $5
	// LIMIT $5
	// OFFSET $6`

	query := "SELECT identitynumber, name, gender, employeeimageuri, departmentid"
	conditions := "WHERE managerid = $1"
	argIndex := 2
	var args []interface{}
	args = append(args, input.ManagerID)

	if input.IdentityNumber != "" {
		args = append(args, input.IdentityNumber)
		conditions += fmt.Sprintf(" AND LOWER(identitynumber) ILIKE $%d || '%s'", argIndex, "%") // eg. "AND LOWER(identity_number) ILIKE $2 || '%'"
		argIndex++
	}
	if input.Name != "" {
		args = append(args, input.Name)
		conditions += fmt.Sprintf(" AND name ILIKE '%s' || $%d || '%s'", "%", argIndex, "%") // eg. "AND name ILIKE '%' || $2 || '%'"
		argIndex++
	}
	if input.Gender != "" {
		args = append(args, input.Gender)
		conditions += fmt.Sprintf(" AND gender = $%d", argIndex)
		argIndex++
	}
	if input.DepartmentID != "" {
		args = append(args, input.DepartmentID)
		conditions += fmt.Sprintf(" AND departmentid = $%d", argIndex)
		argIndex++
	}
	query = strings.TrimRight(query, ",") + " from employees "

	args = append(args, input.Limit)
	conditions += fmt.Sprintf(" LIMIT $%d", argIndex)
	argIndex++

	args = append(args, input.Offset)
	conditions += fmt.Sprintf(" OFFSET $%d;", argIndex)

	query += conditions
	// fmt.Println(query)
	// fmt.Println(args...)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var employees []dto.GetEmployeeResponseItem
	for rows.Next() {
		var employee dto.GetEmployeeResponseItem

		// The employee's attributes depends on the order of written columns in the SQL query
		err := rows.Scan(
			&employee.IdentityNumber,
			&employee.Name,
			&employee.Gender,
			&employee.EmployeeImageUri,
			&employee.DepartmentID,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetEmployeesWithJoin(ctx context.Context, input *dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error) {
	// Membuat query dinamis
	query := "SELECT e.identitynumber, e.name, e.gender, e.employeeimageuri, e.departmentid" // 'e' refer to 'employee e' which will be appended later
	conditions := "WHERE u.identityNumber = $1"                                              // 'u' refer to 'users u' which will be appended later
	argIndex := 2
	var args []interface{}
	args = append(args, input.ManagerID)

	//	The constructed query if all the parameters are sent:
	// 	SELECT
	// 		e.identitynumber,
	// 		e.name,
	// 		e.gender,
	// 		e.employeeimageuri,
	// 		e.departmentid
	//  FROM employees AS e
	//  LEFT JOIN department d
	//  ON e.departmentId = d.departmentId
	// 		LEFT JOIN users u
	// 		ON d.managerId = u.identityNumber
	//  WHERE
	// 		u.identityNumber = $1
	// 		AND LOWER(e.identitynumber) ILIKE $2 || '%'
	// 		AND e.name ILIKE '%' || $3 || '%'
	// 		AND e.gender = $4
	// 		AND e.departmentid = $5
	//  LIMIT $6
	//  OFFSET $7;

	if input.IdentityNumber != "" {
		args = append(args, input.IdentityNumber)
		conditions += fmt.Sprintf(" AND LOWER(e.identitynumber) ILIKE $%d || '%s'", argIndex, "%") // eg. AND identity_number ILIKE $2 || '%'
		argIndex++
	}
	if input.Name != "" {
		args = append(args, input.Name)
		conditions += fmt.Sprintf(" AND e.name ILIKE '%s' || $%d || '%s'", "%", argIndex, "%") // eg. AND name ILIKE %$2%
		argIndex++
	}
	if input.Gender != "" {
		args = append(args, input.Gender)
		conditions += fmt.Sprintf(" AND e.gender = $%d", argIndex)
		argIndex++
	}
	if input.DepartmentID != "" {
		args = append(args, input.DepartmentID)
		conditions += fmt.Sprintf(" AND e.departmentid = $%d", argIndex)
		argIndex++
	}
	query = strings.TrimRight(query, ",") + " FROM employees AS e  LEFT JOIN department d  ON e.departmentId = d.departmentId LEFT JOIN users u  ON d.managerId = u.identityNumber "

	args = append(args, input.Limit)
	conditions += fmt.Sprintf(" LIMIT $%d", argIndex)
	argIndex++

	args = append(args, input.Offset)
	conditions += fmt.Sprintf(" OFFSET $%d;", argIndex)

	query += conditions
	// fmt.Println(query)
	// fmt.Println(args...)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var employees []dto.GetEmployeeResponseItem
	for rows.Next() {
		var employee dto.GetEmployeeResponseItem
		err := rows.Scan(&employee.IdentityNumber, &employee.Name, &employee.EmployeeImageUri, &employee.Gender, &employee.DepartmentID)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}
