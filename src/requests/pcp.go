package requests

import (
	"fmt"
	"strings"
	"os"
	"../config"
	"strconv"
	"time"
)

var port string = config.Get_config("conf", "port")

var connected_on_lan bool = false

func Get_context( host string) string {

	inter_s := ""

	if connected_on_lan {
		context_url := "http://" + host + ":44323/pmapi/context?hostspec=localhost&polltimeout=600"

		resp := Make_http_call( context_url )

		fmt.Println( string(resp) )

		inter := strings.Split( string(resp), ":" )
		inter_s = strings.Replace(inter[1], "}", "", 1)
		inter_s = strings.Replace(inter_s, " ", "", -1)
	} else {

		inter_s = "TheToken"
	}

	return inter_s
}


// Todo : This is not going to be like this, this is log write...
// Log_write
func Log_Write( host *string, context *string, f *os.File ) {

	info := ""


	if connected_on_lan {

		info = string(Make_http_call("http://" + *host + ":" + port + "/pmapi/" + *context + "/_fetch?names=containers.cgroup,containers.name,kernel.all.cpu.sys,kernel.all.cpu.user,hinv.ncpu,kernel.percpu.cpu.sys,kernel.percpu.cpu.user,kernel.all.runnable,kernel.all.load,network.interface.in.bytes,network.interface.out.bytes,network.tcpconn.established,network.tcpconn.time_wait,network.tcpconn.close_wait,network.interface.in.packets,network.interface.out.packets,network.tcp.retranssegs,network.tcp.timeouts,network.tcp.listendrops,network.tcp.fastretrans,network.tcp.slowstartretrans,network.tcp.syncretrans,mem.util.cached,mem.util.used,mem.util.free,mem.util.bufmem,mem.vmstat.pgfault,mem.vmstat.pgmajfault,kernel.all.pswitch,disk.dev.read,disk.dev.write,disk.dev.read_bytes,disk.dev.write_bytes,disk.dev.avactive,disk.dev.read_rawactive,disk.dev.write_rawactive"))
		info = strings.Replace(info, "\n", "", -1)
		info = "\n'{\"data\":" + info + ",\"status\":200,\"config\":{\"method\":\"GET\",\"transformRequest\":[null],\"transformResponse\":[null],\"url\":\"http://192.168.2.129:44323/pmapi/91713890/_fetch\",\"params\":{\"names\":\"containers.cgroup,containers.name,kernel.all.cpu.sys,kernel.all.cpu.user,hinv.ncpu,kernel.percpu.cpu.sys,kernel.percpu.cpu.user,kernel.all.runnable,kernel.all.load,network.interface.in.bytes,network.interface.out.bytes,network.tcpconn.established,network.tcpconn.time_wait,network.tcpconn.close_wait,network.interface.in.packets,network.interface.out.packets,network.tcp.retranssegs,network.tcp.timeouts,network.tcp.listendrops,network.tcp.fastretrans,network.tcp.slowstartretrans,network.tcp.synretrans,mem.util.cached,mem.util.used,mem.util.free,mem.util.bufmem,mem.vmstat.pgfault,mem.vmstat.pgmajfault,kernel.all.pswitch,disk.dev.read,disk.dev.write,disk.dev.read_bytes,disk.dev.write_bytes,disk.dev.avactive,disk.dev.read_rawactive,disk.dev.write_rawactive\"},\"headers\":{\"Accept\":\"application/json, text/plain, */*\"}},\"statusText\":\"OK\"}',"
	} else {

		info = "Test line" + strconv.FormatInt( time.Now().Unix(), 10 )
	}



	if _, err := f.WriteString( info ); err != nil {

		fmt.Println(err)
		fmt.Println(f)
	}
}