package mysql

import (
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/types"
	"database/sql"
	"log"
)

type CommandQueryRepository struct {
	conn                *sql.DB
	insertBackordersSql string
}

func NewCommandQueryRepository(conn *sql.DB) *CommandQueryRepository {
	return &CommandQueryRepository{conn: conn}
}

func (cq *CommandQueryRepository) Close() error {
	return cq.conn.Close()
}

func (cq *CommandQueryRepository) UpdateBackOrders(records []types.BackOrder) error {
	tx, err := cq.conn.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`REPLACE INTO ActiveBackorders (
                               business_area_no, 
                               area_no, 
                               adm_no, 
                               order_no, 
                               customer_no, 
                               material_no,
                               sale_date, 
                               material_eta_date, 
                               backorder_qty) VALUES (?,?,?,?,?,?,?,?,?)`)

	if err != nil {
		return err
	}
	defer func() {
		pkg.CheckAndLogFatal(stmt.Close())
	}()

	for _, b := range records {
		_, err = stmt.Exec(b.BusinessAreaNo, b.AreaNo, b.AdmNo, b.SalesDoc, b.SoldToParty, b.Material, b.SalesDate.Date, b.MatAvailDate.Date, b.BackorderQty)

		if err != nil {
			log.Printf("Failed to insert backorder record: \n %v", b)
		}
	}

	return tx.Commit()

}
