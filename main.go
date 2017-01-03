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
)

func init() {
	flag.StringVar(&dir, "-d", ".", "Directory to output generated config")
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
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: perspective-installer [-d /path/to/directory] generator1 generator2 ...")
		generatorNames := generator.GetNames()
		fmt.Println("The following generators are supported:")
		for _, gn := range generatorNames {
			fmt.Printf("* %s\n", gn)
		}
		os.Exit(1)
	}
	version, clouds := wizard.RunWizards()
	if len(clouds) == 0 {
		fmt.Println("You skipped all clouds. Exiting.")
		os.Exit(1)
	}
	generator.RunGenerators(dir, version, clouds, flag.Args())
	os.Exit(0)
}
