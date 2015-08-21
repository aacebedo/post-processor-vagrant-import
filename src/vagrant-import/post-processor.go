// This package implements a post-processor for Packer that executes
// shell scripts locally.
package vagrantimport
import (
  "bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"os/exec"
	"strings"
	"path/filepath"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
	"github.com/mitchellh/packer/post-processor/vagrant"
)

type Metadata struct {
    Name string `json:"name"`
    Provider string `json:"provider"`
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	// Fields from config file
	ImportName        string `mapstructure:"import_name"`
	KeepInputArtifact bool   `mapstructure:"keep_input_artifact"`
    ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	var errs *packer.MultiError
	if p.config.ImportName == "" {
		errs = packer.MultiErrorAppend(errs,
			errors.New("Import Name is missing."))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {

	keep := p.config.KeepInputArtifact

  importName := p.config.ImportName
  metaData := Metadata{importName,"virtualbox"}
    
  var stdout bytes.Buffer
  var stderr bytes.Buffer
  stdout.Reset()
	stderr.Reset()
	
	outputDir := filepath.Dir(artifact.Files()[0])
	marshalledJson,err := json.Marshal(metaData)
	if err != nil {
    panic(err)
  }
  
  ui.Say(fmt.Sprintf("Creating metadata file: %s", outputDir+"/metadata.json"))
	f, err := os.Create(outputDir+"/metadata.json")
  if err != nil {
    panic(err)
  }
  _ , err = f.Write(marshalledJson)
  if err != nil {
    panic(err)
  }
  f.Sync()
    
	vagrant.DirToBox("./output.box",filepath.Dir(artifact.Files()[0]),ui,0)
   
	ui.Say(fmt.Sprintf("Importing box into vagrant: %s", importName))
  if _, err := os.Stat("./output.box"); err != nil {
    return nil, false, fmt.Errorf("Unable to find box: ./output.box")
  }
  
  cmd := exec.Command("sh", "-c", fmt.Sprintf("vagrant box add --name %s ./output.box",importName))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())
  os.Remove("./output.box")
  if err != nil {
   		return nil, false, fmt.Errorf("Error importing: %s", stderrString)
	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	return artifact, keep, nil
}

    
