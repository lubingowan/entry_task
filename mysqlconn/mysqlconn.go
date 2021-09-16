package mysqlconn

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Profile struct {
    Username string `db:"username"`
    Nickname string `db:"nickname"`
    Password string `db:"password"`
	Picture  string `db:"picture"`
}

const dbconfig string = "debian-sys-maint:v1JyHtzDAaxVUWYY@tcp(localhost:3306)/my_test";

func QueryProfile(username string) (Profile, error) {
	db, err := sqlx.Open("mysql", dbconfig)
	if err != nil {
		fmt.Println("connect database failed, ", err)
		return Profile{}, err;
	}

    defer db.Close()

	var profile []Profile
    err = db.Select(&profile, "select username, nickname, password, picture from profile where username=?", username)

	if err != nil {
        fmt.Println("QueryProfile exec failed, ", err)
        return Profile{}, err
    }

	if len(profile) > 0 {
		return profile[0], nil;
	}
	return Profile{}, nil;
}


func UpdateNickname(username string, nickname string) (error) {
	db, err := sqlx.Open("mysql", dbconfig)
	if err != nil {
		fmt.Println("connect database failed, ", err)
		return err;
	}

    defer db.Close()

    result, err := db.Exec("update profile set nickname=? where username=?", nickname, username)

	if err != nil {
        fmt.Println("exec failed, ", err)
        return err
    }

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return err
	}

	fmt.Println(rows)
	return nil;
}


func UpdatePicture(username string, picture string) (error) {
	db, err := sqlx.Open("mysql", dbconfig)
	if err != nil {
		fmt.Println("connect database failed, ", err)
		return err;
	}

    defer db.Close()

    result, err := db.Exec("update profile set picture=? where username=?", picture, username)

	if err != nil {
        fmt.Println("exec failed, ", err)
        return err
    }

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return err
	}

	fmt.Println(rows)
	return nil;
}
