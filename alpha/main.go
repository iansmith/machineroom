package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coocood/qbs"
	_ "github.com/lib/pq"

	"github.com/igneous-systems/lib/consul" //shared with beta, server only
)

const (
	//unclear if we should really link against beta
	USERPROP = "postgres/host_count/username"
	PWDPROP  = "postgres/host_count/password"
)

type HostCount struct {
	Hostname string `qbs:"pk"`
	Count    int
}

func tryPostgres(user string, pwd string) (*sql.DB, error) {
	app := "alpha"
	url := fmt.Sprintf("postgres://%s:%s@%s.%s/%s?sslmode=disable", user, pwd, app, "postgres.service.consul", app)
	log.Printf("trying postgres url: %s", url)
	return sql.Open("postgres", url)
}

func readConfig() error {
	//try to contact the DB
	resp, err := consul.ReadKV(USERPROP)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("no configuration found for database!")
	}
	user := resp.DecodedValue
	resp, err = consul.ReadKV(PWDPROP)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("can't find a password, but got a username")
	}
	pwd := resp.DecodedValue

	log.Printf("read the username and password for db from consul: %s, %s", user, pwd)

	db, err := tryPostgres(user, pwd)
	if err != nil {
		return fmt.Errorf("failed to connect to the database (networking failed): %v", err)
	}
	qbs.RegisterWithDb("postgres", db, qbs.NewPostgres())

	//at this point we will only have failed if the connectivity is bad
	//not anything with auth because postgres doesn't try that until
	//you do sql
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := readConfig(); err != nil {
		fmt.Fprintf(w, "read config: %v", err)
		return
	}
	q, err := qbs.GetQbs()
	if err != nil {
		fmt.Fprintf(w, "GetQbs: %v", err)
	}
	defer q.Close()

	h, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "hostname() failed: %v", err)

	}

	var count HostCount
	count.Hostname = h
	q.Log = true
	if err := q.Find(&count); err != nil {
		if err != sql.ErrNoRows {
			fmt.Fprintf(w, "find failed (probably your crendentials are bad): %v", err)
			return
		}
	}
	//why doesn't this happen?
	if err == sql.ErrNoRows {
		fmt.Fprintf(w, "no count found for %s", h)
	}
	count.Count++
	fmt.Fprintf(w, "new count for %s is %d", count.Hostname, count.Count)

	if _, err := q.Save(&count); err != nil {
		fmt.Fprintf(w, "save failed: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
