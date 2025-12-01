package consumer

import (
	"context"
	"encoding/json"

	"github.com/Luzin7/alert-service/internal/usecases"
)

type Handler struct {
	useCase *usecases.ProcessAlert
}

func NewHandler(uc *usecases.ProcessAlert) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) Handle(msgBody []byte) error {
	var payload PriceUpdatedPayload

	if err := json.Unmarshal(msgBody, &payload); err != nil {
		return err
	}

	alert, err := payload.ToDomain()
	if err != nil {
		return err
	}

	return h.useCase.Execute(context.Background(), alert)
}
