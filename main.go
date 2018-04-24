package main

import "fmt"
import "unicode/utf8"
import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/spf13/viper"

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func getUrl(short string) string {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	var username string = viper.Get("username").(string)
	var password string = viper.Get("password").(string)
	var database string = viper.Get("database").(string)

	if err != nil {
		panic(fmt.Errorf("Fatal error. Missing config file"))
	}

	db, err := sql.Open("mysql", username+":" + password + "@/" + database)
	if err != nil {
		panic(err.Error())
	}

	short = trimFirstRune(short)	

	// Prior to production use I need to look into how to sanitize this input.
	stmtOut, err := db.Prepare("SELECT url FROM url_redirects WHERE short = ?")
	
	if err !=nil {
		panic(err.Error())
	}

	var results string
	

	err = stmtOut.QueryRow(short).Scan(&results)
	if err != nil {
		panic(err.Error())
	}

	s := fmt.Sprintf("The short code for this link leads to %s", results)
	fmt.Println(s)

	defer db.Close()
	return results
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/"  && req.URL.Path != "/favicon.ico" {
		var url string
		url = getUrl(req.URL.Path)
		http.Redirect(w, req, url, http.StatusSeeOther)
	}
	fmt.Fprint(w,"Krlio is currently upgrading!")
}

func main() {
	fmt.Println("Attempting to listen...")
	fmt.Println("Krlio minimal server listening.")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
