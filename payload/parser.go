package payload

//go:generate protoc --go_out=paths=source_relative:. update_metadata.proto

import (
	"bytes"
	"compress/bzip2"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

var (
	Magic = []byte("CrAU")
)

const (
	manifestOffset = 20
)

type Parser struct {
	DataReader io.ReaderAt
	Manifest   *DeltaArchiveManifest
	Signatures []*Signatures_Signature
	SHA256     []byte
}

func NewParser(r io.ReaderAt) (*Parser, error) {
	buf := make([]byte, manifestOffset)
	if _, err := r.ReadAt(buf, 0); err != nil {
		return nil, err
	}

	magic := buf[:4]
	if !bytes.Equal(magic, Magic) {
		return nil, fmt.Errorf("bad magic: %q", magic)
	}

	manifestSize := int64(binary.BigEndian.Uint64(buf[12:]))
	manifestBlob := make([]byte, manifestSize)
	if _, err := r.ReadAt(manifestBlob, manifestOffset); err != nil {
		return nil, err
	}

	manifest := &DeltaArchiveManifest{}
	if err := proto.Unmarshal(manifestBlob, manifest); err != nil {
		return nil, err
	}

	dataOffset := manifestOffset + manifestSize
	signaturesOffset := int64(manifest.GetSignaturesOffset())
	signaturesBlob := make([]byte, manifest.GetSignaturesSize())
	if _, err := r.ReadAt(signaturesBlob, dataOffset+signaturesOffset); err != nil {
		return nil, err
	}

	signatures := &Signatures{}
	if err := proto.Unmarshal(signaturesBlob, signatures); err != nil {
		return nil, err
	}

	dataReader := io.NewSectionReader(r, dataOffset, signaturesOffset)

	h := sha256.New()
	h.Write(buf)
	h.Write(manifestBlob)
	if _, err := io.Copy(h, dataReader); err != nil {
		return nil, err
	}

	return &Parser{
		DataReader: dataReader,
		Manifest:   manifest,
		Signatures: signatures.Signatures,
		SHA256:     h.Sum(nil),
	}, nil
}

func (p *Parser) Verify(pub *rsa.PublicKey) error {
	err := errors.New("no signatures")
	for _, signature := range p.Signatures {
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, p.SHA256, signature.Data)
		if err == nil {
			break
		}
	}
	return err
}

func (p *Parser) GetData(operation *InstallOperation) ([]byte, error) {
	data := make([]byte, operation.GetDataLength())
	if _, err := p.DataReader.ReadAt(data, int64(operation.GetDataOffset())); err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(data)
	if !bytes.Equal(hashed[:], operation.DataSha256Hash) {
		return nil, errors.New("SHA-256 mismatch")
	}

	return data, nil
}

func (p *Parser) Execute(operations []*InstallOperation, installInfo *InstallInfo, w io.Writer) error {
	h := sha256.New()
	mw := io.MultiWriter(w, h)

	blockSize := int(p.Manifest.GetBlockSize())
	endBlock := uint64(0)

	// The final block might require padding at the end
	blockPadding := 0

	for _, operation := range operations {
		if blockPadding != 0 {
			return errors.New("block padding only supported at the end")
		}

		bz := false
		switch ty := operation.GetType(); ty {
		case InstallOperation_REPLACE:
		case InstallOperation_REPLACE_BZ:
			bz = true
		default:
			return fmt.Errorf("unsupported operation type: %v", ty)
		}

		data, err := p.GetData(operation)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)
		if bz {
			tmp := &bytes.Buffer{}
			if _, err := io.Copy(tmp, bzip2.NewReader(buf)); err != nil {
				return err
			}
			buf = tmp
		}

		for _, extent := range operation.DstExtents {
			if extent.GetStartBlock() != endBlock {
				return errors.New("extents must be contiguous")
			}
			numBlocks := extent.GetNumBlocks()
			endBlock += numBlocks

			blockPadding = int(numBlocks) * blockSize - buf.Len()
			if blockPadding < 0 {
				return errors.New("data larger than number of blocks")
			}

			if _, err := buf.WriteTo(mw); err != nil {
				return err
			}
		}
	}

	if blockPadding != 0 {
		// This should be written to the output file, but not included in the hash
		if _, err := w.Write(make([]byte, blockPadding)); err != nil {
			return err
		}
	}

	if !bytes.Equal(h.Sum(nil), installInfo.Hash) {
		return errors.New("SHA-256 mismatch")
	}

	return nil
}
