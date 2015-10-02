#warning
This is a very early pass of a package with tons of potential. As a result it will remain unstable until I feel as though it satisfies all of my major use cases. But for now, enjoy it. If you have reccomendations for improvements, I'm all ears.

#sugar
Sugar is wrapper around the standard library's testing.T interface that makes tests more beautiful and syntactically clear. There are four main functions

* sugar.New:	makes you some sugar... awww yeaaaaa
* s.Assert:		flags a test as failed but continues execution of the test
* s.Warn:		alerts that something is wrong but the test will pass
* s.Must:		flags a test as failed and prevents subsaquent tests from running 

#example
This is an example of how you could use sugard for an integration test for an endpoint of simple restful api that serves widgets.

**Widget_test.go**

	import (
		"testing"
		"github.com/marksalpeter/sugar"
	)
	
	func TestMain(m *testing.M) {
		// setup your server and db connection in TestMain
		// ...
		sugar.New(nil).Must("database is initialized", func (_ sugar.Log) bool {
			return true
		}).Must("server is running", func (_ sugar.Log) bool {
			return true
		})
		
		exitCode := m.Run()
		
		// nuke everything in the teardown
		...
		
		os.Exit(exitCode)
	}
	
	// test for /widget endpoint
	func TestWidget(t *testing.T) {
		
		// populate the test database with alerts using an orm of your choice
		// ...
		sugar.New(t).Must("db was populated with widgets", func (_ sugar.Log) bool {
			return true
		}).Must("user can authenticate", func (_ sugar.Log) bool {
			return true
		}).Warn("user is vaguely threatended by this warning", func (log sugar.Log) bool {
			log("fml!!")
			return false
		}).Assert("authed user can retreive a collection of widgets", func (_ sugar.Log) bool {
			return false
		}).Assert("authed user cannot create a widget", func (_ sugar.Log) bool {
			return false
		}).Assert("authed user can retreive an individual alert", func (_ sugar.Log) bool {
			return false
		}).Assert("authed user can update an alert", func (_ sugar.Log) bool {
			return false
		}).Assert("authed user can delete an alert", func (_ sugar.Log) bool {
			return false
		}).Assert("unauthed user cannot retreive a collection of alerts", func (_ sugar.Log) bool {
			return false
		}).Assert("unauthed user cannot create an alert", func (_ sugar.Log) bool {
			return false
		}).Assert("unauthed user cannot retreive an individual alert", func (_ sugar.Log) bool {
			return false
		}).Assert("unauthed user cannot update an alert", func (_ sugar.Log) bool {
			return false
		}).Assert("unauthed user cannot delete an alert", func (_ sugar.Log) bool {
			return false
		})
	}	
