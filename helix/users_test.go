package helix

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestGetUsers(t *testing.T) {
	ctx := context.Background()
	in := &GetUsersRequest{
		RequestOptions: requestOptions(),
		IDs:            []string{"1", "2", "3"},
		Logins:         []string{"one", "two", "three"},
	}

	testUser := func(id string) *User {
		return &User{
			ID:              id,
			Login:           "id" + id,
			DisplayName:     "id" + id,
			Type:            "",
			BroadcasterType: "partner",
			Description:     "description",
			ProfileImageURL: "https://twitch.tv",
			OfflineImageURL: "https://twitch.tv",
			Email:           "",
			CreatedAt:       time.Now().UTC(),
		}
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, usersPath+"?id=1&id=2&id=3&login=one&login=two&login=three", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetUsersResponse{
			Data: []*User{
				testUser("1"),
				testUser("2"),
				testUser("3"),
			},
		})
	})

	resp, err := c.GetUsers(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(resp.Data))
}

func TestGetUserFollows(t *testing.T) {
	ctx := context.Background()
	in := &GetUserFollowsRequest{
		RequestOptions: requestOptions(),
		FromID:         "fromID",
		ToID:           "toID",
		After:          "after",
		First:          10,
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, in.FromID, req.URL.Query().Get("from_id"))
		assert.Equal(t, in.ToID, req.URL.Query().Get("to_id"))
		assert.Equal(t, in.After, req.URL.Query().Get("after"))
		assert.Equal(t, strconv.Itoa(in.First), req.URL.Query().Get("first"))
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, userFollowsPath, testutils.WithoutQuery(req.URL))
		return testutils.JSONResponse(t, http.StatusOK, &GetUserFollowsResponse{
			Data: []*UserFollow{
				{
					FromID:     in.FromID,
					FromLogin:  "fromLogin",
					FromName:   "fromName",
					ToID:       in.ToID,
					ToLogin:    "toLogin",
					ToName:     "toName",
					FollowedAt: time.Now().UTC(),
				},
			},
		})
	})

	resp, err := c.GetUserFollows(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, in.FromID, resp.Data[0].FromID)
}

func TestGetUserBlocks(t *testing.T) {
	ctx := context.Background()
	in := &GetUserBlocksRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "id",
		After:          "after",
		First:          10,
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, in.BroadcasterID, req.URL.Query().Get("broadcaster_id"))
		assert.Equal(t, in.After, req.URL.Query().Get("after"))
		assert.Equal(t, strconv.Itoa(in.First), req.URL.Query().Get("first"))
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, userBlocksPath, testutils.WithoutQuery(req.URL))
		return testutils.JSONResponse(t, http.StatusOK, &GetUserBlocksResponse{
			Data: []*UserBlock{
				{
					UserID:      "id",
					UserLogin:   "login",
					DisplayName: "name",
				},
			},
		})
	})

	resp, err := c.GetUserBlocks(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "id", resp.Data[0].UserID)
}

func TestBlockUser(t *testing.T) {
	ctx := context.Background()
	in := &BlockUserRequest{
		RequestOptions: requestOptions(),

		TargetUserID:  "targetUserID",
		SourceContext: "chat",
		Reason:        "some reason",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, in.TargetUserID, req.URL.Query().Get("target_user_id"))
		assert.Equal(t, in.SourceContext, req.URL.Query().Get("source_context"))
		assert.Equal(t, in.Reason, req.URL.Query().Get("reason"))
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Equal(t, userBlocksPath, testutils.WithoutQuery(req.URL))
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.BlockUser(ctx, in))
}

func TestUnblockUser(t *testing.T) {
	ctx := context.Background()
	in := &UnblockUserRequest{
		RequestOptions: requestOptions(),
		TargetUserID:   "targetUserID",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, in.TargetUserID, req.URL.Query().Get("target_user_id"))
		assert.Equal(t, http.MethodDelete, req.Method)
		assert.Equal(t, userBlocksPath, testutils.WithoutQuery(req.URL))
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.UnblockUser(ctx, in))
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	in := &UpdateUserRequest{
		RequestOptions: requestOptions(),
		Description:    "description",
	}

	testUser := func(id string) *User {
		return &User{
			ID:              id,
			Login:           "id" + id,
			DisplayName:     "id" + id,
			Type:            "",
			BroadcasterType: "partner",
			Description:     "description",
			ProfileImageURL: "https://twitch.tv",
			OfflineImageURL: "https://twitch.tv",
			Email:           "",
			CreatedAt:       time.Now().UTC(),
		}
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Equal(t, usersPath+"?description="+in.Description, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetUsersResponse{
			Data: []*User{
				testUser("1"),
			},
		})
	})

	resp, err := c.UpdateUser(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestGetUserExtensions(t *testing.T) {
	ctx := context.Background()
	in := &GetUserExtensionsRequest{
		RequestOptions: requestOptions(),
	}

	testExtension := func(id string) *UserExtension {
		return &UserExtension{
			ID:          id,
			Version:     "0.0.11",
			Name:        "name",
			CanActivate: true,
			Type:        []string{"panel"},
		}
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, userExtensionsListPath, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetUserExtensionsResponse{
			Data: []*UserExtension{
				testExtension("1"),
				testExtension("2"),
			},
		})
	})

	resp, err := c.GetUserExtensions(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Data))
}

func TestGetUserActiveExtensions(t *testing.T) {
	ctx := context.Background()
	in := &GetUserActiveExtensionsRequest{
		RequestOptions: requestOptions(),
		UserID:         "userID",
	}

	testExtension := func(id string) *UserActiveExtension {
		return &UserActiveExtension{
			Active:  true,
			ID:      id,
			Name:    id,
			Version: "0.0.1",
		}
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, userExtensionsPath+"?user_id="+in.UserID, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetUserActiveExtensionsResponse{
			Data: UserActiveExtensionsProperties{
				Component: map[string]*UserActiveExtension{
					"1": testExtension("1"),
				},
				Panel: map[string]*UserActiveExtension{
					"2": testExtension("2"),
					"3": testExtension("3"),
				},
				Overlay: map[string]*UserActiveExtension{
					"4": testExtension("4"),
					"5": testExtension("5"),
					"6": testExtension("6"),
				},
			},
		})
	})

	resp, err := c.GetUserActiveExtensions(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data.Component))
	assert.Equal(t, 2, len(resp.Data.Panel))
	assert.Equal(t, 3, len(resp.Data.Overlay))
}

func TestUpdateUserExtensions(t *testing.T) {
	ctx := context.Background()
	testExtension := func(id string) *UserActiveExtension {
		return &UserActiveExtension{
			Active:  true,
			ID:      id,
			Name:    id,
			Version: "0.0.1",
		}
	}

	data := UserActiveExtensionsProperties{
		Component: map[string]*UserActiveExtension{
			"1": testExtension("1"),
		},
		Panel: map[string]*UserActiveExtension{
			"2": testExtension("2"),
			"3": testExtension("3"),
		},
		Overlay: map[string]*UserActiveExtension{
			"4": testExtension("4"),
			"5": testExtension("5"),
			"6": testExtension("6"),
		},
	}
	in := &UpdateUserExtensionsRequest{
		RequestOptions: requestOptions(),
		Body:           &data,
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Equal(t, userExtensionsPath, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &UpdateUserExtensionsResponse{
			Data: data,
		})
	})

	resp, err := c.UpdateUserExtensions(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data.Component))
	assert.Equal(t, 2, len(resp.Data.Panel))
	assert.Equal(t, 3, len(resp.Data.Overlay))
}
