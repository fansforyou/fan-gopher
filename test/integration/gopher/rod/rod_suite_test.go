package rod_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRod(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rod Suite")
}
