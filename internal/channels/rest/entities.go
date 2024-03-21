package rest

type Response struct {
	Message string `json:"message"`
}

type ResponseAppointments struct {
	Date        string         `json:"date"`
	HoursWorked int            `json:"hours_worked"`
	ClockIns    []Appointments `json:"appointments"`
}

type Appointments struct {
	Time string `json:"time"`
}
