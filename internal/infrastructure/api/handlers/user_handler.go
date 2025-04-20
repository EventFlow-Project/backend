package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/services"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	config        *config.Config
	userService   *services.UserService
	jwtService    *services.JWTService
	minioService  *services.MinioService
	friendService *services.FriendService
}

func NewUserHandler(
	config *config.Config,
	userService *services.UserService,
	jwtService *services.JWTService,
	minioService *services.MinioService,
	friendService *services.FriendService,
) *UserHandler {
	return &UserHandler{
		config:        config,
		userService:   userService,
		jwtService:    jwtService,
		minioService:  minioService,
		friendService: friendService,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")

	users.Get("/getInfo", h.getUserInfo)
	users.Put("/editInfo", h.editUserInfo)
	users.Post("/uploadAvatar", h.uploadImage)
	users.Get("/search", h.searchUsers)

	friends := users.Group("/friends")
	friends.Get("/", h.getFriendsList)
	friends.Get("/incoming", h.getIncomingFriendRequests)
	friends.Post("/request", h.sendFriendRequest)
	friends.Put("/respond", h.respondToFriendRequest)
	friends.Delete("/:friendId", h.removeFriend)
}

func (h *UserHandler) getUserInfo(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	safeUser, err := h.userService.GetUserInfo(token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(safeUser)
}

func (h *UserHandler) editUserInfo(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var req models.EditUserInfo
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	safeUser, err := h.userService.EditUserInfo(token, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(safeUser)
}

func (h *UserHandler) uploadImage(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var req struct {
		Base64Data string `json:"base64_data"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileURL, err := h.minioService.UploadImage(req.Base64Data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"url": fileURL,
	})
}

func (h *UserHandler) sendFriendRequest(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	var req models.SendFriendRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := h.friendService.SendFriendRequest(userID, req.ToID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response)
}

func (h *UserHandler) respondToFriendRequest(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var req models.RespondToFriendRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := h.friendService.RespondToFriendRequest(req.FriendID, req.Accept)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) getFriendsList(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	friends, err := h.friendService.GetFriendsList(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(friends)
}

func (h *UserHandler) removeFriend(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	friendID := c.Params("friendId")
	if friendID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "friend ID is required")
	}

	err = h.friendService.RemoveFriend(userID, friendID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) searchUsers(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	name := c.Query("name")
	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name parameter is required")
	}

	users, err := h.userService.SearchUsersByName(name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}

func (h *UserHandler) getIncomingFriendRequests(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	requests, err := h.friendService.GetIncomingFriendRequests(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := make([]models.IncomingFriendRequestResponse, len(requests))
	for i, request := range requests {
		response[i] = models.IncomingFriendRequestResponse{
			ID:     request.FromID,
			Name:   request.FromName,
			Avatar: request.FromAvatar,
		}
	}

	return c.JSON(response)
}
