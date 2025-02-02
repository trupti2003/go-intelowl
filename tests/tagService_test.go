package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/intelowlproject/go-intelowl/gointelowl"
)

// * Test for TagService.List method!
func TestTagServiceList(t *testing.T) {
	// * table test case
	testCases := make(map[string]TestData)
	testCases["simple"] = TestData{
		Input: nil,
		Data: `[
			{"id": 1,"label": "TEST1","color": "#1c71d8"},
			{"id": 2,"label": "TEST2","color": "#1c71d7"},
			{"id": 3,"label": "TEST3","color": "#1c71d6"}
		]`,
		StatusCode: http.StatusOK,
		Want: []gointelowl.Tag{
			{
				ID:    1,
				Label: "TEST1",
				Color: "#1c71d8",
			},
			{
				ID:    2,
				Label: "TEST2",
				Color: "#1c71d7",
			},
			{
				ID:    3,
				Label: "TEST3",
				Color: "#1c71d6",
			},
		},
	}
	for name, testCase := range testCases {
		//* Subtest
		t.Run(name, func(t *testing.T) {
			testServer := NewTestServer(&testCase)
			defer testServer.Close()
			client := NewTestIntelOwlClient(testServer.URL)
			ctx := context.Background()
			gottenTagList, err := client.TagService.List(ctx)
			if err != nil {
				t.Fatalf("Error listing tags: %v", err)
			}
			diff := cmp.Diff(testCase.Want, (*gottenTagList))
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestTagServiceGet(t *testing.T) {
	// *table test case
	testCases := make(map[string]TestData)
	testCases["simple"] = TestData{
		Input:      1,
		Data:       `{"id": 1,"label": "TEST","color": "#1c71d8"}`,
		StatusCode: http.StatusOK,
		Want: &gointelowl.Tag{
			ID:    1,
			Label: "TEST",
			Color: "#1c71d8",
		},
	}
	testCases["cantFind"] = TestData{
		Input:      9000,
		Data:       `{"detail": "Not found."}`,
		StatusCode: http.StatusNotFound,
		Want: &gointelowl.IntelOwlError{
			StatusCode: http.StatusNotFound,
			Message:    `{"detail": "Not found."}`,
		},
	}
	for name, testCase := range testCases {
		//* Subtest
		t.Run(name, func(t *testing.T) {
			testServer := NewTestServer(&testCase)
			defer testServer.Close()
			client := NewTestIntelOwlClient(testServer.URL)
			ctx := context.Background()
			id, ok := testCase.Input.(int)
			if ok {
				tagId := uint64(id)
				gottenTag, err := client.TagService.Get(ctx, tagId)
				if testCase.StatusCode < http.StatusOK || testCase.StatusCode >= http.StatusBadRequest {
					diff := cmp.Diff(testCase.Want, err, cmpopts.IgnoreFields(gointelowl.IntelOwlError{}, "Response"))
					if diff != "" {
						t.Fatalf(diff)
					}
				} else {
					diff := cmp.Diff(testCase.Want, gottenTag)
					if diff != "" {
						t.Fatalf(diff)
					}
				}
			} else {
				t.Fatalf("Casting failed!")
			}
		})
	}
}

func TestTagServiceCreate(t *testing.T) {
	// *table test case
	testCases := make(map[string]TestData)
	testCases["simple"] = TestData{
		Input: gointelowl.TagParams{
			Label: "TEST TAG",
			Color: "#fffff",
		},
		Data:       `{"id": 1,"label": "TEST TAG","color": "#fffff"}`,
		StatusCode: http.StatusOK,
		Want: &gointelowl.Tag{
			ID:    1,
			Label: "TEST TAG",
			Color: "#fffff",
		},
	}
	testCases["duplicate"] = TestData{
		Input: gointelowl.TagParams{
			Label: "TEST TAG",
			Color: "#fffff",
		},
		Data:       `{"label":["tag with this label already exists."]}`,
		StatusCode: http.StatusBadRequest,
		Want: &gointelowl.IntelOwlError{
			StatusCode: http.StatusBadRequest,
			Message:    `{"label":["tag with this label already exists."]}`,
		},
	}
	for name, testCase := range testCases {
		//* Subtest
		t.Run(name, func(t *testing.T) {
			testServer := NewTestServer(&testCase)
			defer testServer.Close()
			client := NewTestIntelOwlClient(testServer.URL)
			ctx := context.Background()
			tagParams, ok := testCase.Input.(gointelowl.TagParams)
			if ok {
				gottenTag, err := client.TagService.Create(ctx, &tagParams)
				if testCase.StatusCode < http.StatusOK || testCase.StatusCode >= http.StatusBadRequest {
					diff := cmp.Diff(testCase.Want, err, cmpopts.IgnoreFields(gointelowl.IntelOwlError{}, "Response"))
					if diff != "" {
						t.Fatalf(diff)
					}
				} else {
					diff := cmp.Diff(testCase.Want, gottenTag)
					if diff != "" {
						t.Fatalf(diff)
					}
				}
			} else {
				t.Fatalf("Casting failed!")
			}
		})
	}
}

func TestTagServiceUpdate(t *testing.T) {
	// *table test case
	testCases := make(map[string]TestData)
	testCases["simple"] = TestData{
		Input: gointelowl.Tag{
			ID:    1,
			Label: "UPDATED TEST TAG",
			Color: "#f4",
		},
		Data:       `{"id": 1,"label": "UPDATED TEST TAG","color": "#f4"}`,
		StatusCode: http.StatusOK,
		Want: &gointelowl.Tag{
			ID:    1,
			Label: "UPDATED TEST TAG",
			Color: "#f4",
		},
	}
	for name, testCase := range testCases {
		//* Subtest
		t.Run(name, func(t *testing.T) {
			testServer := NewTestServer(&testCase)
			defer testServer.Close()
			client := NewTestIntelOwlClient(testServer.URL)
			ctx := context.Background()
			tag, ok := testCase.Input.(gointelowl.Tag)
			if ok {
				gottenTag, err := client.TagService.Update(ctx, tag.ID, &gointelowl.TagParams{
					Label: tag.Label,
					Color: tag.Color,
				})
				if testCase.StatusCode < http.StatusOK || testCase.StatusCode >= http.StatusBadRequest {
					diff := cmp.Diff(testCase.Want, err, cmpopts.IgnoreFields(gointelowl.IntelOwlError{}, "Response"))
					if diff != "" {
						t.Fatalf(diff)
					}
				} else {
					diff := cmp.Diff(testCase.Want, gottenTag)
					if diff != "" {
						t.Fatalf(diff)
					}
				}
			}
		})
	}
}

func TestTagServiceDelete(t *testing.T) {
	// *table test case
	testCases := make(map[string]TestData)
	testCases["simple"] = TestData{
		Input:      1,
		Data:       "",
		StatusCode: http.StatusNoContent,
		Want:       true,
	}
	for name, testCase := range testCases {
		//* Subtest
		t.Run(name, func(t *testing.T) {
			testServer := NewTestServer(&testCase)
			defer testServer.Close()
			client := NewTestIntelOwlClient(testServer.URL)
			ctx := context.Background()
			id, ok := testCase.Input.(int)
			if ok {
				tagId := uint64(id)
				isDeleted, err := client.TagService.Delete(ctx, tagId)
				if testCase.StatusCode < http.StatusOK || testCase.StatusCode >= http.StatusBadRequest {
					diff := cmp.Diff(testCase.Want, err, cmpopts.IgnoreFields(gointelowl.IntelOwlError{}, "Response"))
					if diff != "" {
						t.Fatalf(diff)
					}
				} else {
					diff := cmp.Diff(testCase.Want, isDeleted)
					if diff != "" {
						t.Fatalf(diff)
					}
				}
			}
		})
	}
}
