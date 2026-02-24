package leasing

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
	"honda-leasing-api/internal/domain/vo"
	"honda-leasing-api/pkg/pagination"
)

const (
	// DefaultMarginRate is the default interest margin for loan calculations.
	DefaultMarginRate = 1.2
	// ContractCodeRandomBound is used for generating unique contract numbers.
	ContractCodeRandomBound = 1000
)

// SubmitOrderInput is a clean service-level input struct (no HTTP tags).
type SubmitOrderInput struct {
	UserID         int64
	MotorID        int64
	ProductID      int64
	NilaiKendaraan float64
	DpDibayar      float64
	TenorBulan     int16
}

type Service interface {
	SubmitOrder(ctx context.Context, req SubmitOrderInput) (*entity.LeasingContract, error)
	GetMyOrders(ctx context.Context, userID int64, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error)
	GetContractProgress(ctx context.Context, contractID int64) ([]entity.LeasingTask, error)
}

type service struct {
	repo contract.LeasingRepository
}

func NewService(repo contract.LeasingRepository) Service {
	return &service{repo: repo}
}

func (s *service) SubmitOrder(ctx context.Context, req SubmitOrderInput) (*entity.LeasingContract, error) {
	// 1. Resolve customer_id dari user_id yang ada di JWT
	customer, err := s.repo.FindCustomerByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: customer profile not found, please complete your profile first", domain.ErrNotFound)
	}

	pokokPinjaman := req.NilaiKendaraan - req.DpDibayar
	if pokokPinjaman <= 0 {
		return nil, fmt.Errorf("%w: invalid DP amount", domain.ErrInvalidInput)
	}

	totalPinjaman := pokokPinjaman * DefaultMarginRate
	cicilan := totalPinjaman / float64(req.TenorBulan)

	contractCode := fmt.Sprintf("CTR-%d-%04d", time.Now().Unix(), rand.Intn(ContractCodeRandomBound))

	newContract := &entity.LeasingContract{
		ContractNumber:  contractCode,
		RequestDate:     time.Now(),
		TenorBulan:      req.TenorBulan,
		NilaiKendaraan:  req.NilaiKendaraan,
		DpDibayar:       req.DpDibayar,
		PokokPinjaman:   pokokPinjaman,
		TotalPinjaman:   totalPinjaman,
		CicilanPerBulan: cicilan,
		Status:          vo.StatusDraft.String(),
		CustomerID:      customer.CustomerID, // ← dari DB, bukan dari request
		MotorID:         req.MotorID,
		ProductID:       req.ProductID,
	}

	// Fetch template tasks to generate the checklist
	templates, err := s.repo.GetTemplateTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch template tasks", domain.ErrInternalServerError)
	}

	var leasingTasks []entity.LeasingTask
	for _, t := range templates {
		status := vo.StatusPending.String()
		var actualStart *time.Time

		// Set the very first task to in-progress
		if t.SequenceNo == 1 {
			status = vo.StatusInProgress.String()
			now := time.Now()
			actualStart = &now
		}

		tetaID := t.TetaID
		leasingTasks = append(leasingTasks, entity.LeasingTask{
			TaskName:        t.TetaName,
			TemplateTaskID:  &tetaID,
			Status:          status,
			RoleID:          t.TetaRoleID,
			SequenceNo:      t.SequenceNo,
			CallFunction:    t.CallFunction,
			ActualStartdate: actualStart,
		})
	}

	// Persist to DB
	err = s.repo.CreateContract(ctx, newContract, leasingTasks)
	if err != nil {
		return nil, err
	}

	return newContract, nil
}

func (s *service) GetMyOrders(ctx context.Context, userID int64, pg contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)
	return s.repo.FindContractsByUserID(ctx, userID, pg)
}

func (s *service) GetContractProgress(ctx context.Context, contractID int64) ([]entity.LeasingTask, error) {
	return s.repo.GetContractProgress(ctx, contractID)
}
