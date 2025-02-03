package tg

import "gymnote/internal/entity"

func (a *API) setUserState(userID string, state entity.UserState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.userStates[userID] = state
}

func (a *API) getUserState(userID string) entity.UserState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.userStates[userID]
}

func (a *API) clearUserState(userID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.userStates, userID)
}
