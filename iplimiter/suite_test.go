package iplimiter

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPeripcounter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Peripcounter Suite")
}
