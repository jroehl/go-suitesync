package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PrHeaderf print underlined header
func TestPrHeaderf(t *testing.T) {
	assert.NotPanics(t, func() { PrHeaderf("%s", "foobar") })
}

func TestPrNoticef(t *testing.T) {
	assert.NotPanics(t, func() { PrNoticef("foobar") })
}

func TestPrWarnf(t *testing.T) {
	assert.NotPanics(t, func() { PrWarnf("") })
}

func TestPrFatalf(t *testing.T) {
	// can't be tested easily due to os.Exit
	// assert.Panics(t, func() { PrFatalf("") })
}

func TestPrResultf(t *testing.T) {
	assert.NotPanics(t, func() { PrResultf("") })
}

func TestPrettyList(t *testing.T) {
	assert.NotPanics(t, func() { PrettyList("", []string{"foo", "bar"}) })
}

func TestPrettyHash(t *testing.T) {
	assert.NotPanics(t, func() { PrettyHash("", []Hash{Hash{}}) })
}

func TestPrintResponse(t *testing.T) {
	assert.NotPanics(t, func() { PrintResponse("", []DeleteResult{DeleteResult{Code: "1"}, DeleteResult{Code: "2"}}) })
}

func TestPrettyTable(t *testing.T) {
	assert.NotPanics(t, func() { PrettyTable("", []string{"foo", "bar"}, [][]string{[]string{"bar", "foo"}}) })
}

func TestPrHeaderfCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrHeaderf("%s", "foobar") })
}

func TestPrNoticefCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrNoticef("foobar") })
}

func TestPrWarnfCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrWarnf("") })
}

func TestPrResultfCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrResultf("") })
}

func TestPrettyListCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrettyList("", []string{"foo", "bar"}) })
}

func TestPrettyHashCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrettyHash("", []Hash{Hash{}}) })
}

func TestPrintResponseCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrintResponse("", []DeleteResult{DeleteResult{Code: "1"}, DeleteResult{Code: "2"}}) })
}

func TestPrettyTableCI(t *testing.T) {
	os.Setenv("CI", "true")
	assert.NotPanics(t, func() { PrettyTable("", []string{"foo", "bar"}, [][]string{[]string{"bar", "foo"}}) })
}
