package database

import (
	"errors"
	"go-api-server/internal/models"
	"sync"
)

// InMemoryDB represents an in-memory database for storing users.
// This is a simple implementation using Go maps for learning purposes.
// In production, you'd use a real database like PostgreSQL, MySQL, or MongoDB.
type InMemoryDB struct {
	// users stores all users with email as the key for quick lookups
	// map[email]User allows us to quickly check if an email already exists
	users map[string]*models.User

	// goals stores all goals with ID as the key
	goals map[string]*models.Goal
	
	// mu is a read-write mutex to protect concurrent access to the users map
	// This prevents race conditions when multiple goroutines access the database
	// RWMutex allows multiple readers or one writer at a time
	mu sync.RWMutex
}

// NewInMemoryDB creates and initializes a new in-memory database instance.
// This function is called when setting up the application.
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		// Initialize the users map with make()
		users: make(map[string]*models.User),
		goals: make(map[string]*models.Goal),
	}
}

// CreateUser adds a new user to the database.
// It returns an error if a user with the same email already exists.
// Parameters:
//   - user: pointer to the User struct to be stored
// Returns:
//   - error: nil if successful, error if email already exists
func (db *InMemoryDB) CreateUser(user *models.User) error {
	// Lock the database for writing (exclusive access)
	// This prevents other goroutines from reading or writing while we're creating a user
	db.mu.Lock()
	// Defer unlocking to ensure it happens even if there's an error
	// defer means "run this when the function exits"
	defer db.mu.Unlock()
	
	// Check if a user with this email already exists
	if _, exists := db.users[user.Email]; exists {
		// Return an error if the email is already registered
		return errors.New("user with this email already exists")
	}
	
	// Store the user in the map with email as the key
	db.users[user.Email] = user
	
	// Return nil to indicate success (no error)
	return nil
}

// GetUserByEmail retrieves a user from the database by their email address.
// Parameters:
//   - email: the email address to search for
// Returns:
//   - *models.User: pointer to the found user, or nil if not found
//   - error: nil if found, error if user doesn't exist
func (db *InMemoryDB) GetUserByEmail(email string) (*models.User, error) {
	// Lock the database for reading (shared access allowed)
	// Multiple goroutines can read at the same time, but not while someone is writing
	db.mu.RLock()
	// Unlock when the function exits
	defer db.mu.RUnlock()
	
	// Look up the user by email
	// The comma-ok idiom: user gets the value, exists is true/false
	user, exists := db.users[email]
	if !exists {
		// Return nil user and an error if not found
		return nil, errors.New("user not found")
	}
	
	// Return the found user and no error
	return user, nil
}

// GetUserByID retrieves a user from the database by their ID.
// Note: This is less efficient than GetUserByEmail because we have to iterate
// through all users. In a real database, you'd have an index on the ID field.
// Parameters:
//   - id: the user ID to search for
// Returns:
//   - *models.User: pointer to the found user, or nil if not found
//   - error: nil if found, error if user doesn't exist
func (db *InMemoryDB) GetUserByID(id string) (*models.User, error) {
	// Lock for reading
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	// Iterate through all users to find matching ID
	// range iterates over map: key (email), value (user pointer)
	for _, user := range db.users {
		if user.ID == id {
			// Found the user, return it
			return user, nil
		}
	}
	
	// No user found with this ID
	return nil, errors.New("user not found")
}

// DeleteUser removes a user from the database by their email.
// This is used for the logout functionality (though JWT logout is typically
// handled differently - see comments in auth handler).
// Parameters:
//   - email: the email of the user to delete
// Returns:
//   - error: nil if successful, error if user doesn't exist
func (db *InMemoryDB) DeleteUser(email string) error {
	// Lock for writing (exclusive access)
	db.mu.Lock()
	defer db.mu.Unlock()
	
	// Check if user exists before trying to delete
	if _, exists := db.users[email]; !exists {
		return errors.New("user not found")
	}
	
	// Delete the user from the map using the built-in delete function
	delete(db.users, email)
	
	return nil
}

// GetAllUsers returns a slice of all users in the database.
// This is useful for admin functionality or testing.
// Returns:
//   - []*models.User: slice containing pointers to all users
func (db *InMemoryDB) GetAllUsers() []*models.User {
	// Lock for reading
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	// Create a slice to hold all users
	// Pre-allocate with the correct capacity for efficiency
	users := make([]*models.User, 0, len(db.users))
	
	// Iterate through the map and append each user to the slice
	for _, user := range db.users {
		users = append(users, user)
	}
	
	return users
}

// CreateGoal adds a new goal to the database.
// It returns an error if a goal with the same ID already exists.
// Parameters:
//   - goal: pointer to the Goal struct to be stored
// Returns:
//   - error: nil if successful, error if ID already exists
func (db *InMemoryDB) CreateGoal(goal *models.Goal) error {
	// Lock the database for writing (exclusive access)
	// This prevents other goroutines from reading or writing while we're creating a goal
	db.mu.Lock()
	// Defer unlocking to ensure it happens even if there's an error
	// defer means "run this when the function exits"
	defer db.mu.Unlock()

	// Check if a goal with this ID already exists
	if _, exists := db.goals[goal.ID]; exists {
		// Return an error if the ID is already registered
		return errors.New("goal with this ID already exists")
	}

	// Store the goal in the map with ID as the key
	db.goals[goal.ID] = goal

	// Return nil to indicate success (no error)
	return nil
}

// GetGoalsByUserID retrieves all goals for a specific user.
// Parameters:
//   - userID: the ID of the user whose goals we want to retrieve
// Returns:
//   - []*models.Goal: slice containing pointers to the user's goals
//   - error: nil if successful, error if user has no goals or other error
func (db *InMemoryDB) GetGoalsByUserID(userID string) ([]*models.Goal, error) {
	// Lock the database for reading (shared access allowed)
	// Multiple goroutines can read at the same time, but not while someone is writing
	db.mu.RLock()
	// Unlock when the function exits
	defer db.mu.RUnlock()

	// Create a slice to hold the user's goals
	var userGoals []*models.Goal
	// Iterate through all goals to find those belonging to the user
	for _, goal := range db.goals {
		if goal.UserID == userID {
			// Add the goal to the user's goals slice
			userGoals = append(userGoals, goal)
		}
	}

	// Return the slice of user goals (empty if none found) and no error
	return userGoals, nil
}

// GetGoalByID retrieves a goal from the database by its ID.
// Parameters:
//   - id: the ID of the goal to retrieve
// Returns:
//   - *models.Goal: pointer to the found goal, or nil if not found
//   - error: nil if found, error if goal doesn't exist
func (db *InMemoryDB) GetGoalByID(id string) (*models.Goal, error) {
	// Lock the database for reading (shared access allowed)
	// Multiple goroutines can read at the same time, but not while someone is writing
	db.mu.RLock()
	// Unlock when the function exits
	defer db.mu.RUnlock()

	// Look up the goal by ID
	goal, exists := db.goals[id]
	if !exists {
		// Return nil goal and an error if not found
		return nil, errors.New("goal not found")
	}

	// Return the found goal and no error
	return goal, nil
}

// UpdateGoal updates an existing goal in the database.
// Parameters:
//   - goal: pointer to the Goal struct with updated information
// Returns:
//   - error: nil if successful, error if goal doesn't exist
func (db *InMemoryDB) UpdateGoal(goal *models.Goal) error {
	// Lock the database for writing (exclusive access)
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the goal exists before trying to update
	if _, exists := db.goals[goal.ID]; !exists {
		return errors.New("goal not found")
	}

	// Update the goal in the map
	db.goals[goal.ID] = goal

	return nil
}

// DeleteGoal removes a goal from the database by its ID.
// Parameters:
//   - id: the ID of the goal to delete
// Returns:
//   - error: nil if successful, error if goal doesn't exist
func (db *InMemoryDB) DeleteGoal(id string) error {
	// Lock the database for writing (exclusive access)
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the goal exists before trying to delete
	if _, exists := db.goals[id]; !exists {
		return errors.New("goal not found")
	}

	// Delete the goal from the map using the built-in delete function
	delete(db.goals, id)

	return nil
}
