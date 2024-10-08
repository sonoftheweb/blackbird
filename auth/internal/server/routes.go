package server

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.Handle("/admin", s.jwtMiddleware(http.HandlerFunc(s.HelloWorldHandler)))

	mux.HandleFunc("/register", s.registerHandler)
	mux.HandleFunc("/login", s.loginHandler)

	mux.HandleFunc("/.well-known/jwks.json", s.jwkHandler)

	mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	_, err = s.db.Exec(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]bool{"status": true})
	if err != nil {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedPassword string
	ctx := r.Context()
	var userID int
	err := s.db.QueryRow(ctx, "SELECT id, password FROM users WHERE email = $1", creds.Email).Scan(&userID, &storedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(userID, creds.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Error marshaling token", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(tokenJSON)
	if err != nil {
		http.Error(w, "Error writing token", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) jwkHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/server/jwk_set.json")
}
