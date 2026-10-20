package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"resty.dev/v3"
)

func serveReleasePages(t *testing.T, pages map[int][]GithubRelease) func() []int {
	t.Helper()
	var mu sync.Mutex
	var requested []int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || r.URL.Query().Get("per_page") != strconv.Itoa(maxPerPage) {
			http.Error(w, "invalid pagination query", http.StatusBadRequest)
			return
		}
		mu.Lock()
		requested = append(requested, page)
		mu.Unlock()
		releases, ok := pages[page]
		if !ok {
			// Bound baseline failures: an unexpected page must never loop forever.
			http.Error(w, fmt.Sprintf("unexpected page %d", page), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(releases); err != nil {
			t.Errorf("encode releases: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	originalClient, originalURL := client, repoURL
	testClient := resty.New().SetTimeout(2 * time.Second)
	client = testClient
	repoURL = server.URL + "/repos/%s/releases?per_page=%d&page=%d"
	t.Cleanup(func() {
		testClient.Close()
		client, repoURL = originalClient, originalURL
	})
	return func() []int {
		mu.Lock()
		defer mu.Unlock()
		return append([]int(nil), requested...)
	}
}

func newerReleasePage() []GithubRelease {
	releases := make([]GithubRelease, maxPerPage)
	date := time.Date(2025, 9, 10, 12, 0, 0, 0, time.UTC)
	for i := range releases {
		releases[i].TagName = minioReleasePrefix + date.Add(-time.Duration(i)*time.Second).Format("2006-01-02T15-04-05Z")
	}
	return releases
}

func TestGetReleaseByDatePagination(t *testing.T) {
	const matchingTag = "RELEASE.2025-09-07T12-00-00Z"
	for _, tc := range []struct {
		name      string
		pages     map[int][]GithubRelease
		wantTag   string
		wantError bool
		wantPages []int
	}{
		{
			name: "match on second page",
			pages: map[int][]GithubRelease{
				1: newerReleasePage(),
				2: {{TagName: matchingTag}},
			},
			wantTag:   matchingTag,
			wantPages: []int{1, 2},
		},
		{
			name:      "empty first page",
			pages:     map[int][]GithubRelease{1: {}},
			wantError: true,
			wantPages: []int{1},
		},
		{
			name: "no match before empty final page",
			pages: map[int][]GithubRelease{
				1: newerReleasePage(),
				2: {},
			},
			wantError: true,
			wantPages: []int{1, 2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requested := serveReleasePages(t, tc.pages)
			tag, err := getReleaseByDate("minio/mc", "2025-09-08T00:00:00Z")
			if (err != nil) != tc.wantError || tag != tc.wantTag {
				t.Errorf("getReleaseByDate() = %q, %v; want tag %q, error %v", tag, err, tc.wantTag, tc.wantError)
			}
			if got := requested(); !reflect.DeepEqual(got, tc.wantPages) {
				t.Errorf("requested pages = %v, want %v", got, tc.wantPages)
			}
		})
	}
}

func TestGetAllMinioReleasesUntilLastPagination(t *testing.T) {
	const recentTag = "RELEASE.2025-09-08T12-00-00Z"
	const cutoffTag = "RELEASE.2025-09-07T12-00-00Z"
	firstPage := newerReleasePage()
	firstPageTags := make([]string, len(firstPage))
	for i, release := range firstPage {
		firstPageTags[i] = release.TagName
	}
	for _, tc := range []struct {
		name      string
		pages     map[int][]GithubRelease
		wantTags  []string
		wantPages []int
	}{
		{
			name: "cutoff on second page",
			pages: map[int][]GithubRelease{
				1: firstPage,
				2: {{TagName: recentTag}, {TagName: cutoffTag}},
			},
			wantTags:  append(append([]string{}, firstPageTags...), recentTag),
			wantPages: []int{1, 2},
		},
		{
			name:      "empty first page",
			pages:     map[int][]GithubRelease{1: {}},
			wantTags:  []string{},
			wantPages: []int{1},
		},
		{
			name: "preserve releases when final page is empty",
			pages: map[int][]GithubRelease{
				1: firstPage,
				2: {},
			},
			wantTags:  firstPageTags,
			wantPages: []int{1, 2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requested := serveReleasePages(t, tc.pages)
			tags, err := getAllMinioReleasesUntilLast("minio/minio", convertReleaseStr(cutoffTag))
			if err != nil {
				t.Errorf("getAllMinioReleasesUntilLast() error: %v", err)
			}
			if !reflect.DeepEqual(tags, tc.wantTags) {
				t.Errorf("returned %d tags, want %d with matching contents", len(tags), len(tc.wantTags))
			}
			if got := requested(); !reflect.DeepEqual(got, tc.wantPages) {
				t.Errorf("requested pages = %v, want %v", got, tc.wantPages)
			}
		})
	}
}
