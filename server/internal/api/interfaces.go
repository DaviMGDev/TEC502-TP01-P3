package api

import (
// "cod-server/internal/
)

type EventHandlerInterface interface {
	OnRegister(event Event) Event
	OnLogin(event Event) Event

	OnGetCards(event Event) Event
	OnBuyPack(event Event) Event
	OnOfferTrade(event Event) Event
	OnAcceptTrade(event Event) Event

	OnStartMatch(event Event) Event
	OnJoinMatch(event Event) Event
	OnSurrenderMatch(event Event) Event
	OnMakeMove(event Event) Event
}
