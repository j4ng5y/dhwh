package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/j4ng5y/dhwh/pkg/dkr"
)

type PushData struct {
	Images   []string  `json:"images"`
	PushedAt time.Time `json:"pushed_at"`
	Pusher   string    `json:"pusher"`
	Tag      string    `json:"tag"`
}

type Repository struct {
	CommentCount    int       `json:"comment_count"`
	DateCreated     time.Time `json:"date_created"`
	Description     string    `json:"description"`
	Dockerfile      string    `json:"dockerfile"`
	FullDescription string    `json:"full_description"`
	IsOfficial      bool      `json:"is_official"`
	IsPrivate       bool      `json:"is_private"`
	IsTrusted       bool      `json:"is_trusted"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	Owner           string    `json:"owner"`
	RepoName        string    `json:"repo_name"`
	RepoURL         url.URL   `json:"repo_url"`
	StarCount       int       `json:"star_count"`
	Status          string    `json:"status"`
}

type Request struct {
	CallbackURL url.URL    `json:"callback_url"`
	PushData    PushData   `json:"push_data"`
	Repository  Repository `json:"respoitory"`
}

func (R *Request) Unmarshal(requestBody io.ReadCloser) error {
	b, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, R)
}

func (R Request) Validate() (ok bool, err error) {
	var (
		expectedNamespace = "my_expected_namespace"
		actualNamespace   = R.Repository.Namespace
		expectedReponame  = "my_expected_reponame"
		actualReponame    = R.Repository.Name
		expectedOwner     = "my_expected_owner"
		actualOwner       = R.Repository.Owner
	)

	switch {
	case actualNamespace != expectedNamespace:
		return false, fmt.Errorf("webhook repository namespace validation failed, expected %s, got %s", expectedNamespace, actualNamespace)
	case actualReponame != expectedReponame:
		return false, fmt.Errorf("webhook repository name validation failed, expected %s, got %s", expectedReponame, actualReponame)
	case actualOwner != expectedOwner:
		return false, fmt.Errorf("webhook repository owner validation failed, expected %s, got %s", expectedOwner, actualOwner)
	// Other validation logic
	default:
		return true, nil
	}
}

type Response struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	TargetURL   string `json:"target_url"`
}

func (R *Response) Marshal() ([]byte, error) {
	return json.Marshal(R)
}

func (R *Response) SendCallback(url string) error {
	b, err := R.Marshal()
	if err != nil {
		return err
	}

	if _, err = http.Post(url, "appliction/json", bytes.NewBuffer(b)); err != nil {
		return err
	}

	return nil
}

func (S *Server) NewWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var resp = &Response{
		State:       "success",
		Description: "docker container successfully pulled and restarted",
		Context:     "my webhook handler",
		TargetURL:   "",
	}

	req := new(Request)

	if err := req.Unmarshal(r.Body); err != nil {
		log.Println(err)
		http.Error(w, "unable to deserialize the json request", http.StatusInternalServerError)
		return
	}

	if _, err := req.Validate(); err != nil {
		log.Println(err)
		http.Error(w, "webhook validation failed", http.StatusUnauthorized)
		resp = &Response{
			State:       "error",
			Description: "webhook validation failed",
			Context:     "my webhook handler",
			TargetURL:   "",
		}
		resp.SendCallback(req.CallbackURL.String())
		return
	}

	if err := dkr.PullNewContainerImage(context.Background(), "my_container"); err != nil {
		log.Println(err)
		http.Error(w, "unable to pull new container iamge", http.StatusInternalServerError)
		resp = &Response{
			State:       "error",
			Description: "docker image pull failed",
			Context:     "my webhook handler",
			TargetURL:   "",
		}
		resp.SendCallback(req.CallbackURL.String())
		return
	}

	if err := dkr.RestartContainer(context.Background(), "my_container"); err != nil {
		log.Println(err)
		http.Error(w, "unable to restart container", http.StatusInternalServerError)
		resp = &Response{
			State:       "error",
			Description: "docker container restart failed",
			Context:     "my webhook handler",
			TargetURL:   "",
		}
		resp.SendCallback(req.CallbackURL.String())
		return
	}

	if err := resp.SendCallback(req.CallbackURL.String()); err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (S *Server) WebhookStatusHandler(w http.ResponseWriter, r *http.Request) {

}
