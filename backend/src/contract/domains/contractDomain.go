package domains

import (
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"errors"
)

type UseCase interface {
	GetContracts(filter requests.ContractFilter) ([]responses.ContractResponse, error)
	GetContractByID(id int) (*responses.ContractResponse, error)
	CreateContract(userID int, userRole string, req requests.CreateContractRequest) (*responses.ContractResponse, error)
	ApproveContract(id int, userID int) (*responses.ContractResponse, error)
	RejectContract(id int, userID int, reason string) error
	TerminateContract(id int, userID int) error
}

var ErrApplicationNotApproved = errors.New("application is not approved for this stall")

type Repository interface {
	FindAll(filter requests.ContractFilter) ([]models.Contract, error)
	FindByID(id int) (*models.Contract, error)
	Create(contract *models.Contract) error
	Approve(contractID int, userID int) (*models.Contract, error)
	Reject(contractID int, userID int, reason string) error
	UpdateStallStatus(stallIdent string, status string) error
	Terminate(contractID int) error
}
