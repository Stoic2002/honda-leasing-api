package finance

import (
	"context"
	"fmt"
	"time"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
)

type Service interface {
	GetPaymentSchedules(ctx context.Context, contractID int64) ([]PaymentScheduleResponse, error)
	ProcessPayment(ctx context.Context, req PaymentRequest) error

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

func (s *service) ProcessPayment(ctx context.Context, req PaymentRequest) error {
	// 1. Validasi schedule
	schedule, err := s.repo.FindScheduleByID(ctx, req.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to find schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("%w: payment schedule not found", domain.ErrNotFound)
	}

	if schedule.StatusPembayaran == "paid" {
		return fmt.Errorf("%w: payment schedule is already paid", domain.ErrInvalidInput)
	}

	if schedule.ContractID != req.ContractID {
		return fmt.Errorf("%w: schedule does not belong to the specified contract", domain.ErrInvalidInput)
	}

	// 2. Update status
	now := time.Now()
	schedule.StatusPembayaran = "paid"
	schedule.TanggalBayar = &now

	// 3. Create payment record
	payment := &entity.Payment{
		NomorBukti:       req.NomorBukti,
		JumlahBayar:      req.JumlahBayar,
		TanggalBayar:     now,
		MetodePembayaran: req.MetodePembayaran,
		ContractID:       req.ContractID,
		ScheduleID:       &req.ScheduleID,
		CreatedAt:        now,
	}

	return s.repo.CreatePaymentAndUpdateSchedule(ctx, payment, schedule)
}

func (s *service) GeneratePaymentSchedule(ctx context.Context, contractID int64) error {
	// Ambil data contract untuk tahu tenor dan nominal cicilan
	contract, err := s.repo.FindContractByID(ctx, contractID)
	if err != nil {
		return fmt.Errorf("failed fetched contract: %w", err)
	}
	if contract == nil {
		return fmt.Errorf("%w: contract not found", domain.ErrNotFound)
	}

	var schedules []entity.PaymentSchedule
	now := time.Now()

	// Generate schedule per bulan berdasarkan tenor
	for i := int16(1); i <= contract.TenorBulan; i++ {
		jatuhTempo := time.Date(now.Year(), now.Month()+time.Month(i), now.Day(), 0, 0, 0, 0, time.Local)

		schedules = append(schedules, entity.PaymentSchedule{
			ContractID:       contract.ContractID,
			AngsuranKe:       i,
			JatuhTempo:       jatuhTempo,
			Pokok:            contract.PokokPinjaman / float64(contract.TenorBulan),
			Margin:           (contract.TotalPinjaman - contract.PokokPinjaman) / float64(contract.TenorBulan),
			TotalTagihan:     contract.CicilanPerBulan,
			StatusPembayaran: "unpaid",
			CreatedAt:        now,
		})
	}

	err = s.repo.CreatePaymentSchedules(ctx, schedules)
	if err != nil {
		return fmt.Errorf("failed to create schedules: %w", err)
	}

	fmt.Printf("[Finance Service] Successfully generated %d Payment Schedules for ContractID: %d\n", contract.TenorBulan, contractID)
	return nil
}

func (s *service) CreatePurchaseOrder(ctx context.Context, contractID int64) error {
	// Implement real logic for communicating with Dealer API.
	fmt.Printf("[Finance Service] Creating Purchase Order for ContractID: %d\n", contractID)
	return nil
}
