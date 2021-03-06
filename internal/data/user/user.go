package user

import (
	"context"
	"log"

	"go-tutorial-2020/pkg/errors"

	userEntity "go-tutorial-2020/internal/entity/user"

	"github.com/jmoiron/sqlx"
)

type (
	// Data ...
	Data struct {
		db   *sqlx.DB
		stmt map[string]*sqlx.Stmt
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

const (
	getAllUsers  = "GetAllUsers"
	qGetAllUsers = "SELECT * FROM user_test"

	insertAllUsers  = "insertAllUsers"
	qinsertAllUsers = "INSERT INTO user_test VALUES (?,?,?,?,?,?)"

	updateAllUsers  = "updateAllUsers"
	qupdateAllUsers = "UPDATE user_test set nama_lengkap = ?, tanggal_lahir = ?, jabatan = ?, email = ? WHERE nip = ?"

	deleteAllUsers  = "deleteAllUsers"
	qdeleteAllUsers = "Delete FROM user_test where nip = ?"
)

var (
	readStmt = []statement{
		{getAllUsers, qGetAllUsers},
		{insertAllUsers, qinsertAllUsers},

		{updateAllUsers, qupdateAllUsers},
		{deleteAllUsers, qdeleteAllUsers},
	}
)

// New ...
func New(db *sqlx.DB) Data {
	d := Data{
		db: db,
	}

	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// GetAllUsers digunakan untuk mengambil semua data user
func (d Data) GetAllUsers(ctx context.Context) ([]userEntity.User, error) {
	var (
		user  userEntity.User
		users []userEntity.User
		err   error
	)

	// Query ke database
	rows, err := d.stmt[getAllUsers].QueryxContext(ctx)

	// Looping seluruh row data
	for rows.Next() {
		// Insert row data ke struct user
		if err := rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA][GetAllUsers] ")
		}
		// Tambahkan struct user ke array user
		users = append(users, user)
	}
	// Return users array
	return users, err
}

// InsertAllUsers
func (d Data) InsertAllUsers(ctx context.Context, user userEntity.User) error {
	_, err := d.stmt[insertAllUsers].ExecContext(ctx,
		user.UserID,
		user.UserNip,
		user.UserNama,
		user.UserTanggalLahir,
		user.UserJabatan,
		user.UserAlamatEmail)
	return err
}

func (d Data) UpdateAllUsers(ctx context.Context, user userEntity.User) error {

	_, err := d.stmt[updateAllUsers].ExecContext(ctx,
		user.UserNama,
		user.UserTanggalLahir,
		user.UserJabatan,
		user.UserAlamatEmail,
		user.UserNip)

	return err
}

func (d Data) DeleteAllUsers(ctx context.Context, user userEntity.User) error {

	_, err := d.stmt[deleteAllUsers].ExecContext(ctx,
		user.UserNip)

	return err
}
