package domain

type Auditorium struct {
	AuditoriumID int64  `db:"auditorium_id" json:"auditoriumId"`
	CantRows     int    `db:"cant_rows" json:"cantRows"`
	CantCols     int    `db:"cant_cols" json:"cantCols"`
	Name         string `db:"name" json:"name"`
}

func NewAuditorium(cantRows, cantCols int, name string) *Auditorium {
	return &Auditorium{
		CantRows: cantRows,
		CantCols: cantCols,
		Name:     name,
	}
}
