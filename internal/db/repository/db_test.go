package repository

import (
    "testing"
    "p2pbot/internal/db/models"
    "p2pbot/internal/app"
    "github.com/jmoiron/sqlx"
    "fmt"
    "os"
    "path/filepath"
)

var (
    DB *sqlx.DB
    userRepo *UserRepository
    trackerRepo *TrackerRepository
)

func findProjectRoot() (string, error) {
    // Get the current working directory of the test
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }

    // Traverse upwards to find the project root (where .env file is located)
    for {
        if _, err := os.Stat(filepath.Join(dir, ".env")); os.IsNotExist(err) {
            parent := filepath.Dir(dir)
            if parent == dir {
                // Reached the root of the filesystem, and .env wasn't found
                return "", os.ErrNotExist
            }
            dir = parent
        } else {
            return dir, nil
        }
    }
}

func TestMain(m *testing.M) {
    root, err := findProjectRoot()
    fmt.Println("Root: ", root)
    if err != nil {
        panic("Error finding project root: " + err.Error())
    }

    err = os.Chdir(root)
    if err != nil {
        panic(err)
    }

    DB, _, err := app.Init()
    if err != nil {
        panic(err)
    }
    userRepo = NewUserRepository(DB)
    trackerRepo = NewTrackerRepository(DB)

    code := m.Run()

    //DB.MustExec("DELETE FROM users");
    
    os.Exit(code)
}

// User tests
func TestSaveUserEmail(t *testing.T) {
    email := "ab@gmail.com"
    password := "123456"

    // Create a new user
    user := &models.User{
        Email:       &email,
        Password_en: &password,
    }
    id, err := userRepo.Save(user)
    if err != nil {
        t.Fatalf("error saving user: %v", err)
    }
    fmt.Println("User ID: ", id)
}

func TestSaveUserChatID(t *testing.T) {
    chatID := int64(12345)

    user := &models.User{
        ChatID:      &chatID,
    }
    id, err := userRepo.Save(user)
    if err != nil {
        t.Fatalf("error saving user: %v", err)
    }
    fmt.Println("User ID: ", id)
}

func TestSaveUserEmail2(t *testing.T) {
    email := "a@mail.ru"
    password := "123456"

    // Create a new user
    user := &models.User{
        Email:       &email,
        Password_en: &password,
    }
    id, err := userRepo.Save(user)
    if err != nil {
        t.Fatalf("error saving user: %v", err)
    }
    fmt.Println("User ID: ", id)
}

func TestSaveUserChatID2(t *testing.T) {
    chatID := int64(123)

    user := &models.User{
        ChatID:      &chatID,
    }
    id, err := userRepo.Save(user)
    if err != nil {
        t.Fatalf("error saving user: %v", err)
    }
    fmt.Println("User ID: ", id)
}

func TestSaveUserNull(t *testing.T) {
    user := &models.User{}
    _, err := userRepo.Save(user)
    if err == nil {
        t.Fatalf("null email and chatid user can't be saved")
    }
}

func TestSaveUserEmailNoPassword(t *testing.T) {
    email := "nopass@mail.ri"

    user := &models.User{
        Email: &email,
    }
    _, err := userRepo.Save(user)
    if err == nil {
        t.Fatalf("user without password can't be saved")
    }
}

func TestGetUserByChatID(t *testing.T) {
    chatID := int64(123)
    user, err := userRepo.GetByChatID(chatID)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }
    fmt.Println("User: ", user)
}

func TestGetUserByID(t *testing.T) {
    ID := 1
    user, err := userRepo.GetByID(ID)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }
    fmt.Println("User: ", user)
}

func TestGetUserByEmail(t *testing.T) {
    email := "ab@gmail.com"
    user, err := userRepo.GetByEmail(email)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }
    fmt.Println("User: ", user)
}

func TestUpdateUser(t *testing.T) {
    newEmail := "abob@fmail.cz"
    password := "123"

    userToUpdate, err := userRepo.GetByChatID(123)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }
    userToUpdate.Email = &newEmail
    userToUpdate.Password_en = &password
    id, err := userRepo.Save(userToUpdate)
    if err != nil {
        t.Fatalf("error updating user: %v", err)
    }
    fmt.Println("updated user with ID: ", id)
}

// Tracker tests

func TestCreateTracker(t *testing.T) {
    err := trackerRepo.Save(&models.Tracker{
        UserID:   1,
        Exchange: "binance",
        Username: "anton",
        Currency: "BTC",
        Side:     "buy",
        Waiting:  false,
    })

    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    fmt.Println("Tracker created")
}

func TestUpdateTracker(t *testing.T) {
    err := trackerRepo.Save(&models.Tracker{
        ID:       1,
        UserID:   1,
        Exchange: "binance",
        Username: "anton",
        Currency: "BTC",
        Side:     "buy",
        Payment:  []string{"zen", "wise"},
        Waiting:  true,
    })

    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    t.Logf("Tracker updated")
}

func TestSetOutbided(t *testing.T) {
    err := trackerRepo.UpdateOutbided(1, true);
    if err != nil {
        t.Fatalf("error updating outbided flag")
    }
    t.Logf("Outbided flag set")
}

func TestGetAllTrackers(t *testing.T) {
    trackers, err := trackerRepo.GetAllTrackers()
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    for _, tracker := range trackers {
        fmt.Println("Tracker: ", tracker.ID,
                                    tracker.Exchange,
                                    tracker.Currency,
                                    tracker.Side,
                                    tracker.Username,
                                    tracker.Outbided,
                                    tracker.Waiting)
        for _, method := range tracker.Payment {
            fmt.Println("Payment method: ", method)
        }
    }   
}

