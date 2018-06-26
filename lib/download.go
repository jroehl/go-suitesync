//+build !test

package lib

import (
	"net/http"
	"time"

	"github.com/cavaliercoder/grab"
)

var licenseCookie = &http.Cookie{Name: "oraclelicense",
	Value:    "accept-securebackup-cookie",
	Expires:  time.Now().Add(356 * 24 * time.Hour),
	HttpOnly: true}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.AddCookie(licenseCookie)
	return nil
}

func downloadFile(filepath string, uri string, redirect bool, ignoreStatus bool) string {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filepath, uri)
	// start download
	PrResultf("\nDownloading %v...\n", req.URL())

	if redirect {
		// nasty java download site redirect
		client.HTTPClient.CheckRedirect = redirectPolicyFunc
	}

	req.NoResume = ignoreStatus
	req.IgnoreBadStatusCodes = ignoreStatus

	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			PrNoticef("\r   transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	PrNoticef("\r  transferred %v / %v bytes (100.00%%)\n", resp.Size, resp.Size)

	// check for errors
	if err := resp.Err(); err != nil {
		PrFatalf("Download failed: %v\n", err)
	}

	PrResultf("Download saved to ./%v \n\n", resp.Filename)
	return resp.Filename
}
