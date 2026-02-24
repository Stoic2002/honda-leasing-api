package finance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
)

type Service interface {
	GetPaymentSchedules(ctx context.Context, contractID int64) ([]PaymentScheduleResponse, error)
	ProcessPaymentWebhook(ctx context.Context, req WebhookPaymentRequest) error

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

func (s *service) ProcessPaymentWebhook(ctx context.Context, req WebhookPaymentRequest) error {
	// 1. Validasi schedule
	schedule, err := s.repo.FindScheduleByID(ctx, req.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to find schedule: %w", err)
	}
	if schedule == nil {
		return errors.New("payment schedule not found")
	}

	if schedule.StatusPembayaran == "paid" {
		return errors.New("payment schedule is already paid")
	}

	// 2. Cek apakah invoice sesuai
	// (Di dunia nyata harus cek jumlah vs tagihan + late fee, untuk simplifikasi kita anggap lunas)

	// 3. Update status
	now := time.Now()
	schedule.StatusPembayaran = "paid"
	schedule.TanggalBayar = &now

	// 4. Create payment record
	payment := &entity.Payment{
		NomorBukti:       req.NomorBukti,
		JumlahBayar:      req.JumlahBayar,
		TanggalBayar:     now,
		MetodePembayaran: req.MetodePembayaran,
		ContractID:       req.ContractID,
		ScheduleID:       &req.ScheduleID,
		CreatedAt:        now,
	}

	if req.Provider != "" {
		payment.Provider = &req.Provider
	}

	return s.repo.CreatePaymentAndUpdateSchedule(ctx, payment, schedule)
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
