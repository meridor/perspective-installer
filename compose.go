package main

import "fmt"

type DockerComposeGenerator struct {
	
}

func (g DockerComposeGenerator) Config(answers Answers) string {
	panic("Not implemented")
}

func (g DockerComposeGenerator) Command(filename string) string {
	return fmt.Sprintf("docker-compose -f %s up", filename);
}