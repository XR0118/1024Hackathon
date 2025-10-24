package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/XR0118/1024Hackathon/deployment-trigger/internal/service"
)

type WebhookHandler struct {
	service *service.VersionService
	secret  string
}

func NewWebhookHandler(service *service.VersionService, secret string) *WebhookHandler {
	return &WebhookHandler{
		service: service,
		secret:  secret,
	}
}

func (h *WebhookHandler) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	signature := r.Header.Get("X-Hub-Signature-256")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if h.secret != "" && !h.verifySignature(body, signature) {
		log.Printf("Invalid signature for delivery %s", deliveryID)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	if eventType != "push" {
		log.Printf("Ignoring event type: %s", eventType)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Event type not supported"}`))
		return
	}

	var pushEvent PushEvent
	if err := json.Unmarshal(body, &pushEvent); err != nil {
		log.Printf("Failed to parse push event: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(pushEvent.Ref, "refs/tags/") {
		log.Printf("Ignoring non-tag push: %s", pushEvent.Ref)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Not a tag creation event"}`))
		return
	}

	tagName := strings.TrimPrefix(pushEvent.Ref, "refs/tags/")
	log.Printf("Processing tag: %s", tagName)

	result, err := h.service.ProcessTagEvent(r.Context(), &service.TagEvent{
		TagName:    tagName,
		Repository: pushEvent.Repository.CloneURL,
		Commit:     pushEvent.After,
		Pusher:     pushEvent.Pusher.Name,
	})

	if err != nil {
		log.Printf("Failed to process tag event: %v", err)
		http.Error(w, "Failed to process tag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *WebhookHandler) verifySignature(payload []byte, signature string) bool {
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(payload)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

type PushEvent struct {
	Ref    string `json:"ref"`
	After  string `json:"after"`
	Before string `json:"before"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Repository struct {
		FullName string `json:"full_name"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
}
