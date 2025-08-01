package http

import (
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
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
