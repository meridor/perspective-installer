package main

import (
	"flag"
	"fmt"
	"github.com/meridor/perspective-installer/generator"
	"github.com/meridor/perspective-installer/wizard"
	"os"
	"os/signal"
	"syscall"
)

var (
	dir string
	dryRun bool
)

func init() {
	flag.StringVar(&dir, "d", ".", "Directory to output generated config")
	flag.BoolVar(&dryRun, "dryRun", false, "Dry run (only show what is going to be done)")
	flag.Parse()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		for {
			<-sig
			fmt.Println("Interrupted.")
			os.Exit(1)
		}
	}()
}

func main() {
	generator.InitGenerators(dir, dryRun)
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: perspective-installer [-dir /path/to/directory] [-dryRun] generator1 generator2 ...")
		generatorNames := generator.GetNames()
		fmt.Println("The following generators are supported:")
		for _, gn := range generatorNames {
			fmt.Printf("* %s\n", gn)
		}
		os.Exit(1)
	}
	if (dryRun) {
		fmt.Println("Warning: working in \"dry run\" mode. No data will be actually written to disk.")
	}
	clusterConfig := wizard.RunWizards()
	if len(clusterConfig.Clouds) == 0 {
		fmt.Println("You skipped all clouds. Exiting.")
		os.Exit(1)
	}
	generator.RunGenerators(clusterConfig, flag.Args())
	os.Exit(0)
}
