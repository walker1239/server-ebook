package main

import (
	"log"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"encoding/json"
	"io/ioutil"
	//"time"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"strconv"
)


type UserFavorites struct{
	Id int 						`json:"id"`
	User_Id int				`json:"user_id"`		
	Favorite_Id   string		`json:"favorite_id"`
	Json_book string	`json:"json_book"`
}

type Favorite struct{
	Id string 			`json:"id"`
	Json_book string	`json:"json_book"`
}

type User struct{
	Id int
	Email string		`json:"email"`		
	Password string		`json:"password"`	
	Token string			
}


var users []User
var favorites []Favorite
//var db *sql.DB

var jwtKey = []byte("13Akfq195g")

func dbConn() (db *sql.DB) {
	//db, errdb := sql.Open("mysql", "root:97477731@tcp(127.0.0.1:3306)/test")
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "974777331"
    dbName := "ebookreader"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func holaName(w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]
	fmt.Fprint(w, "hola "+name)
}

func getFavorites(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(favorites)
}

func getUsers(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(users)
}

func getFavorite(w http.ResponseWriter, r *http.Request){
	//id_String := mux.Vars(r)["id"]
	//id,_ := strconv.Atoi(id_String)
	/*for i := range favorites{
		/*if favorites[i].ID == id_String{
			json.NewEncoder(w).Encode(favorites[i])
		}*/
	//}
	
}

func getFavoritesByUser(w http.ResponseWriter, r *http.Request){
	//var favoritesUser []Favorite

	tokenHeader := r.Header.Get("Authorization")
	//id_String := mux.Vars(r)["id"]
	db := dbConn()

	selDB, err := db.Query("SELECT fa.json_book FROM favorites fa inner join favorites_users fu on fa.id = fu.favorite_id inner join users us on us.id=fu.user_id WHERE us.token=?", tokenHeader)
    if err != nil {
        panic(err.Error())
	}
	var favoritesf []string
    for selDB.Next() {
        var json_fav string
        err = selDB.Scan(&json_fav)
        if err != nil {
            panic(err.Error())
		}
		favoritesf=append(favoritesf,json_fav)
	}

	json.NewEncoder(w).Encode(favoritesf)
}

func assignFavoriteToUser(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	var temp UserFavorites
	//id_String := mux.Vars(r)["id"]
	body,_:=ioutil.ReadAll(r.Body)
	
	json.Unmarshal(body,&temp)
	query1, _ := db.Prepare("INSERT INTO  favorites(id, json_book) VALUES(?,?)")
	query1.Exec(temp.Favorite_Id,temp.Json_book)
	insForm, err := db.Prepare("INSERT INTO favorites_users(user_id, favorite_id) VALUES(?,?)")
	if err != nil {
		message1 := map[string]string{
			"code":"Error",
			"message":"No se pudo insertar",
		}
		json.NewEncoder(w).Encode(message1)
		//fmt.Println(err)
		return
	}
	fmt.Println(temp.User_Id)
	fmt.Println(temp.Favorite_Id)
	insForm.Exec(temp.User_Id,temp.Favorite_Id)
	defer db.Close()
	message1 := map[string]string{
		"code":"Correcto",
		"message":"Se inserto",
	}
	json.NewEncoder(w).Encode(message1)
}

func createFavorite(w http.ResponseWriter, r *http.Request){
    db := dbConn()
	var temp Favorite
	//id_String := mux.Vars(r)["id"]
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	insForm, err := db.Prepare("INSERT INTO favorites(id, json_book) VALUES(?,?)")
	if err != nil {
		message1 := map[string]string{
			"code":"Error",
			"message":"No se pudo insertar",
		}
		json.NewEncoder(w).Encode(message1)
		//fmt.Println(err)
		return
	}
	//fmt.Println(temp.Json_book)
	insForm.Exec(temp.Id,temp.Json_book)
	defer db.Close()
	message1 := map[string]string{
		"code":"Correcto",
		"message":"Se inserto",
	}
	json.NewEncoder(w).Encode(message1)
}

func remove(favoritesdel []Favorite, i int) []Favorite {
    favoritesdel[len(favoritesdel)-1], favoritesdel[i] = favoritesdel[i], favoritesdel[len(favoritesdel)-1]
    return favoritesdel[:len(favoritesdel)-1]
}

func deleteFavorite(w http.ResponseWriter, r *http.Request){
	//id_String := mux.Vars(r)["id"]
	//id,_ := strconv.Atoi(id_String)

	tokenHeader := r.Header.Get("Authorization")

	for i := range users{
		if users[i].Token == tokenHeader{
			/*for j := range favorites{
				if favorites[j].Id == id_String && favorites[j].EmailUser == users[i].Email {
					favorites=remove(favorites,j);
					message1 := map[string]string{
						"code":"correcto",
						"message":"Favorito eliminada correctamente",
					}
					json.NewEncoder(w).Encode(message1)
				}
			}*/
			message1 := map[string]string{
				"code":"error",
				"message":"Favorito no encontrada",
			}
			json.NewEncoder(w).Encode(message1)
			return
		}
	}
	message1 := map[string]string{
		"code":"error",
		"message":"Token no reconocido, vuelva a iniciar sesión",
	}
	json.NewEncoder(w).Encode(message1)
}

func createUser(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	var temp User
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	insForm, err := db.Prepare("INSERT INTO users(email, password, token) VALUES(?,?,?)")
	if err != nil {
		message1 := map[string]string{
			"code":"Error",
			"message":"No se pudo insertar",
		}
		json.NewEncoder(w).Encode(message1)
		//fmt.Println(err)
		return
	}
	//fmt.Println(temp.Json_book)
	//account.Token="none"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(temp.Password), bcrypt.DefaultCost)
	temp.Password = string(hashedPassword)
	//users = append(users,account)
	insForm.Exec(temp.Email,temp.Password,"null")
	defer db.Close()
	message1 := map[string]string{
		"code":"Correcto",
		"message":"Se inserto",
	}
	json.NewEncoder(w).Encode(message1)
}

func loginUser(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	var account User
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&account)
	//favorites = append(favorites,temp)
	selDB, err := db.Query("SELECT password, id FROM users WHERE email=?", account.Email)
    if err != nil {
        panic(err.Error())
    }
    for selDB.Next() {
        var password string
        var idUser int
        err = selDB.Scan(&password, &idUser)
        if err != nil {
            panic(err.Error())
		}
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(account.Password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { 
			message := map[string]string{
				"code":"error",
				"message":"Contraseña o correo invalidos",
			}
			json.NewEncoder(w).Encode(message)
			return
		}
		userClaims := jwt.MapClaims{}
		userClaims["email"] = account.Email
		userClaims["password"] = password
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
		token, _ := at.SignedString(jwtKey)
		query, errdb := db.Prepare("UPDATE users SET token=? WHERE email=?")
        if errdb != nil {
            panic(errdb.Error())
        }
		query.Exec(token, account.Email)
		defer db.Close()
		message := map[string]string{
			"code":"correcto",
			"message":"Login correcto",
			"Api": token,
			"Id":strconv.Itoa(idUser),
		}
		json.NewEncoder(w).Encode(message)
		return
    }
}

func main(){
	

	//favorites=append(favorites,Favorite{Id:1,Title:"titulo1",Description:"description1"})
	//favorites=append(favorites,Favorite{Id:2,Title:"titulo2",Description:"description2"})

	r := mux.NewRouter()
	r.HandleFunc("/favorite", createFavorite).Methods("POST")
	r.HandleFunc("/favorite-user", assignFavoriteToUser).Methods("POST")
	r.HandleFunc("/favorites", getFavorites).Methods("GET")
	r.HandleFunc("/user/favorites", getFavoritesByUser).Methods("GET")
	r.HandleFunc("/favorites/{id}", getFavorite).Methods("GET")
	//r.HandleFunc("/favorites/{id}", updateFavorite).Methods("PUT")
	r.HandleFunc("/favorites/{id}", deleteFavorite).Methods("DELETE")
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/login", loginUser).Methods("POST")

    //r.HandleFunc("/{name}", holaName).Methods("GET")

	log.Print("Corriendo en el puerto 8085")
	err := http.ListenAndServe(":8085",r)
	if err != nil{
		log.Fatal("error: ",err)
	}
}