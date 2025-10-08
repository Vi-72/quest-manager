package authentication

import (
	"context"
	"net/http"

	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
	"quest-manager/tests/integration/mock"

	"github.com/google/uuid"
)

// AUTHENTICATION MIDDLEWARE TESTS - EDGE CASES
// Tests for various authentication edge cases and scenarios

// TestUserFromTokenForAssignQuest проверяет что user ID действительно берется из токена для assign endpoint
func (s *Suite) TestUserFromTokenForAssignQuest() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest using handler
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Setup - create auth client with specific user ID
	expectedUserID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	customAuthClient := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, expectedUserID)
	routerWithCustomAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(customAuthClient)

	// Act - assign quest (user ID will be taken from token)
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithCustomAuth, assignReq)

	// Assert - quest should be assigned to user from token
	assignResult := httpAssertions.QuestHTTPAssignedSuccessfully(assignResp, err)
	s.Assert().Equal(expectedUserID.String(), assignResult.Assignee, "Quest should be assigned to user from JWT token")

	// Verify in database
	updatedQuest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
	s.Require().NoError(err)
	s.Assert().NotNil(updatedQuest.Assignee, "Assignee should be set")
	s.Assert().Equal(expectedUserID, *updatedQuest.Assignee, "Database should have user ID from token")
}

// TestUserFromTokenForListAssignedQuests проверяет что user ID берется из токена для list assigned endpoint
func (s *Suite) TestUserFromTokenForListAssignedQuests() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Setup - create quests and assign to specific user
	userID1 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	userID2 := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	// Create and assign quest to user1
	quest1, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest1.ID(), userID1)
	s.Require().NoError(err)

	// Create and assign quest to user2
	quest2, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest2.ID(), userID2)
	s.Require().NoError(err)

	// Act - list assigned quests for user1 (using token)
	authClientForUser1 := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, userID1)
	routerForUser1 := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClientForUser1)

	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, routerForUser1, listReq)

	// Assert - should return only quests assigned to user1
	quests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	s.Assert().Len(quests, 1, "Should return exactly 1 quest for user1")
	s.Assert().Equal(quest1.ID().String(), quests[0].Id, "Should return quest assigned to user1")
	s.Assert().Equal(userID1.String(), quests[0].Assignee, "Returned quest should be assigned to user1")

	// Act - list assigned quests for user2 (using token)
	authClientForUser2 := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, userID2)
	routerForUser2 := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClientForUser2)

	listReq2 := casesteps.ListAssignedQuestsHTTPRequest()
	listResp2, err := casesteps.ExecuteHTTPRequest(ctx, routerForUser2, listReq2)

	// Assert - should return only quests assigned to user2
	quests2 := httpAssertions.QuestHTTPListSuccessfully(listResp2, err)
	s.Assert().Len(quests2, 1, "Should return exactly 1 quest for user2")
	s.Assert().Equal(quest2.ID().String(), quests2[0].Id, "Should return quest assigned to user2")
	s.Assert().Equal(userID2.String(), quests2[0].Assignee, "Returned quest should be assigned to user2")
}

// TestDifferentUsersCannotSeeEachOthersAssignedQuests проверяет изоляцию данных между пользователями
func (s *Suite) TestDifferentUsersCannotSeeEachOthersAssignedQuests() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Setup - create quests for two different users
	userA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	// Assign 3 quests to user A
	for i := 0; i < 3; i++ {
		quest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
		s.Require().NoError(err)
		_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest.ID(), userA)
		s.Require().NoError(err)
	}

	// Assign 2 quests to user B
	for i := 0; i < 2; i++ {
		quest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
		s.Require().NoError(err)
		_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest.ID(), userB)
		s.Require().NoError(err)
	}

	// Act & Assert - user A should see only their 3 quests
	authClientA := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, userA)
	routerForUserA := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClientA)

	listReqA := casesteps.ListAssignedQuestsHTTPRequest()
	listRespA, err := casesteps.ExecuteHTTPRequest(ctx, routerForUserA, listReqA)

	questsA := httpAssertions.QuestHTTPListSuccessfully(listRespA, err)
	s.Assert().Len(questsA, 3, "User A should see exactly 3 assigned quests")
	for _, q := range questsA {
		s.Assert().Equal(userA.String(), q.Assignee, "All quests should be assigned to user A")
	}

	// Act & Assert - user B should see only their 2 quests
	authClientB := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, userB)
	routerForUserB := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClientB)

	listReqB := casesteps.ListAssignedQuestsHTTPRequest()
	listRespB, err := casesteps.ExecuteHTTPRequest(ctx, routerForUserB, listReqB)

	questsB := httpAssertions.QuestHTTPListSuccessfully(listRespB, err)
	s.Assert().Len(questsB, 2, "User B should see exactly 2 assigned quests")
	for _, q := range questsB {
		s.Assert().Equal(userB.String(), q.Assignee, "All quests should be assigned to user B")
	}
}

// TestMultipleSequentialRequestsWithDifferentTokens проверяет что каждый запрос использует свой токен
func (s *Suite) TestMultipleSequentialRequestsWithDifferentTokens() {
	ctx := context.Background()

	// Setup - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// User 1 assigns the quest
	user1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	authClient1 := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, user1)
	router1 := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClient1)

	assignReq1 := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp1, err := casesteps.ExecuteHTTPRequest(ctx, router1, assignReq1)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, assignResp1.StatusCode, "User 1 should successfully assign quest")

	// Verify quest is assigned to user 1
	quest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
	s.Require().NoError(err)
	s.Assert().Equal(user1, *quest.Assignee, "Quest should be assigned to user 1")

	// User 2 tries to assign the same quest (should fail - quest already assigned)
	user2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	authClient2 := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, user2)
	router2 := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClient2)

	assignReq2 := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp2, err := casesteps.ExecuteHTTPRequest(ctx, router2, assignReq2)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusBadRequest, assignResp2.StatusCode, "Should fail - quest already assigned to different user")

	// Verify quest is still assigned to user 1
	questAfter, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
	s.Require().NoError(err)
	s.Assert().Equal(user1, *questAfter.Assignee, "Quest should still be assigned to user 1")
}

// TestAuthenticationPersistenceAcrossMultipleEndpoints проверяет что аутентификация работает последовательно
func (s *Suite) TestAuthenticationPersistenceAcrossMultipleEndpoints() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Setup - use consistent user ID across all requests
	consistentUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	authClient := mock.NewConfigurableAuthClient(mock.BehaviorSuccess, consistentUserID)
	router := s.TestDIContainer.NewHTTPRouterWithAuthClient(authClient)

	// 1. Create quest
	questData := testdatagenerators.SimpleQuestData(
		"Consistent User Quest",
		"Testing auth consistency",
		"medium",
		3,
		120,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)
	createReq := casesteps.CreateQuestHTTPRequest(questData)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, router, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// 2. Assign quest (to the same user from token)
	questID := uuid.UUID(createdQuest.Id) // Convert openapi_types.UUID to uuid.UUID
	assignReq := casesteps.AssignQuestHTTPRequest(questID)
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, router, assignReq)
	assignResult := httpAssertions.QuestHTTPAssignedSuccessfully(assignResp, err)
	s.Assert().Equal(consistentUserID.String(), assignResult.Assignee, "Should be assigned to consistent user")

	// 3. List assigned quests
	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, router, listReq)
	quests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	s.Assert().GreaterOrEqual(len(quests), 1, "Should have at least 1 assigned quest")

	found := false
	for _, q := range quests {
		if q.Id == createdQuest.Id {
			found = true
			s.Assert().Equal(consistentUserID.String(), q.Assignee, "Listed quest should be assigned to consistent user")
			break
		}
	}
	s.Assert().True(found, "Created quest should be in assigned quests list")
}

// TestEmptyUserIDFromToken проверяет обработку случая когда токен валиден но user ID пустой
func (s *Suite) TestEmptyOrInvalidUserIDFromAuthClient() {
	ctx := context.Background()

	// Setup - auth client that returns empty/nil user ID (simulating missing user in auth response)
	missingUserAuthClient := mock.NewConfigurableAuthClient(mock.BehaviorMissingUser, uuid.Nil)
	routerWithMissingUser := s.TestDIContainer.NewHTTPRouterWithAuthClient(missingUserAuthClient)

	// Act - try to list quests with token that has missing user
	listReq := casesteps.ListQuestsHTTPRequest("")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithMissingUser, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 for missing user in auth response")
}
