package v1

import (
	"net/http"

	"github.com/renderview-inc/backend/internal/app/application/services"
)

type UserAccountHandler struct {
	accountService *services.UserAccountService
}

func (uah *UserAccountHandler) HandleRegister(w http.ResponseWriter, r *http.Response) {
	
}
