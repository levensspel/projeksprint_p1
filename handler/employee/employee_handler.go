package employeeHandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	service "github.com/levensspel/go-gin-template/service/employee"
	"github.com/levensspel/go-gin-template/validation"
)

type EmployeeHandler interface {
	GetEmployees(ctx *gin.Context)
	GetEmployeesWithJoin(ctx *gin.Context)
}

type handler struct {
	service service.EmployeeService
	logger  logger.Logger
}

func NewEmployeeHandler(service service.EmployeeService, logger logger.Logger) EmployeeHandler {
	return &handler{service: service, logger: logger}
}

// Get employee
// @Tags employees
// @Summary Get employees
// @Description Get employees
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.EmployeeRequestGet true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /v1/employee [GET]
func (h handler) GetEmployees(ctx *gin.Context) {
	input := new(dto.GetEmployeesRequest)

	setGetEmployeeRequest(ctx, input)

	err := validation.ValidateEmployeeGet(input)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			helper.NewResponse(
				helper.ErrorResponse{
					Code:    helper.GetErrorStatusCode(err),
					Message: err.Error(),
				},
				err,
			),
		)
		return
	}

	response, err := h.service.GetEmployees(*input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.NewResponse(response, nil))
}

func (h handler) GetEmployeesWithJoin(ctx *gin.Context) {
	input := new(dto.GetEmployeesRequest)

	setGetEmployeeRequest(ctx, input)

	err := validation.ValidateEmployeeGet(input)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			helper.NewResponse(
				helper.ErrorResponse{
					Code:    helper.GetErrorStatusCode(err),
					Message: err.Error(),
				},
				err,
			),
		)
		return
	}

	response, err := h.service.GetEmployeesWithJoin(*input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.NewResponse(response, nil))
}

func setGetEmployeeRequest(ctx *gin.Context, input *dto.GetEmployeesRequest) {
	managerId, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}
	input.ManagerID = managerId

	gender := ctx.Request.URL.Query().Get("gender")
	input.Gender = strings.ToLower(gender)

	idNumber := ctx.Request.URL.Query().Get("identityNumber")
	input.IdentityNumber = strings.ToLower(idNumber)

	name := ctx.Request.URL.Query().Get("name")
	input.Name = strings.ToLower(name)

	departmentId := ctx.Request.URL.Query().Get("departmentId")
	input.DepartmentID = strings.ToLower(departmentId)

	limitParam := ctx.Request.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 0 {
		input.Limit = dto.DefaultLimit
	} else {
		input.Limit = limit
	}

	offsetParam := ctx.Request.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		input.Offset = dto.DefaultOffset
	} else {
		input.Offset = offset
	}
}
