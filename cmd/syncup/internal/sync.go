package internal

import "database/sql"

func UpdateMysql(db *sql.DB, records []BackOrder) {
	tx, err := db.Begin()
	checkAndPanic(err)

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
	checkAndPanic(err)
	defer stmt.Close()

	for _, b := range records {
		_, err = stmt.Exec(b.BusinessAreaNo, b.AreaNo, b.AdmNo, b.SalesDoc, b.SoldToParty, b.Material, b.SalesDate.Date, b.MatAvailDate.Date, b.BackorderQty)
		checkAndPanic(err)
	}

	err = tx.Commit()
	checkAndPanic(err)
}

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
