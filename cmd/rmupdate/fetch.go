package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/saleemrashid/rmupdate"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	fetchFilename = ""
	omahaURL      = ""

	requestParams = &rmupdate.RequestParams{}

	fetchCmd = &cobra.Command{
		Use:   "fetch [-o FILE] -p PLATFORM",
		Short: "Fetch the latest update payload from the reMarkable update server",
		Long: `Send an update request to the reMarkable update server and fetch the latest update payload.
The "extract" subcommand can be used to verify the payload and extract the root filesystem image.`,
		Args: cobra.NoArgs,
		Run:  fetchMain,
	}
)

func init() {
	fetchCmd.Flags().StringVarP(&fetchFilename, "out", "o", "", "output filename (default is from the payload URL)")
	fetchCmd.Flags().StringVar(&requestParams.AppID, "app-id", rmupdate.AppID, "GUID for firmware on Omaha server")
	fetchCmd.Flags().StringVar(&requestParams.Group, "group", "Prod", "release channel")
	fetchCmd.Flags().StringVar(&requestParams.MachineType, "machine", "armv7l", "CPU architecture")
	fetchCmd.Flags().StringVar(&requestParams.OSIdentifier, "os-id", "codex", "OS identifier")
	fetchCmd.Flags().StringVar(&requestParams.OSVersion, "os-version", "2.5.2", "OS version")
	fetchCmd.Flags().StringVarP(&requestParams.Platform, "platform", "p", "", "firmware platform (\"reMarkable\" or \"reMarkable2\")")
	fetchCmd.MarkFlagRequired("platform")
	fetchCmd.Flags().StringVar(&requestParams.ReleaseVersion, "release-version", "2.5.0.0", "current firmware version")
	fetchCmd.Flags().StringVar(&requestParams.SerialNumber, "serial-number", "RM100-000-00000", "serial number")
	fetchCmd.Flags().StringVar(&omahaURL, "omaha-url", rmupdate.OmahaURL, "URL for the Omaha server")

	rootCmd.AddCommand(fetchCmd)
}

func fetchMain(cmd *cobra.Command, args []string) {
	request := requestParams.Build()

	response, err := request.Send(omahaURL)
	if err != nil {
		panic(fmt.Errorf("sending update request: %w", err))
	}

	app := response.GetApp(requestParams.AppID)
	if app == nil {
		panic(fmt.Errorf("missing app in response: %s", requestParams.AppID))
	}

	u := app.UpdateCheck
	if len(u.Manifest.Packages) != 1 {
		panic("expected one package")
	}

	pkg := &u.Manifest.Packages[0]
	urls, err := u.PayloadURLs(pkg)
	if err != nil {
		panic(fmt.Errorf("parsing payload URLs: %w", err))
	}
	if len(urls) < 1 {
		panic("no payload URLs in response")
	}
	log.Printf("Payload URLs: %#v", urls)
	url := urls[0]

	action := u.Manifest.GetAction("postinstall")
	if action == nil {
		panic("missing postinstall action")
	}

	sha1Expected := pkg.SHA1
	sha256Expected := action.SHA256

	log.Printf("URL: %s", url)
	log.Printf("Size: %d bytes", pkg.Size)
	log.Printf("Expected SHA-1: %x", sha1Expected)
	log.Printf("Expected SHA-256: %x", sha256Expected)

	if fetchFilename == "" {
		fetchFilename = path.Base(pkg.Name)
	}
	log.Printf("Output filename: %s", fetchFilename)

	sha1Ctx := sha1.New()
	sha256Ctx := sha256.New()

	f, err := os.Create(fetchFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	progress := progressbar.DefaultBytes(pkg.Size, "Fetching")
	w := io.MultiWriter(bw, sha1Ctx, sha256Ctx, progress)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		panic(err)
	}

	sha1Computed := sha1Ctx.Sum(nil)
	sha256Computed := sha256Ctx.Sum(nil)

	log.Printf("Computed SHA-1: %x", sha1Computed)
	log.Printf("Computed SHA-256: %x", sha256Computed)

	if !bytes.Equal(sha1Computed, sha1Expected) {
		panic("SHA-1 mismatch")
	}
	if !bytes.Equal(sha256Computed, sha256Expected) {
		panic("SHA-256 mismatch")
	}
}
