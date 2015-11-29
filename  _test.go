package sugar_test
// 
// import (
// 	"testing"
// 	"github.com/marksalpeter/sugar"
// )
// 
// // create a writer that streams all of the output of the tests
// // to a channel for us to examine
// type OutputChecker struct {}
// func (cw *channelWriter)  Write(p []byte) (n int, err error) {
// 	cw.Channel <- p
// 	return len(p), nil
// }
// var cw channelWriter
//
//
// func TestVerboseSugar(t *testing.T) {
//
// 	// TODO: set verbose flag
//
// 	testT := &testing.T{}
//
// 	sugar.New(testT).
//
// 	Must("this must pass", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return true
// 	}).
//
// 	Assert("this passes and continues", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return true
// 	}).
//
// 	Assert("this fails but continues", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return false
// 	}).
//
// 	Title("this must add a title").
//
// 	Warn("this passes and continues", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return true
// 	}).
//
// 	Warn("this passes and continues but warns", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return false
// 	}).
//
// 	Must("this must fail now", func (log sugar.Log) bool {
// 		log("and log this sentence")
// 		log(sugar.NewLogger().Log("and log this nested sentence"))
// 		return false
// 	}).
//
// 	Assert("this must never execute", func (log sugar.Log) bool {
// 		log("and not log this sentence")
// 		log(sugar.NewLogger().Log("and not log this nested sentence"))
// 		return true
// 	})
//
// 	sugar.New(t).
//
// 	Asset("the output from " func () bool {
//
// 	})
//
// }
//
// expectedOutput := `
// PASS [] this must pass
//  ┃ and log this sentence
//  ┖ ┖ and log this nested sentence
// PASS [] this passes and continues
//  ┃ and log this sentence
//  ┖ ┖ and log this nested sentence
// FAIL [] this fails but continues
//  ┃ and log this sentence
//  ┖ ┖ and log this nested sentence
// ======== this must add a title ========
// PASS [] this passes and continues
//  ┃ and log this sentence
//  ┖ ┖ and log this nested sentence
// WARN [] this passes and continues but warns
//  ┃ and log this sentence
//  ┖ ┖ and log this nested sentence
// `