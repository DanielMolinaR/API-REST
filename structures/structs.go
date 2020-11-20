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
	User Users	`json:"data"`
	Birthdate  int 	`json:"Birthdate"`
}

type Employee struct{
	User   Users	`json:"data"`
	Active bool		`json:"Active"`
	Admin  bool		`json:"Admin"`
}

type Habits struct{
	Exercise	string	`json:"exercise"`
	Activities	string	`json:"activities"`
	Sleep		string	`json:"sleep"`
	Nutrition	string	`json:"nutrition"`
}

type Anamnesis struct{
	AntedecentesMQ	string	`json:"mq_antecedentes"`	//Antecedentes personales médico-quirúrgicos
	AntecedentesF	string	`json:"f_antecedentes"`	 	//Antedecentes familiares
	Ftreatments		string	`json:"f_treatments"` 		//Tratamientos farmacológicos
	Mtreatments		string	`json:"m_treatments"`		//Tratamientos médicos
	Habits			Habits	`json:"habits"`				//Hábitos de vida
	Treatments		string	`json:"treatments"`			//Tratamientos multidisciplinarios accesorios
}

type Exploration struct{
	Observation		string	`json:"observation"`
	PainDetails		string	`json:"pain_details"`
	Range			string	`json:"range"`				//Rango articular
	Daniels			string	`json:"daniels"`
	Deficits		string	`json:"deficits"`
}

type ClinicalBackgroundData struct{
	Reason			string		`json:"reason"`
	Anamnesis 		Anamnesis 	`json:"anamnesis"`
	Exploration		Exploration	`json:"exploration"`
	YellowFlags		string		`json:"yellow_flags"`
	Diagnosis		string		`json:"diagnosis"`		//Diagnóstico de fisioterapia
	Pobjectives		string		`json:"p_objectives"`	//Objetivos del paciente
	Tobjectives		string		`json:"t_objectives"`	//Objetivos del tratamiento
	Treatment		string		`json:"treatment"`		//Tratamiento del fisioterapeuta
	Progression		string		`json:"progression"`	//Evolución/Progreso del paciente
	Observations	string		`json:"observations"`
}