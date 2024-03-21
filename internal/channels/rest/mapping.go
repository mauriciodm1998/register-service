package rest

import (
	"register-service/internal/domain"
)

func toResponse2(weekReports []domain.DailyRegister) []ResponseAppointments {
	result := []ResponseAppointments{}

	for _, register := range weekReports {
		result = append(result, toResponse(register))
	}

	return result
}

func toResponse(daily domain.DailyRegister) ResponseAppointments {
	var appointments []Appointments
	var date string

	if len(daily.Clocks) > 0 {
		date = daily.Clocks[0].Date.Format("02/01/2006")
	}

	for _, value := range daily.Clocks {
		appointment := Appointments{
			Time: value.Time.Format("15:04"),
		}
		appointments = append(appointments, appointment)
	}

	return ResponseAppointments{
		ClockIns:    appointments,
		HoursWorked: daily.Hours,
		Date:        date,
	}
}
