package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDb(t *testing.T) *sql.DB {

	conn, err := sql.Open("pgx", "host=localhost user=postgres dbname=bluehorse password= postgres sslmode=disable port=5432")
	if err != nil {
		t.Fatalf("error opening testDB: %v", err)
	}

	err = Migrate(conn, "../../migration")
	if err != nil {
		t.Fatalf("error migrating: %v", err)
	}

	_, err = conn.Exec(`TRUNCATE lawyer , addresses CASCADE`)
	if err != nil {
		t.Fatalf("error executing: %v", err)
	}
	return conn
}

func TestCreateLawyer(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()
	store := NewPostgresDB(db)
	tests := []struct {
		name    string
		lawyer  *Lawyer
		wantErr bool
	}{
		{
			name: "valid lawyer",
			lawyer: &Lawyer{
				First_name:   "nati",
				Last_name:    "abebe",
				Email:        "nati@369",
				Phone_number: "0983315117"},
			wantErr: false,
		},
		{
			name: "invalid lawyer",
			lawyer: &Lawyer{
				Last_name:    "abebe",
				Email:        "nati@369",
				Phone_number: "0983315117"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdLawyer, err := store.CreateLawyer(tt.lawyer)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.lawyer.First_name, createdLawyer.First_name)
			assert.Equal(t, tt.lawyer.Last_name, createdLawyer.Last_name)
			assert.Equal(t, tt.lawyer.Email, createdLawyer.Email)
			assert.Equal(t, tt.lawyer.Phone_number, createdLawyer.Phone_number)
			retreivedLawyer, err := store.GetLawyerById(int64(createdLawyer.ID))
			require.NoError(t, err)
			assert.Equal(t, retreivedLawyer.ID, createdLawyer.ID)
		})

	}
}

func IntPtr(I int) *int {
	return &I
}

func FloatPrt(F float64) *float64 {
	return &F
}
