package usecases

import (
	"context"
	"fmt"

	"github.com/Luzin7/alert-service/internal/domain"
)

type ProcessAlert struct {
	linkGen domain.LinkGenerator
	repo    domain.AlertRepository
	sender  domain.TempEmailSender
}

func (u *ProcessAlert) Execute(ctx context.Context, alert *domain.Alert) error {
	link := u.linkGen.Generate(alert.Origin, alert.Destination, alert.OutboundDate, alert.ReturnDate)

	alert.Link = link

	userEmail, err := u.repo.GetUserEmail(ctx, alert.ID)
	if err != nil {
		return err
	}

	alertEmail := &domain.AlertEmail{
		To:      userEmail,
		Subject: "Price Alert Updated",
		Body:    fmt.Sprintf("O preço do seu alerta foi atualizado. Novo preço: %.2f %s. Link: %s", alert.NewPrice, alert.Currency, alert.Link),
	}

	err = u.sender.Send(alertEmail.To, alertEmail.Subject, alertEmail.Body)
	if err != nil {
		return err
	}

	return nil
}
