package middleware

import (
	"TFG/structures"
	"github.com/robfig/cron/v3"
	"strconv"
)

var c *cron.Cron

func init(){
	c = cron.New()

	c.Start()
	
	c.AddFunc("0 9 * * 1-5", func(){ sendReminderEmailToEmployees()})
}

func setAppointmentReminder(appointmentData structures.Appointment) {
	c.AddFunc(strconv.Itoa(appointmentData.Minute) + " " + strconv.Itoa(appointmentData.Hour - 1) + " " + strconv.Itoa(appointmentData.Day) + " " +
		strconv.Itoa(appointmentData.Month) + " *", func(){appointmentReminder(appointmentData)})
}

func setExerciseReminder(exerciseData structures.Exercise) {
	c.AddFunc(strconv.Itoa(exerciseData.Minute) + " " + strconv.Itoa(exerciseData.Hour - 1) + " " + strconv.Itoa(exerciseData.Day) + " " +
		strconv.Itoa(exerciseData.Month) + " *", func(){exerciseRreminder(exerciseData)})
}

func appointmentReminder(appointmentData structures.Appointment) {
	minute := strconv.Itoa(appointmentData.Minute)
	if (appointmentData.Minute >= 0 && appointmentData.Minute<10){
		minute = "0" + minute
	}
	_, employee_dni := getStringFromField("employee", "dni", "email", appointmentData.Employee_email)
	_, employee_name := getStringFromField("employee", "name", "dni", employee_dni)
	sendReminder("CITA FISIOTERAPIA", "Tienes una cita pendiente con " + employee_name + ".",
		strconv.Itoa(appointmentData.Day), strconv.Itoa(appointmentData.Hour) + ":" + minute,
		"http://localhost:8081/calendar", appointmentData.Patient_email, appointmentData.Month)

}

func exerciseRreminder(exerciseData structures.Exercise) {
	minute := strconv.Itoa(exerciseData.Minute)
	if (exerciseData.Minute >= 0 && exerciseData.Minute<10){
		minute = "0" + minute
	}
	sendReminder("EJERCICIO: " + exerciseData.Exercise_name + " pendiente", exerciseData.Description,
		strconv.Itoa(exerciseData.Day), strconv.Itoa(exerciseData.Hour) + ":" + minute,
		"http://localhost:8081/calendar", exerciseData.Patient_email, exerciseData.Month)

}

func sendReminderEmailToEmployees() {
	var employeeDni, email string

	_, rows := getAllEmployeeDnis()

	for (rows.Next()){
		rows.Scan(&employeeDni, &email)
		_, appointmentsRows := getAllApointmentsOfTheDay(employeeDni)
		sendDailyReminder(appointmentsRows, email)
	}
}