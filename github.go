package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}

func createPR(info *asnInfo) error {

	base := "master"
	head := fmt.Sprintf("add-%d-%s", info.AsNumber, info.FirstIP)
	title := fmt.Sprintf("Add ASN %d", info.AsNumber)
	body := fmt.Sprintf("This PR adds ASN `%d` which is known as `%s` in `%s` to the list of mobile ASNs.\n\nThis AS has IPs from `%s` to `%s`.\n\nPlease verify this is not a false positive\n\n\n*This is an automated PR*", info.AsNumber, info.AsDescription, info.AsCountryCode, info.FirstIP, info.LastIP)

	_, _, err := client.PullRequests.Create(context.Background(), "innovate-technologies", "mobile-asn", &github.NewPullRequest{
		Base:  &base,
		Head:  &head,
		Title: &title,
		Body:  &body,
	})

	return err
}
