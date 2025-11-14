package handler

import (
	"net/http"

	svc "moment/internal/service"

	"github.com/gin-gonic/gin"
)

type EncodeRequest struct {
	Text      string `json:"text" binding:"required"`
	ScreatKey string `json:"screat_key" binding:"required"`
}

type EncodeResponse struct {
	EncryptedText string `json:"encrypted_text"`
}

type DecodeRequest struct {
	EncryptedText string `json:"encrypted_text" binding:"required"`
	ScreatKey     string `json:"screat_key" binding:"required"`
}

type DecodeResponse struct {
	DecryptedText string `json:"decrypted_text"`
}

type TextCryptoHandler struct{}

func NewTextCryptoHandler() *TextCryptoHandler {
	return &TextCryptoHandler{}
}

// GET /textcrypto/random-screat-key
func (h *TextCryptoHandler) RandomScreatKey(c *gin.Context) {
	key := svc.RandomScreatKey()
	c.JSON(http.StatusOK, gin.H{
		"screat_key": key,
	})
}

// POST /textcrypto/encrypt
func (h *TextCryptoHandler) Encrypt(c *gin.Context) {
	var req EncodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, EncodeResponse{
		EncryptedText: svc.Encrypt(req.Text, req.ScreatKey),
	})
}

// POST /textcrypto/decrypt
func (h *TextCryptoHandler) Decrypt(c *gin.Context) {
	var req DecodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, DecodeResponse{
		DecryptedText: svc.Decrypt(req.EncryptedText, req.ScreatKey),
	})
}
