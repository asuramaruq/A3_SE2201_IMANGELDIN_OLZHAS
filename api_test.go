package greenlight_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readResponseBody(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
}

type User struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Activated bool   `json:"activated"`
}

var testToken string

func TestAccountCreationRequest(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Name":     "test3211",
		"Email":    "test3211@gmail.com",
		"Password": "test321321",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/users", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response struct {
		User User `json:"user"`
	}

	str := readResponseBody(resp)

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		assert.Fail(t, "Error due Response Unmarshall")
		return
	}
	assert.Positive(t, response.User.ID)
	assert.Equal(t, "test3211", response.User.Name)
	assert.Equal(t, "test3211@gmail.com", response.User.Email)
}

type ErrorIncorrectAccountCreationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestIncorrectAccountCreationRequest(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Name":     "asdfj",
		"Email":    "asdfjk",
		"Password": ""})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/users", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response struct {
		Error ErrorIncorrectAccountCreationRequest `json:"error"`
	}

	str := readResponseBody(resp)

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}
	assert.Equal(t, "must be provided", response.Error.Password)
	assert.Equal(t, "must be valid email address", response.Error.Email)
}

type TokenForTestLogin struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}

func TestLogin(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Email":    "test3211@gmail.com",
		"Password": "test321321",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/tokens/authentication", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response struct {
		Token TokenForTestLogin `json:"authentication_token"`
	}

	str := readResponseBody(resp)

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}
	testToken = response.Token.Token
	assert.NotNil(t, response.Token.Token)
	assert.NotNil(t, response.Token.Expiry)
}

func TestIncorrectLogin(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Email":    "aksldf@gmail.com", // This email does not exist
		"Password": "password123",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:4000/v1/tokens/authentication", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response struct {
		Error string `json:"error"`
	}

	str := readResponseBody(resp)

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}
	assert.Equal(t, "invalid authentication credentials", response.Error)
}

func TestValidToken(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Token": "NY5A4Z7K256WVXYR57HLFE3ZEM", // ACTIVATION TOKEN
	})
	responseBody := bytes.NewBuffer(postBody)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:4000/v1/users/activated", responseBody)
	if err != nil {
		log.Fatalf("An Error Occurred while creating request: %v", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occurred while sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	str := readResponseBody(resp)

	var response struct {
		Error struct {
			Token string `json:"token"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}
	assert.Equal(t, "invalid or expired activation token", response.Error.Token)
}

func TestInvalidToken(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"Token": "aaaaaaaaaaaaaaaaaaaaaaaaaa", // INVALID ACTIVATION TOKEN
	})
	responseBody := bytes.NewBuffer(postBody)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:4000/v1/users/activated", responseBody)
	if err != nil {
		log.Fatalf("An Error Occurred while creating request: %v", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occurred while sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	str := readResponseBody(resp)

	var response struct {
		Error struct {
			Token string `json:"token"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(str), &response); err != nil {
		fmt.Println("Error:", err)
		return
	}
	assert.Equal(t, "invalid or expired activation token", response.Error.Token)
}

type Movie struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Runtime string   `json:"runtime"`
	Genres  []string `json:"genres"`
}

func TestInsertingMoviesIntoDatabase(t *testing.T) {
	moviePayload := Movie{
		Title:   "Inception",
		Year:    2010,
		Runtime: "144 mins",
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

	bearerToken := "NY5A4Z7K256WVXYR57HLFE3ZEM"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	readResponseBody(resp)
}

func TestInsertingMoviesIntoDatabaseWithWrongYear(t *testing.T) {
	moviePayload := Movie{
		Title:   "Inception",
		Year:    2029, // This is a future date
		Runtime: "144 mins",
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

	bearerToken := "NY5A4Z7K256WVXYR57HLFE3ZEM"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	readResponseBody(resp)
}

func TestInsertingMoviesIntoDatabaseWithWrongRuntime(t *testing.T) {
	moviePayload := Movie{
		Title:   "Inception",
		Year:    2020,
		Runtime: "144", // This is not a valid runtime
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

	bearerToken := "NY5A4Z7K256WVXYR57HLFE3ZEM"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	readResponseBody(resp)
}

func TestMovieDeletionById(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:4000/v1/movies/5", nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	bearerToken := "NY5A4Z7K256WVXYR57HLFE3ZEM"
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	readResponseBody(resp)
}
