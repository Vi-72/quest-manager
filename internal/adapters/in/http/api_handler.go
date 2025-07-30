package http

import (
	"context"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"
	"quest-manager/internal/pkg/errs"
)

type ApiHandler struct {
	createQuestHandler        commands.CreateQuestCommandHandler
	listQuestsHandler         queries.ListQuestsQueryHandler
	getQuestByIDHandler       queries.GetQuestByIDQueryHandler
	changeQuestStatusHandler  commands.ChangeQuestStatusCommandHandler
	searchQuestsByRadius      queries.SearchQuestsByRadiusQueryHandler
	listAssignedQuestsHandler queries.ListAssignedQuestsQueryHandler
	assignQuestHandler        commands.AssignQuestCommandHandler
}

func NewApiHandler(
	createQuestHandler commands.CreateQuestCommandHandler,
	listQuestsHandler queries.ListQuestsQueryHandler,
	getQuestByIDHandler queries.GetQuestByIDQueryHandler,
	changeQuestStatusHandler commands.ChangeQuestStatusCommandHandler,
	searchQuestsByRadius queries.SearchQuestsByRadiusQueryHandler,
	listAssignedQuestsHandler queries.ListAssignedQuestsQueryHandler,
	assignQuestHandler commands.AssignQuestCommandHandler,
) (*ApiHandler, error) {
	if createQuestHandler == nil {
		return nil, errs.NewValueIsRequiredError("createQuestHandler")
	}
	if listQuestsHandler == nil {
		return nil, errs.NewValueIsRequiredError("listQuestsHandler")
	}
	if getQuestByIDHandler == nil {
		return nil, errs.NewValueIsRequiredError("getQuestByIDHandler")
	}
	if changeQuestStatusHandler == nil {
		return nil, errs.NewValueIsRequiredError("changeQuestStatusHandler")
	}
	if searchQuestsByRadius == nil {
		return nil, errs.NewValueIsRequiredError("searchQuestsByRadius")
	}
	if listAssignedQuestsHandler == nil {
		return nil, errs.NewValueIsRequiredError("listAssignedQuestsHandler")
	}
	if assignQuestHandler == nil {
		return nil, errs.NewValueIsRequiredError("assignQuestHandler")
	}

	return &ApiHandler{
		createQuestHandler:        createQuestHandler,
		listQuestsHandler:         listQuestsHandler,
		getQuestByIDHandler:       getQuestByIDHandler,
		changeQuestStatusHandler:  changeQuestStatusHandler,
		searchQuestsByRadius:      searchQuestsByRadius,
		listAssignedQuestsHandler: listAssignedQuestsHandler,
		assignQuestHandler:        assignQuestHandler,
	}, nil
}

// StrictServerInterface implementation
func (a *ApiHandler) ListQuests(ctx context.Context, request servers.ListQuestsRequestObject) (servers.ListQuestsResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) CreateQuest(ctx context.Context, request servers.CreateQuestRequestObject) (servers.CreateQuestResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request servers.ListAssignedQuestsRequestObject) (servers.ListAssignedQuestsResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request servers.SearchQuestsByRadiusRequestObject) (servers.SearchQuestsByRadiusResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) GetQuestById(ctx context.Context, request servers.GetQuestByIdRequestObject) (servers.GetQuestByIdResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	// TODO: Implement
	return nil, nil
}

func (a *ApiHandler) ChangeQuestStatus(ctx context.Context, request servers.ChangeQuestStatusRequestObject) (servers.ChangeQuestStatusResponseObject, error) {
	// TODO: Implement
	return nil, nil
}
