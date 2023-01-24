package employee

type Employee struct {
	EmployeeName string  `json:"employeeName"`
	Department   string  `json:"department"`
	Salary       float64 `json:"salary"`
	Address      string  `json:"address"`
}
