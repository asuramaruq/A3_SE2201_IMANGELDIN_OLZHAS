package greenlight_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
)

func printResponseBody(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Print(sb)
}

func printRequestBody(req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Print(sb)

}

// Unit tests for the API
func TestCreationNewAccount(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Name":     "Olzhas",
		"Email":    "example@test.com",
		"Password": "password123",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/users", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestCreationNewAccount2(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Name":     "Olzhas",
		"Email":    "asdfasdf",
		"Password": ""})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/users", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestGettingAuthenticationToken(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Email":    "example@test.com",
		"Password": "password123",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/tokens/authentication", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestGettingAuthenticationToken2(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Email":    "fakeexample@test.com", // This email does not exist
		"Password": "password123",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/tokens/authentication", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

//ATTENTION FOR THE LAST 2 TESTS RUN THE TESTS IN ORDER

func TestActivationAccount(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Token": "MQJNOQ3BTN6MHBENQBRROMA2FY", // ACTIVATION TOKEN
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.NewRequest(http.MethodPost, "http://localhost:4000/v1/users/activated", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printRequestBody(resp)
}

func TestActivationAccount2(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Token": "NOT A VALID TOKEN", // ACTIVATION TOKEN
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.NewRequest(http.MethodPost, "http://localhost:4000/v1/users/activated", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	printRequestBody(resp)
}

//Integration tests for the API

type Movie struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Runtime string   `json:"runtime"`
	Genres  []string `json:"genres"`
}

func TestInsertMovie(t *testing.T) {
	moviePayload := Movie{
		Title:   "Fight Club",
		Year:    1999,
		Runtime: "139 mins",
		Genres:  []string{"drama"},
	}

	payloadBytes, err := json.Marshal(moviePayload)
	if err != nil {
		log.Fatalf("Error marshaling movie payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:4000/v1/movies", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	bearerToken := "YXJXRFN44TZTZJ4OES3BVCR2RQ"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestInsertMovieWithWrongYear(t *testing.T) {
	moviePayload := Movie{
		Title:   "Fight Club",
		Year:    2000, // This is a future date
		Runtime: "139 mins",
		Genres:  []string{"drama"},
	}

	payloadBytes, err := json.Marshal(moviePayload)
	if err != nil {
		log.Fatalf("Error marshaling movie payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:4000/v1/movies", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	bearerToken := "YXJXRFN44TZTZJ4OES3BVCR2RQ"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestInsertMovieWithWrongRuntime(t *testing.T) {
	moviePayload := Movie{
		Title:   "Fight Club",
		Year:    2020,
		Runtime: "139", // This is not a valid runtime
		Genres:  []string{"thriller"},
	}

	payloadBytes, err := json.Marshal(moviePayload)
	if err != nil {
		log.Fatalf("Error marshaling movie payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:4000/v1/movies", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	bearerToken := "YXJXRFN44TZTZJ4OES3BVCR2RQ"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}

func TestMovieDeletionById(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:4000/v1/movies/3", nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	bearerToken := "YXJXRFN44TZTZJ4OES3BVCR2RQ"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	printResponseBody(resp)
}
