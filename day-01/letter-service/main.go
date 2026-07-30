package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidLetterID    = errors.New("invalid letter ID")
	ErrAccessDenied       = errors.New("access denied")
	ErrToFindLetter       = errors.New("cannot find letter")
	ErrCannotSendToSelf   = errors.New("cannot send to self")
	ErrInvalidSenderID    = errors.New("invalid sender ID")
	ErrInvalidRecipientID = errors.New("invalid recipient ID")
	ErrEmptyMessage       = errors.New("message connot be empty")
)

type LetterRepository interface {
	GetLetterByID(letterID int64) (Letter, error)
	CreateLetter(letter Letter) (Letter, error)
}

type MemoryLetterRepository struct {
	letters map[int64]Letter
}

type LetterService struct {
	letterRepo LetterRepository
}

func NewLetterService(repo LetterRepository) *LetterService {
	return &LetterService{
		letterRepo: repo,
	}
}

type Letter struct {
	ID        int64
	Sender    int64
	Recipient int64
	Message   string
	CreatedAt time.Time
}

func (r *MemoryLetterRepository) GetLetterByID(letterID int64) (Letter, error) {
	letter, ok := r.letters[letterID]
	if !ok {
		return Letter{}, ErrToFindLetter
	}
	return letter, nil
}

func (s *LetterService) GetLetterByID(letterID int64, userID int64) (Letter, error) {
	if letterID <= 0 {
		return Letter{}, ErrInvalidLetterID
	}

	letter, err := s.letterRepo.GetLetterByID(letterID)
	if err != nil {
		return Letter{}, err
	}

	if userID != letter.Sender && userID != letter.Recipient {
		return Letter{}, ErrAccessDenied
	}

	return letter, nil
}

func (r *MemoryLetterRepository) CreateLetter(letter Letter) (Letter, error) {
	newID := int64(len(r.letters) + 1)
	letter.ID = newID
	letter.CreatedAt = time.Now()

	r.letters[newID] = letter

	return letter, nil
}

func (s *LetterService) CreateLetter(senderID int64, recipientID int64, message string) (Letter, error) {
	if senderID <= 0 {
		return Letter{}, ErrInvalidSenderID
	}
	if recipientID <= 0 {
		return Letter{}, ErrInvalidRecipientID
	}
	if senderID == recipientID {
		return Letter{}, ErrCannotSendToSelf
	}
	if message == "" {
		return Letter{}, ErrEmptyMessage
	}

	letter := Letter{
		Sender:    senderID,
		Recipient: recipientID,
		Message:   message,
	}

	createdLetter, err := s.letterRepo.CreateLetter(letter)
	if err != nil {
		return Letter{}, err
	}

	return createdLetter, nil

}
func main() {
	repo := &MemoryLetterRepository{
		letters: map[int64]Letter{
			1: {
				ID:        1,
				Sender:    10,
				Recipient: 20,
				Message:   "Привет",
				CreatedAt: time.Now(),
			},
			2: {
				ID:        2,
				Sender:    20,
				Recipient: 10,
				Message:   "и вам привет",
				CreatedAt: time.Now(),
			},
		},
	}

	service := NewLetterService(repo)

	newletter, err := service.CreateLetter(10, 20, "как дела?")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println(newletter)

	letter, err := service.GetLetterByID(newletter.ID, 20)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(letter)
}
