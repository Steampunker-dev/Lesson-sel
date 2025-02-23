package handler

import (
	"awesomeProject/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegUser
// @Description create user
// @Tags user
// @Produce  json
// @Param user body models.CreateUserRequest true "User info"
// @Success 201 {object} models.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/register [post]
func (h *Handler) RegUser(ctx *gin.Context) {
	var request models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.CreateUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateUserResponse{
		User: user,
	})
}

// UpdateUser
// @Description update user
// @Tags user
// @Produce  json
// @Param user body models.UpdateUserRequest true "User info"
// @Success 200 {object} models.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/update [put]
func (h *Handler) UpdateUser(ctx *gin.Context) {
	var request models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.UpdateUser(request.ID, request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateUserResponse{
		User: user,
	})
}

// AuthUser
// @Description auth user
// @Tags user
// @Produce  json
// @Param user body models.CreateUserRequest true "User info"
// @Success 200 {object} models.AuthUserRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func (h *Handler) AuthUser(ctx *gin.Context) {
	var request models.AuthUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Repository.AuthUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.AuthUserResponse{
		Token: token,
	})
}

// LogoutUser
// @Description logout user
// @Tags user
// @Produce  json
// @Param user body models.LogoutUserRequest true "User info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
	var request models.LogoutUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repository.LogoutUser(request.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user logged out"})

	/*// получаем заголовок
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, prefix) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	// отрезаем префикс
	jwtStr = jwtStr[len(prefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		log.Println(err)
		return
	}

	// сохраняем в блеклист редиса
	//err = h.Repository.WriteJWTToBlacklist(jwtStr, time.Hour)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cant write to blacklist"})
		log.Println(err)
		return
	}

	ctx.Status(http.StatusOK)*/

}

/*
func (h *Handler) LoginUser(gCtx *gin.Context) {
	cfg := h.
	req := &models.LoginUserRequest{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Login == login && req.Password == password {
		// значит проверка пройдена
		// генерируем ему jwt
		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "DeliVeryWell-admin",
			},
			UserUUID: uuid.New(), // test uuid
			Scopes:   []string{}, // test data
		})

		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		gCtx.JSON(http.StatusOK, models.LoginUserResponse{
			ExpiresIn:   cfg.JWT.ExpiresIn,
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
	}

	gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен
}

*/
