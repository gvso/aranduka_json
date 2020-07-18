package main

type record struct {
	Expenses []*expense `json:"egresos"`
	Incomes  []*income  `json:"ingresos"`
}

type expense struct {
	Period            string `json:"periodo"`
	RUC               string `json:"ruc"`
	DocumentType      string `json:"tipo"`
	DocumentTypeText  string `json:"tipoTexto"`
	Date              string `json:"fecha"`
	TimbradoNumber    string `json:"timbradoNumero"`
	TimbradoDocument  string `json:"timbradoDocumento"`
	TimbradoCondition string `json:"timbradoCondicion"`
	EntityIDType      string `json:"relacionadoTipoIdentificacion"`
	EntityID          string `json:"relacionadoNumeroIdentificacion"`
	Entity            string `json:"relacionadoNombres"`
	Amount            int64  `json:"egresoMontoTotal"`
	ExpenseType       string `json:"tipoEgreso"`
	ExpenseSubtype    string `json:"subtipoEgreso"`
}

type income struct{}

func documentTypeID(value string) string {
	switch value {
	case "Factura":
		return "1"
	default:
		return "0"
	}
}

func expenseType(value string) string {
	switch value {
	case "Gasto":
		return "gasto"
	default:
		return ""
	}
}

func expenseSubtype(value string) string {
	switch value {
	case "Gastos personales y de familiares a cargo realizados en el país":
		return "GPERS"
	default:
		return ""
	}
}
