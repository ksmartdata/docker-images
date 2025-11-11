package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/samber/lo"
)

const (
	// https://hub.docker.com/r/minio/minio/tags
	lastImageMinioPublished = "RELEASE.2025-09-07T16-13-09Z"
)

type GhcrImage struct {
	Metadata struct {
		PackageType string `json:"package_type"`
		Container   struct {
			Tags []string `json:"tags"`
		} `json:"container"`
	} `json:"metadata"`
}

func runDailyJob() {
	// get all minio releases and compare with the image we have published in the ghcr.io registry
	minioPublishedTags, err := getAllMinioReleasesUntilLast("minio/minio", convertReleaseStr(lastImageMinioPublished))
	if err != nil {
		log.Fatalf("❌ get all released tags from minio faild, err: %s", err)
	}
	ghcrImageTags, err := getAllPublishedImageTagsInGhcr()
	if err != nil {
		log.Fatalf("❌ get all released tags in ghcr faild, err: %s", err)
	}
	diff, _ := lo.Difference(minioPublishedTags, ghcrImageTags)
	log.Println("✅ the image to build:", diff)
	err = os.WriteFile("/tmp/daily-job.txt", []byte(strings.Join(diff, "\n")+"\n"), 0644)
	if err != nil {
		panic(err)
	}
}

func getAllMinioReleasesUntilLast(repo, dateStr string) ([]string, error) {
	inputDateUnix, err := parseDateStr(dateStr)
	if err != nil {
		log.Fatalf("❌ invalid input date string format: '%s', expect format: '%s'", dateStr, timeFormat)
	}
	releases := []GithubRelease{}
	startPage := initPage
	result := []string{}
	for {
		res, err := client.R().
			SetResult(&releases).
			EnableTrace().
			SetHeader("Accept", "application/vnd.github.v3+json").
			Get(fmt.Sprintf(repoURL, repo, maxPerPage, startPage))
		if err != nil {
			return nil, err
		}
		resCode := res.StatusCode()
		if resCode == http.StatusUnprocessableEntity {
			return nil, maybeTheLastError
		}
		if resCode != http.StatusOK {
			return nil, fmt.Errorf("status code is not %d", resCode)
		}
		for _, item := range releases {
			tmp, err := parseDateStr(item.PublishedAt)
			if err != nil {
				log.Printf("'%s' parse error: %s", item.PublishedAt, err)
				continue
			}
			result = append(result, item.TagName)
			// The first one that is earlier than the input date is the release we need.
			if tmp <= inputDateUnix {
				return result, nil
			}
		}
		startPage += maxPerPage
	}
	return nil, nil
}

func getAllPublishedImageTagsInGhcr() ([]string, error) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatal("❌ GH_TOKEN environment variable is not set")
	}
	url := "https://api.github.com/orgs/ksmartdata/packages/container/minio/versions?per_page=%d&page=%d"
	startPage := initPage
	res := []string{}
	for {
		ghchImages := []GhcrImage{}
		response, err := client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetResult(&ghchImages).
			Get(fmt.Sprintf(url, maxPerPage, startPage))
		if err != nil {
			return nil, err
		}
		resCode := response.StatusCode()
		if resCode != http.StatusOK {
			if resCode == http.StatusNotFound {
				return res, nil
			}
			return nil, fmt.Errorf("status code is %d", resCode)
		}
		if len(ghchImages) == 0 {
			// reached the last page
			break
		}
		for _, item := range ghchImages {
			if len(item.Metadata.Container.Tags) == 0 {
				continue
			}
			res = append(res, item.Metadata.Container.Tags...)
		}
		startPage++
	}
	return res, nil
}
