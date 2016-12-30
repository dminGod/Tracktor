package main

import (
	"./config"
	"time"
	"fmt"
	"./requests"
	"os"
	"strconv"
	"flag"
	"net/rpc"
	"net"
	"net/http"
)

// Response of the first Context call
type Context struct {

	Context int `json:"context"`
}

type Test struct {

	test_tag string
	file_name string
}

type CurrentTests struct {

	Tests []Test
}

func (t *CurrentTests) Removetest( test string, tp *string ) error  {

	var out CurrentTests

	// If the size of the count is 0 we dont want to try to remove this
	countExists := len( curTests.find_test(test).Tests )

	if countExists == 0 {  fmt.Println("Tried to remove a test", test, "that does not exist");  return nil }

	// Remove the tests range
	for _, item := range t.Tests {

		if item.test_tag != test {

			out.Tests = append(out.Tests, item)
		} else {

			close_custom_files( item.file_name )
		}
	}

	curTests = &out
	return nil
}

func (t *CurrentTests) Addtest( test string,  tp *string ) error {

	// We dont want to keep adding stuff of the same tag, its going to confuse stuff
	countExists := len( curTests.find_test(test).Tests )
	if countExists > 0 {  fmt.Println(curTests.find_test(test).Tests);  return nil }

	test_keyword := test + "_" + strconv.FormatInt( time.Now().Unix(), 10 )
	make_test_files( test_keyword )

	curTests.Tests = append(curTests.Tests, Test{ test, test_keyword })

	return nil
}

// Local method to find tests
func (t *CurrentTests) find_test ( test string )  ( out CurrentTests ) {

	for _, item := range t.Tests {
		if item.test_tag == test {
			out.Tests = append(out.Tests, item)
		}
	}

	return out
}

// The tests that are currently on..
var curTests = new(CurrentTests)

// Stop the scanning
var stop_hosts_c = make( chan bool )

// Signal to close the deamon, never... for now..
var done = make( chan bool )

// Get config items
var LogDirectory string = config.Get_config("conf", "log_directory")

var serverAddress  string = "127.0.0.1"
var hosts []string = config.Get_hosts()

// Holder for files
var ff map[string]*os.File

// For now we are dropping the CLI functionality.
// var is_cli *bool = flag.Bool("launch-cli", false, "If you want to start a CLI")
var is_daemon *bool = flag.Bool("start_server", false, "Start the server, you need to run this only once.")

var start_test *bool = flag.Bool("start_test", false, "To start a test, pass start_test use with the --test_name param")
var stop_test *bool = flag.Bool("stop_test", false, "This is to stop the test, use with the --test_name param")

var test_name *string = flag.String("test_name", "", "This is the test name usage, test_name=\"MY_TEST_NAME\"")

var stop_scanning *bool = flag.Bool("stop_scanning", false, "This is to stop the scanning of all the tests")


func init() {

	flag.Parse()

	// Make the hash to hold files
	ff = make( map[string]*os.File )

	if len(hosts) == 0 {

		fmt.Println("There are no hosts to monitor. Please configure some hosts to monitor.")
		os.Exit(1)
	}

	if *is_daemon && (*start_test || *stop_test) {
		fmt.Println("Server and test Start, Stop can not be in the same instance"); os.Exit(1) }

	if *start_test && *stop_test {
		fmt.Println("Start and Stop test can not be together.."); os.Exit(1)}

}


/*

A Deamon must start -- we will run this as & to start with that is good enough
It will have a list of hosts that it monitors -- that we are doing currently
 If any tests are running then it will send the web calls for monitoring
	 Make a struct for a test
	 Each test need not have its own go func -- Just maintaining the struct is good enough

	Do a count and register new tests as they come
	// In the ticker we have to modify
	After the successful call is made then it will write to all the files it needs to write to

todo: Make the initial calls that need to be made so vector can get the data
todo: Change Vector to accept the files that are given to it in a dropdown.. show multiple hosts at the same time..?

todo: Do you want to do graphite with PCP?  problem with vector is that it shows static data, doesnt have

 */

func main() {


	if *is_daemon {

//		defer stop_test()
		// Start the service as a deamon, let it remain on indefinitely..

		// Close the files once the main gets over
		defer close_files()

		start_hosts()

		// Start RPC only if started as daemon
		rpc.Register( curTests )
		rpc.HandleHTTP()

		fmt.Println("Starting daemon")
		l, _ := net.Listen("tcp", ":1234")

		go http.Serve(l, nil)

		go func(){
			for {
				time.Sleep(3 * time.Second)
				fmt.Println(  *curTests  )
			}
		}()


		<-done
	}

	// Manage tests.. both are same only one method is different, put them together
	if *start_test || *stop_test {

		manage_tests()
	}


	if *stop_scanning {

		stop_hosts()
	}


	// go func(){ for i := 0; i < 300; i++ {time.Sleep(time.Second); fmt.Println(".") } }()  // Dummy ToDO: Remove this...


}


// This handles the start and stop requests
func manage_tests() {

	// Make an RPC call
	client, _ := rpc.DialHTTP("tcp", serverAddress + ":1234")

	//		getTests := new(CurrentTests)
	var resp string = ""

	test_to_call := ""

	if *start_test {  test_to_call =  "CurrentTests.Addtest"  }
	if *stop_test {  test_to_call =  "CurrentTests.Removetest" }

	//var reply int
	err  := client.Call(test_to_call, test_name, &resp)

	if err != nil {

		fmt.Println("Issue with making the RPC call..")
		os.Exit(1)
	} else {

		fmt.Println(resp)
		os.Exit(0)
	}
}


func start_hosts()  {

	// Loop over all the hosts
	for  ind, host := range hosts {

		_ = ind // kachara

		context_variable := requests.Get_context( host )

		// For each host make a Go Func
		go func( host string, context_variable string ) {

			ticker := time.NewTicker( 2 * time.Second )

			for {
				select {
					// Timer
					case <-ticker.C:
						 if len(curTests.Tests) > 0 {

							 // Get the result for this host

							// Write to as many files as needed by the tests..
							for _, item := range curTests.Tests {

								requests.Log_Write(&host, &context_variable, ff[host + "_" + item.file_name])
							}
						 } else {

							 fmt.Println("No tests right now...")
						 }

					case <-stop_hosts_c:
						fmt.Println("Stopping hosts")
						return
				}
			}
		}( host, context_variable )
	}
}

func stop_hosts() {

	for  range hosts {

		stop_hosts_c <- true
	}

	fmt.Println("Stopped all hosts")
}

func make_test_files(test_keyword string)  {

	for _, host := range hosts {

		ff[host + "_" + test_keyword], _ = os.OpenFile(LogDirectory + host + "_" + test_keyword + ".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	}



}

func close_custom_files( keyword string) {

	fmt.Println("Closing custom files")

	for _, host := range hosts {
		ff[host + "_" + keyword].Close()
	}
}

func close_files() {

	fmt.Println("Closing files")
	for _, host := range hosts {

		ff[host].Close()
	}
}