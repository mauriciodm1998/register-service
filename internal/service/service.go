package service

import (
	"context"
	"fmt"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/integration/sqs_publisher"
	"register-service/internal/repository"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
)

type RegisterService interface {
	QueueToDatabase(ctx context.Context, register domain.ClockInRegister) error
	ClockIn(ctx context.Context, userId int) error
	GetDayAppointments(ctx context.Context, userId int) (*domain.DailyRegister, error)
	GetWeekAppointments(ctx context.Context, userId int) ([]domain.DailyRegister, error)
	GetMonthAppointments(ctx context.Context, userId int) ([]domain.ClockInRegister, error)
}

type registerService struct {
	publisher  sqs_publisher.Publisher
	repository repository.Repository
}

func NewRegisterService() RegisterService {
	return &registerService{
		publisher:  sqs_publisher.NewSQS(),
		repository: repository.New(),
	}
}

func (s *registerService) ClockIn(ctx context.Context, userId int) error {
	if err := s.publisher.SendMessage(
		domain.ClockInRegister{
			Id:        uuid.New().String(),
			UserId:    userId,
			Date:      time.Date(2024, time.March, 20, 0, 0, 0, 0, time.UTC),
			Time:      time.Date(2024, time.March, 20, 18, 15, 0, 0, time.UTC),
			CreatedAt: time.Now(),
		},
		// domain.ClockInRegister{
		// 	Id:        uuid.New().String(),
		// 	UserId:    userId,
		// 	Date:      time.Now().Truncate(24 * time.Hour).UTC(),
		// 	Time:      time.Now().UTC(),
		// 	CreatedAt: time.Now(),
		// },
		config.Get().SQS.ClockInQueue,
	); err != nil {
		return fmt.Errorf("an error occurred when creating a clock in register: %w", err)
	}

	return nil
}

func (s *registerService) QueueToDatabase(ctx context.Context, register domain.ClockInRegister) error {
	err := s.repository.Create(ctx, register)
	if err != nil {
		log.Err(err).Msg("an error occurred when save the registry")
	}

	return nil
}

func (s *registerService) GetDayAppointments(ctx context.Context, userId int) (*domain.DailyRegister, error) {
	appointments, err := s.repository.GetDayAppointments(ctx, userId, time.Now().Truncate(24*time.Hour).UTC())
	if err != nil {
		return nil, err
	}

	return &domain.DailyRegister{
		Clocks: appointments,
		Hours:  calculeHoursWorked(appointments),
	}, nil
}

func calculeHoursWorked(appointments []domain.ClockInRegister) int {
	if appointmentsAreOdd(appointments) {
		return 0
	}

	var hoursWorked time.Duration

	for i := 0; i < len(appointments); i += 2 {
		hoursWorked += appointments[i+1].Time.Sub(appointments[i].Time)
	}

	return int(hoursWorked.Hours())
}

func appointmentsAreOdd(appointments []domain.ClockInRegister) bool {
	return len(appointments)%2 == 1
}

func (s *registerService) GetWeekAppointments(ctx context.Context, userId int) ([]domain.DailyRegister, error) {
	appointments, err := s.repository.GetWeekAppointments(ctx, userId, time.Now().Truncate(24*time.Hour).UTC())
	if err != nil {
		return nil, err
	}

	weekAppointments := map[string][]domain.ClockInRegister{}
	for _, appointment := range appointments {

		if existingAppointments, ok := weekAppointments[appointment.Date.String()]; ok {
			weekAppointments[appointment.Date.String()] = append(existingAppointments, appointment)
		} else {
			weekAppointments[appointment.Date.String()] = []domain.ClockInRegister{appointment}
		}
	}

	var weekRegisters []domain.DailyRegister
	for _, appointments := range weekAppointments {

		dailyRegister := domain.DailyRegister{
			Clocks: appointments,
			Hours:  calculeHoursWorked(appointments),
		}

		weekRegisters = append(weekRegisters, dailyRegister)
	}

	return weekRegisters, nil
}

func (s *registerService) GetMonthAppointments(ctx context.Context, userId int) ([]domain.ClockInRegister, error) {
	return s.repository.GetMonthAppointments(ctx, userId, time.Now())
}
