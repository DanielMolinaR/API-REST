package middleware

import (
	"TFG/API-REST/structures"
	"github.com/robfig/cron/v3"
	"strconv"
)

var c *cron.Cron

func init(){
	c = cron.New()

	c.Start()
	
	c.AddFunc("0 9 * * 1-5", func(){ sendReminderEmailToEmployees()})
}

func setReminder(appointmentData structures.Appointment, isAppointment bool) {
	c.AddFunc(strconv.Itoa(appointmentData.Minute) + " " + strconv.Itoa(appointmentData.Hour - 1) + " " + strconv.Itoa(appointmentData.Day) + " " +
		strconv.Itoa(appointmentData.Month) + " *", func(){reminder(appointmentData, isAppointment)})
}

func reminder(appointmentData structures.Appointment, isAppointment bool) {
	if isAppointment{
		minute := strconv.Itoa(appointmentData.Minute)
		if (appointmentData.Minute >= 0 && appointmentData.Minute<10){
			minute = "0" + minute
		}
		_, employee_dni := getStringFromField("employee", "dni", "email", appointmentData.Employee_email)
		_, employee_name := getStringFromField("employee", "name", "dni", employee_dni)
		sendReminder("Cita fisioterapia", "Tienes una cita pendiente con " + employee_name + ".",
			string(appointmentData.Day), string(appointmentData.Hour) + ":" + minute,
			"http://localhost:8081/calendar", appointmentData.Patient_email, appointmentData.Month)
	} /*else {
		sendReminder("Ejercicio fisioterapéutico", exerciseDescription,
			string(appointmentData.Day), string(appointmentData.Hour) + ":" + string(appointmentData.Minute),
			appointmentData.Patient_name, "http://localhost:8081/calendar", appointmentData.Patient_email)
	}*/
}