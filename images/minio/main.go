package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"resty.dev/v3"
)

const (
	timeFormat         = "2006-01-02T15:04:05Z"
	minioReleasePrefix = "RELEASE."
	initPage           = 1
	maxPerPage         = 100
	tarbarFileName     = "/tmp/mc.tar.gz"
)

var (
	// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#list-releases
	// max page is 100
	// Only the first 1000 results are available.
	repoURL                 = "https://api.github.com/repos/%s/releases?per_page=%d&page=%d"
	maybeTheLastError error = errors.New("maybe the last page, return code: 422")
	client                  = resty.New()
)

type GithubRelease struct {
	PublishedAt string `json:"published_at"`
	TarballURL  string `json:"tarball_url"`
	TagName     string `json:"tag_name"`
}

func main() {
	var (
		minioRelease = flag.String("minio_release", "", "The tag corresponding to a specific release on the MinIO release page.")
		dailyJob     = flag.Bool("daily_job", false, "Run daily-job logic (mutually exclusive with -minio_release)")
	)
	flag.Parse()

	// Mutual exclusion check
	if *minioRelease != "" && *dailyJob {
		log.Fatalf("❌ flag -minio_release and -daily_job are mutually exclusive")
	}

	switch {
	case *dailyJob:
		log.Println("✅ daily_job mode enabled")
		runDailyJob()
	case *minioRelease != "":
		log.Printf("✅ MinIO release version: %s\n", *minioRelease)
		runManuallyBuildImage(*minioRelease)
	default:
		log.Fatalln("❌ neither -minio_release nor -daily_job specified")
	}
	// https://api.github.com/repos/${REPO}/releases?per_page=${PER_PAGE}&page=${PAGE}
	defer client.Close()
}

func runManuallyBuildImage(s string) {
	if err := checkMinioTagExists(s); err != nil {
		log.Fatalf("❌ check minio tag %s failed: %s", s, err.Error())
	}
	tagName, err := getReleaseByDate("minio/mc", convertReleaseStr(s))
	if err != nil {
		log.Fatalf("❌ get the mc version closest to the release date of MinIO version %s failed: %s", s, err.Error())
	}
	log.Printf("✅ found the tag name: %s, write to /tmp/mc.txt", tagName)
	err = os.WriteFile("/tmp/mc.txt", []byte(tagName), 0644)
	if err != nil {
		panic(err)
	}
}
func checkMinioTagExists(tagName string) error {
	minioTagURL := "https://api.github.com/repos/minio/minio/git/ref/tags/%s"
	res, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		Get(fmt.Sprintf(minioTagURL, tagName))
	if err != nil {
		return err
	}
	resCode := res.StatusCode()
	if resCode != http.StatusOK {
		return fmt.Errorf("status code is not %d", http.StatusOK)
	}
	return nil
}

func getReleaseByDate(repo, dateStr string) (string, error) {
	inputDateUnix, err := parseDateStr(dateStr)
	if err != nil {
		log.Fatalf("❌ invalid input date string format: '%s', expect format: '%s'", dateStr, timeFormat)
	}
	releases := []GithubRelease{}
	startPage := initPage
	for {
		res, err := client.R().
			SetResult(&releases).
			EnableTrace().
			SetHeader("Accept", "application/vnd.github.v3+json").
			Get(fmt.Sprintf(repoURL, repo, maxPerPage, startPage))
		if err != nil {
			return "", err
		}
		resCode := res.StatusCode()
		if resCode == http.StatusUnprocessableEntity {
			return "", maybeTheLastError
		}
		if resCode != http.StatusOK {
			return "", fmt.Errorf("status code is not %d", resCode)
		}
		for _, item := range releases {
			tmp, err := parseDateStr(item.PublishedAt)
			if err != nil {
				log.Printf("❌ '%s' parse error: %s", item.PublishedAt, err)
				continue
			}
			// The first one that is earlier than the input date is the release we need.
			if tmp <= inputDateUnix {
				return item.TagName, nil
			}
		}
		startPage += maxPerPage
	}
	return "", nil
}

func parseDateStr(dateStr string) (int64, error) {
	t, err := time.Parse(timeFormat, dateStr)
	if err != nil {
		return -1, err
	}
	return t.Unix(), nil
}

func convertReleaseStr(s string) string {
	res := strings.TrimPrefix(s, minioReleasePrefix)
	parts := strings.SplitN(res, "T", 2)
	if len(parts) != 2 {
		log.Fatalf("❌ invalid input date minio release: '%s', expect format: '%s', "+
			"check the minio release page: https://github.com/minio/minio/releases",
			s, minioReleasePrefix+"."+timeFormat)
	}
	datePart := parts[0]
	timePart := parts[1]
	// Replace the '-' in timePart with ':'.
	timePart = strings.ReplaceAll(timePart, "-", ":")
	return datePart + "T" + timePart
}
