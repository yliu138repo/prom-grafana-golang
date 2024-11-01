package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"goprometheus/internal/pgdatabase"
	"goprometheus/internal/validator"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type UserHandler struct {
	pgdb *pgdatabase.Postgresdb
}

func NewUserHandler(pg *pgdatabase.Postgresdb) *UserHandler {
	if pg == nil {
		log.Fatal("database connection cannot be nil")
	}
	return &UserHandler{pgdb: pg}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid CreateUserRequest payload", http.StatusBadRequest)
		return
	}

	// Validate the request DTO
	validationErrors := validator.ValidateStruct(req)
	if len(validationErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": validationErrors,
		})
		return
	}

	user, err := h.createUserInDB(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05-0700"),
	})
}

// GetUser handles fetching a user by ID.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.getUserFromDB(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05-0700"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05-0700"),
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.getAllUsersFromDB(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) getAllUsersFromDB(ctx context.Context) ([]UserModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//log.Printf("Handler PG Pool Pointer: %p\n", h.pgdb.Pool)

	stats := h.pgdb.GetStats()
	_ = stats

	if err := h.pgdb.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database connection check failed: %w", err)
	}

	query := "SELECT id, name, email, age,created_at, updated_at FROM public.pgusers"
	rows, err := h.pgdb.Pool.Query(ctx, query)
	if err != nil {
		return []UserModel{}, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	users := []UserModel{}

	for rows.Next() {
		user := UserModel{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return []UserModel{}, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (h *UserHandler) getUserFromDB(ctx context.Context, id string) (UserModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "SELECT id, name, email, age, created_at, updated_at FROM public.pgusers WHERE id = $1"
	row := h.pgdb.Pool.QueryRow(ctx, query, id)

	var user UserModel
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

func (h *UserHandler) createUserInDB(ctx context.Context, req CreateUserRequest) (UserModel, error) {
	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.pgusers (id, name, email, age, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, name, email, age, created_at, updated_at
	`

	var user UserModel
	err := h.pgdb.Pool.QueryRow(ctx, query,
		uuid.New().String(),
		req.Name,
		req.Email,
		req.Age,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == context.DeadlineExceeded {
			return UserModel{}, fmt.Errorf("database operation timed out: %w", err)
		}
		return UserModel{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,gte=18,lte=100"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
}

type GetUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
type UserModel struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt *time.Time // nullable
}
