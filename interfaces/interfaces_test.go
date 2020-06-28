package interfaces

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInterfaces(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "interfaces")
}

/*
 * no tests because no statements;
 * this file is in order to measure test coverage
 */