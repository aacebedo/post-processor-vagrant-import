// This package implements a post-processor for Packer that executes
// shell scripts locally.
package vagrantimport
import (
//    "bufio"
    "bytes"
// 	"io/ioutil"
	"errors"
	"fmt"
	"log"
	"os"
//	"os/exec"
	"strings"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	// Fields from config file
	ImportName        string `mapstructure:"import_name"`
	KeepInputArtifact bool   `mapstructure:"keep_input_artifact"`
	BoxFile string   `mapstructure:"box_file"`
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
    boxFile := p.config.BoxFile
    
	if _, err := os.Stat(boxFile); err != nil {
      return nil, false, fmt.Errorf("Unable to find box: %s", boxFile)
	}
	
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    stdout.Reset()
	stderr.Reset()
	ui.Say(strings.Join(artifact.Files(),","))
	ui.Say(fmt.Sprintf("Importing box: %s", boxFile))	
    cmd := fmt.Sprintf("Command: vagrant box add --name %s %s",importName,boxFile)
	ui.Say(cmd)

    //cmd := exec.Command("sh", "-c", "vagrant box add "+importName+)
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr
	//err = cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

    //if err != nil {
   //		return nil, false, fmt.Errorf("Error importing: %s", stderrString)
//	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	return artifact, keep, nil
}

    
