package service

import (
	"context"
	"fmt"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/integration/mail"
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
	GetMonthAppointments(ctx context.Context, userId int, userEmail string) error
	ReportAppointments(ctx context.Context, report domain.MonthReportRequest) error
}

type registerService struct {
	publisher  sqs_publisher.Publisher
	repository repository.Repository
	mailer     mail.Mailer
}

func NewRegisterService() RegisterService {
	return &registerService{
		publisher:  sqs_publisher.NewSQS(),
		repository: repository.New(),
		mailer:     mail.NewMailer(),
	}
}

func (s *registerService) ClockIn(ctx context.Context, userId int) error {
	if err := s.publisher.SendMessage(
		domain.ClockInRegister{
			Id:        uuid.New().String(),
			UserId:    userId,
			Date:      time.Date(2024, time.March, 22, 0, 0, 0, 0, time.UTC),
			Time:      time.Date(2024, time.March, 22, 9, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		// domain.ClockInRegister{
		// 	Id:        uuid.New().String(),
		// 	UserId:    userId,
		// 	Date:      time.Now().Truncate(24 * time.Hour).UTC(),
		// 	Time:      time.Now().UTC(),
		// 	CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
		// },
		config.Get().SQS.ClockInQueue,
	); err != nil {
		return fmt.Errorf("an error occurred when sending clock in to queue: %w", err)
	}

	return nil
}

func (s *registerService) QueueToDatabase(ctx context.Context, register domain.ClockInRegister) error {
	register.Id = uuid.NewString()

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

func (s *registerService) GetWeekAppointments(ctx context.Context, userId int) ([]domain.DailyRegister, error) {
	appointments, err := s.repository.GetWeekAppointments(ctx, userId, time.Now().Truncate(24*time.Hour).UTC())
	if err != nil {
		return nil, err
	}

	return s.mountDailyRegister(s.mergeIntoDays(appointments)), nil
}

func (s *registerService) GetMonthAppointments(ctx context.Context, userId int, userEmail string) error {
	report := domain.MonthReportRequest{
		UserId: userId,
		Time:   time.Now(),
		Mail:   userEmail,
	}

	if err := s.publisher.SendMessage(report, config.Get().SQS.ReportQueue); err != nil {
		return fmt.Errorf("an error occurred when sending report to queue: %w", err)
	}

	return nil
}

func (s *registerService) ReportAppointments(ctx context.Context, reportRequest domain.MonthReportRequest) error {
	if s.mailer == nil {
		return fmt.Errorf("mailer is not set")
	}

	report, err := s.repository.GetMonthAppointments(ctx, reportRequest.UserId, reportRequest.Time)
	if err != nil {
		return err
	}

	monthRegisters := s.mountDailyRegister(s.mergeIntoDays(report))

	//report, err := s.repository.SaveReport(ctx, appointments)
	//if err != nil {
	//	return err
	//}

	message, err := s.mailer.MountHTMLBody(struct {
		ID      int
		Time    string
		Message string
	}{
		ID:      reportRequest.UserId,
		Time:    fmt.Sprintf("%d hours", getMonthTotalTime(monthRegisters)),
		Message: mountMonthReport(monthRegisters),
	})
	if err != nil {
		return err
	}

	if err = s.mailer.SendEmail("Month Clock In Report", message, []string{reportRequest.Mail}, nil, nil, nil); err != nil {
		return err
	}

	return nil
}

func (*registerService) mountDailyRegister(periodAppointments map[string][]domain.ClockInRegister) []domain.DailyRegister {
	var periodRegister []domain.DailyRegister

	for _, appointments := range periodAppointments {

		dailyRegister := domain.DailyRegister{
			Clocks: appointments,
			Hours:  calculeHoursWorked(appointments),
		}

		periodRegister = append(periodRegister, dailyRegister)
	}

	return periodRegister
}

func (*registerService) mergeIntoDays(appointments []domain.ClockInRegister) map[string][]domain.ClockInRegister {
	periodAppointments := map[string][]domain.ClockInRegister{}

	for _, appointment := range appointments {

		if existingAppointments, ok := periodAppointments[appointment.Date.String()]; ok {
			periodAppointments[appointment.Date.String()] = append(existingAppointments, appointment)
		} else {
			periodAppointments[appointment.Date.String()] = []domain.ClockInRegister{appointment}
		}
	}

	return periodAppointments
}

func mountMonthReport(registers []domain.DailyRegister) string {
	var report string

	for _, register := range registers {
		report += "<tr>"

		var clocks string

		for _, clock := range register.Clocks {
			clocks += fmt.Sprintf(" (%s) ", clock.Time.String())
		}

		report += fmt.Sprintf("<td>%s</td><td>%s</td><td>%d</td>\n", register.Clocks[0].Date.String(), clocks, register.Hours)
		report += "</tr>"
	}

	return report
}

func getMonthTotalTime(registers []domain.DailyRegister) int {
	var totalTime int

	for _, register := range registers {
		totalTime += register.Hours
	}

	return totalTime
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
