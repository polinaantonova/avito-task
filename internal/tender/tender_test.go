package tender

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestValidation(t *testing.T) {
	t.Run("new tender", func(t *testing.T) {
		tender := NewTender()

		//подключаюсь к postgres
		user := os.Getenv("POSTGRES_USERNAME")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		dbName := os.Getenv("POSTGRES_DATABASE")

		if user == "" || password == "" || host == "" || port == "" || dbName == "" {

			//errorText := fmt.Sprintf("empty env variables\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
			//return errors.New(errorText)

			user = "cnrprod1725725190-team-78136"
			password = "cnrprod1725725190-team-78136"
			host = "rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net"
			port = "6432"
			dbName = "cnrprod1725725190-team-78136"
		}

		psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
			host, port, user, password, dbName)

		db, err := sql.Open("postgres", psqlInfo)

		if err != nil {
			errorText := fmt.Sprintf("cannot connect to database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
			t.Fatal(errorText)
		}

		defer db.Close()

		err = db.Ping()
		if err != nil {
			errorText := fmt.Sprintf("cannot ping database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
			t.Fatal(errorText)
		}

		fmt.Println("Successfully connected!")

		t.Run("unvalid service type", func(t *testing.T) {
			tender.ServiceType = "Music"
			err := tender.ValidateTenderServiceType()
			require.Error(t, err)
		})

		t.Run("too long string field", func(t *testing.T) {
			tender.Name = "a_very_very_loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong_name"
			err := tender.ValidateStringFieldsLen()
			require.Error(t, err)
		})

		t.Run("username validation error", func(t *testing.T) {
			tender.CreatorUsername = "user100"
			err := tender.ValidateUser(db)
			require.Error(t, err)
		})

		t.Run("username validation ok", func(t *testing.T) {
			tender.CreatorUsername = "user1"
			err := tender.ValidateUser(db)
			require.Nil(t, err)
		})

		t.Run("username validation empty", func(t *testing.T) {
			tender.CreatorUsername = ""
			err := tender.ValidateUser(db)
			require.Error(t, err)
		})
	})

}
