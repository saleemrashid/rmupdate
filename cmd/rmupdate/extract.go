package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/saleemrashid/rmupdate"
	"github.com/saleemrashid/rmupdate/payload"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	extractInput     = ""
	extractOutput    = ""
	extractPublicKey = ""

	extractCmd = &cobra.Command{
		Use:   "extract -i FILE -o FILE",
		Short: "Extract the update payload",
		Args:  cobra.NoArgs,
		Run:   extractMain,
	}
)

func init() {
	extractCmd.Flags().StringVarP(&extractInput, "in", "i", "", "payload filename")
	extractCmd.MarkFlagRequired("in")
	extractCmd.Flags().StringVarP(&extractOutput, "out", "o", "", "output filename")
	extractCmd.MarkFlagRequired("out")
	extractCmd.Flags().StringVar(&extractPublicKey, "public-key", "", "custom RSA public key for payload verification (default is official reMarkable public key)")

	rootCmd.AddCommand(extractCmd)
}

func loadPublicKey(filename string) (*rsa.PublicKey, error) {
	var data []byte
	if filename == "" {
		data = []byte(rmupdate.PublicKey)
	} else {
		tmp, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("loading public key: %w", err)
		}
		data = tmp
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("loading public key: invalid PEM")
	}

	tmp, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("load public key: %w", err)
	}

	pub, ok := tmp.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("load public key: not an RSA key")
	}

	return pub, nil
}

func extractMain(cmd *cobra.Command, args []string) {
	pub, err := loadPublicKey(extractPublicKey)
	if err != nil {
		panic(err)
	}

	r, err := os.Open(extractInput)
	if err != nil {
		panic(err)
	}

	log.Printf("Parsing manifest")
	parser, err := payload.NewParser(r)
	if err != nil {
		panic(fmt.Errorf("parsing manifest: %w", err))
	}

	log.Printf("Verifying RSA signature")
	if err := parser.Verify(pub); err != nil {
		panic(fmt.Errorf("verifying signature: %w", err))
	}

	installInfo := parser.Manifest.GetNewPartitionInfo()
	if installInfo == nil {
		panic("extracting payload: missing install info")
	}

	f, err := os.Create(extractOutput)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	progress := progressbar.DefaultBytes(int64(installInfo.GetSize()), "Extracting")
	w := io.MultiWriter(bw, progress)

	log.Printf("Extracting payload")
	if err := parser.Execute(parser.Manifest.PartitionOperations, installInfo, w); err != nil {
		panic(fmt.Errorf("extracting payload: %w", err))
	}
}
