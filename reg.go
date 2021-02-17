package reg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beevik/etree"
	"github.com/noshto/dsig/pkg/safenet"
	"github.com/noshto/sep"
)

// Params represents collection of parameters needed for IIC function
type Params struct {
	SafenetConfig *safenet.Config
	SepConfig     *sep.Config
	InFile        string
	OutFile       string
}

// Register sends given request to SEP service
func Register(params *Params) error {

	// Initialize Signer
	signer := safenet.SafeNet{}
	if err := signer.Initialize(params.SafenetConfig); err != nil {
		return err
	}
	defer signer.Finalize()

	client, err := signer.NewClient()
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(params.InFile)
	if err != nil {
		return err
	}

	var url string
	if params.SepConfig.Environment == sep.TEST {
		url = sep.TestingURL
	} else if params.SepConfig.Environment == sep.PROD {
		url = sep.ProductionURL
	} else {
		return fmt.Errorf("unknown environment")
	}

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(buf),
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(buf); err != nil {
		return err
	}
	doc.IndentTabs()
	doc.Root().SetTail("")

	if err = doc.WriteToFile(params.OutFile); err != nil {
		return err
	}
	return nil
}
