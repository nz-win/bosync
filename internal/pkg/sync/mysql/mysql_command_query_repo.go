package mysql

import (
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/types"
	"github.com/jmoiron/sqlx"
	"log"
)

type CommandQueryRepository struct {
	conn                *sqlx.DB
	insertBackordersSql string
}

func NewCommandQueryRepository(conn *sqlx.DB) *CommandQueryRepository {
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

	doDrop, err := tx.Prepare(`DELETE FROM ActiveBackorders;`)

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`REPLACE INTO ActiveBackorders (
                               updated_at,
                               business_area_no, 
                               area_no, 
                               adm_no, 
                               order_no, 
                               customer_no, 
                               material_no,
                               sale_date, 
                               material_eta_date, 
                               backorder_qty) VALUES (DEFAULT,?,?,?,?,?,?,?,?,?)`)

	if err != nil {
		return err
	}
	defer func() {
		pkg.CheckAndLogFatal(doDrop.Close())
		pkg.CheckAndLogFatal(stmt.Close())
	}()

	_, err = doDrop.Exec()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	for _, b := range records {
		_, err = stmt.Exec(b.BusinessAreaNo, b.AreaNo, b.AdmNo, b.SalesDoc, b.SoldToParty, b.Material, b.SalesDate.Date, b.MatAvailDate.Date, b.BackorderQty)

		if err != nil {
			log.Printf("Failed to insert backorder record: \n %v", b)
		}
	}

	return tx.Commit()

}
