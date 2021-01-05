package structures

type Users struct{
	DNI         string 	`json:"DNI"`
	Password	string	`json:"Password"`
	Email       string 	`json:"Email"`
	Name 		string 	`json:"Name"`
	Surname 	string 	`json:"Surname"`
	Phone		string 	`json:"Phone"`
}

type Patient struct{
	User Users	`json:"Data"`
	Birthdate  string 	`json:"Birthdate"`
}

type Employee struct{
	User   Users	`json:"Data"`
	Active bool		`json:"Active"`
	Admin  bool		`json:"Admin"`
}

type Habits struct{
	Exercise	string	`json:"Exercise"`
	Activities	string	`json:"Activities"`
	Sleep		string	`json:"Sleep"`
	Nutrition	string	`json:"Nutrition"`
}

type Anamnesis struct{
	AntedecentesMQ	string	`json:"Mq_antecedentes"`	//Antecedentes personales médico-quirúrgicos
	AntecedentesF	string	`json:"F_antecedentes"`	 	//Antedecentes familiares
	Ftreatments		string	`json:"F_treatments"` 		//Tratamientos farmacológicos
	Mtreatments		string	`json:"M_treatments"`		//Tratamientos médicos
	Habits			Habits	`json:"Habits"`				//Hábitos de vida
	Treatments		string	`json:"Treatments"`			//Tratamientos multidisciplinarios accesorios
}

type Exploration struct{
	Observation		string	`json:"Observation"`
	PainDetails		string	`json:"Pain_details"`
	Range			string	`json:"Range"`				//Rango articular
	Daniels			string	`json:"Daniels"`
	Deficits		string	`json:"Deficits"`
}

type ClinicalBackgroundData struct{
	Patient_dni 	string 		`json:"Patient_dni"`
	Reason			string		`json:"Reason"`
	Anamnesis 		Anamnesis 	`json:"Anamnesis"`
	Exploration		Exploration	`json:"Exploration"`
	YellowFlags		string		`json:"Yellow_flags"`
	Diagnosis		string		`json:"Diagnosis"`		//Diagnóstico de fisioterapia
	P_objectives	string		`json:"P_objectives"`	//Objetivos del paciente
	T_objectives	string		`json:"T_objectives"`	//Objetivos del tratamiento
	Treatment		string		`json:"Treatment"`		//Tratamiento del fisioterapeuta
	Progression		string		`json:"Progression"`	//Evolución/Progreso del paciente
	Observations	string		`json:"Observations"`
}

type Appointment struct{
	New_user		bool	`json:"New_user"`		//if false phone and name must to be empty
	Patient_phone	string	`json:"Patient_phone"`
	Patient_name	string	`json:"Patient_name"`
	Patient_email 	string 	`json:"Patient_email"`
	Employee_email	string	`json:"Employee_email"`	//This never going to be empty because when the user log in we send back the email so here we recieve the email from the store
	Hour			int		`json:"Hour"`
	Minute			int		`json:"Minute"`
	Day				int 	`json:"Day"`
	Month			int		`json:"Month"`
	Year			int		`json:"Year"`
}

type Exercise struct{
	Exercise_name	string	`json:"Exercise_name"`
	Description		string 	`json:"Description"`
	Patient_email 	string 	`json:"Patient_email"`
	Hour			int		`json:"Hour"`
	Minute			int		`json:"Minute"`
	Day				int 	`json:"Day"`
	Month			int		`json:"Month"`
	Year			int		`json:"Year"`
}

type AppointmentResponse struct{
	Date 			string
	Employee_name 	string
	Patient_name	string
}

type ExercisesResponse struct{
	Date 			string
	Description 	string
	Exercise_name	string
}

type Date struct{
	Date 			string
}