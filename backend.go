package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type FormData struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Savelogin bool   `json:"savelogin"`
}

type LoginData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Savelogin bool   `json:"savelogin"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

var tokenStore = make(map[int]string)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("123")
	mux := http.NewServeMux()

	db, err := sql.Open("mysql", "")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("LOGIN-REGISTER API STARTED")

	mux.HandleFunc("/register/newuser", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		var formData FormData
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		username, email, password := formData.Username, formData.Email, formData.Password
		if len(username) == 0 || len(username) > 128 {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}
		if len(email) < 4 || len(email) > 128 || !isValidEmail(email) {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}
		if len(password) < 8 || len(password) > 128 {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO Users (Name, Email, HashedPassword) VALUES (?, ?, ?)", username, email, hashedPassword)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Registration successful"}`))
	})

	//LOGIN FUNCTION
	mux.HandleFunc("/login/newrequest", func(writer http.ResponseWriter, request *http.Request) {
		enableCORS(writer)
		if request.Method != http.MethodPost {
			http.Error(writer, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		if request.Header.Get("Content-Type") != "application/json" {
			http.Error(writer, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		//decode the login data
		var loginData LoginData
		err := json.NewDecoder(request.Body).Decode(&loginData)
		if err != nil {
			http.Error(writer, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		email := loginData.Email
		password := loginData.Password
		savelogin := loginData.Savelogin

		//perform same checks as client
		if len(email) < 4 || len(email) > 128 || !isValidEmail(email) {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": "Invalid credentials"}`))
			return
		}
		if len(password) < 8 || len(password) > 128 {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": "Invalid credentials"}`))
			return
		}

		// get the hashed password
		query := "SELECT HashedPassword FROM USERS WHERE Email = ?"
		var hashedPassword string
		err = db.QueryRow(query, email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				writer.WriteHeader(http.StatusUnauthorized)
				writer.Write([]byte(`{"message": "Invalid credentials"}`))
				return
			}
			log.Printf("SQL ERROR: %v", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		//compare passwords
		passMatchErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if passMatchErr != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": "Invalid credentials"}`))
			return
		}
		// get the userid
		var uid int
		query = "SELECT ID FROM USERS WHERE Email = ?"
		db.QueryRow(query, email).Scan(&uid)

		//delete the original token
		deleteToken(uid)

		//generate token
		token, err := generateJWT(uid, savelogin)
		if err != nil {
			fmt.Println("ERROR! ", err)
		}

		//save token
		saveToken(uid, token)
		fmt.Println("New token generated!", token)

		response := fmt.Sprintf(`{"message": "Login successful", "token": "%s"}`, token)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(response))
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: corsMiddleware(mux),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func isValidEmail(email string) bool {
	// Regular expression for validating an email
	re := regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	return re.MatchString(strings.ToLower(email))
}

func generateJWT(uid int, saveLogin bool) (string, error) {
	// Define the secret key
	secretKey := "secretest key evar!!!!!"

	//when savelogin is true, generate a token which lasts 1 week, instead of 1h
	var tokenTime int
	if saveLogin {
		tokenTime = 168
	} else {
		tokenTime = 24
	}

	// set the claims for tokens
	claims := jwt.MapClaims{
		"userID": uid,
		"exp":    time.Now().Add(time.Hour * time.Duration(tokenTime)).Unix(),
		"iat":    time.Now().Unix(), //time the token was issued at
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func saveToken(uid int, token string) {
	tokenStore[uid] = token
}

func deleteToken(uid int) {
	delete(tokenStore, uid)
}

/* currently unused
func getToken(uid int) (string, bool) {
	token, exists := tokenStore[uid]
	return token, exists
}
*/
