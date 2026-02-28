package finance

import (
	"context"
	"fmt"
	"time"

	"honda-leasing-api/internal/domain/contract"
)

type Service interface {
	GetPaymentSchedules(ctx context.Context, contractID int64) ([]PaymentScheduleResponse, error)

	// Call functions for Officer Task transitions
	GeneratePaymentSchedule(ctx context.Context, contractID int64) error
	CreatePurchaseOrder(ctx context.Context, contractID int64) error
}

type service struct {
	repo contract.FinanceRepository
}

func NewService(repo contract.FinanceRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetPaymentSchedules(ctx context.Context, contractID int64) ([]PaymentScheduleResponse, error) {
	schedules, err := s.repo.FindPaymentSchedules(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedules: %w", err)
	}

	var res []PaymentScheduleResponse
	now := time.Now()

	for _, sch := range schedules {
		lateFee := 0.0

		// Jika belum lunas dan sudah lewat jatuh tempo, hitung denda
		if sch.StatusPembayaran != "paid" && now.After(sch.JatuhTempo) {
			daysLate := int(now.Sub(sch.JatuhTempo).Hours() / 24)
			if daysLate > 0 {
				lateFee = float64(daysLate) * 2000.0 // Denda statis Rp 2.000 per hari
			}
		}

		res = append(res, PaymentScheduleResponse{
			ScheduleID:       sch.ScheduleID,
			AngsuranKe:       sch.AngsuranKe,
			JatuhTempo:       sch.JatuhTempo,
			Pokok:            sch.Pokok,
			Margin:           sch.Margin,
			LateFee:          lateFee,
			TotalTagihan:     sch.TotalTagihan + lateFee, // Total tagihan dinamis ditambah denda
			StatusPembayaran: sch.StatusPembayaran,
			TanggalBayar:     sch.TanggalBayar,
		})
	}

	return res, nil
}

func (s *service) GeneratePaymentSchedule(ctx context.Context, contractID int64) error {
	// Implement real logic for generating payment schedule.
	// For simulation, we log this execution.
	fmt.Printf("[Finance Service] Generating Payment Schedule for ContractID: %d\n", contractID)
	return nil
}

func (s *service) CreatePurchaseOrder(ctx context.Context, contractID int64) error {
	// Implement real logic for communicating with Dealer API.
	fmt.Printf("[Finance Service] Creating Purchase Order for ContractID: %d\n", contractID)
	return nil
}
