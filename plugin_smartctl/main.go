package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Satainfo struct {
	Endpoint    string                `json:"endpoint"`
	Id          string                `json:"id"`
	Time        int64                 `json:"time"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	DeviceModel string                `json:"devicemodel"`
	Data        map[string][4]float64 `json:"data"`
}
type Sasinfo struct {
	Endpoint      string `json:"endpoint"`
	Id            string `json:"id"`
	Time          int64  `json:"time"`
	Type          string `json:"type"`
	Version       string `json:"version"`
	VendorProduct string `json:"devicemodel"`
	Data          struct {
		//Sasdata  map[string]string        `json:"sasdata"`
		Current_Drive_Temperature                                  string `json:"Current Drive Temperature"`
		Drive_Trip_Temperature                                     string `json:"Drive Trip Temperature"`
		Manufactured_in_week                                       string `json:"Manufactured in week"`
		Specified_cycle_count_over_device_lifetime                 string `json:"Specified cycle count over device lifetime"`
		Accumulated_start_stop_cycles                              string `json:"Accumulated start-stop cycles"`
		Specified_load_unload_count_over_device_lifetime           string `json:"Specified load-unload count over device lifetime"`
		Accumulated_load_unload_cycles                             string `json:"Accumulated load-unload cycles"`
		Elements_in_grown_defect_list                              string `json:"Elements in grown defect list"`
		Blocks_sent_to_initiator                                   string `json:"Blocks sent to initiator"`
		Blocks_received_from_initiator                             string `json:"Blocks received from initiator"`
		Blocks_read_from_cache_and_sent_to_initiator               string `json:"Blocks read from cache and sent to initiator"`
		Number_of_read_and_write_commands_whose_size__segment_size string `json:"Number of read and write commands whose size <= segment size"`
		Number_of_read_and_write_commands_whose_size_segment_size  string `json:"Number of read and write commands whose size > segment size"`
		number_of_hours_powered_up                                 string `json:"number of hours powered up"`
		number_of_minutes_until_next_internal_SMART_test           string `json:"number of minutes until next internal SMART test"`
		Errors_Corrected_by_ECC_fast_read                          string `json:"Errors Corrected by ECC fast read"`
		Errors_Corrected_by_ECC_delayed_read                       string `json:"Errors Corrected by ECC delayed read"`
		Errors_Corrected_by_rereads_rewrites_read                  string `json:"Errors Corrected by rereads/rewrites read"`
		Total_errors_corrected_read                                string `json:"Total errors corrected read"`
		Correction_algorithm_invocations_read                      string `json:"Correction algorithm invocations read"`
		Gigabytes_processed_read                                   string `json:"Gigabytes processed [10^9 bytes] read"`
		Total_uncorrected_errors_read                              string `json:"Total uncorrected errors read"`
		Errors_Corrected_by_ECC_fast_write                         string `json:"Errors Corrected by ECC fast write"`
		Errors_Corrected_by_ECC_delayed_write                      string `json:"Errors Corrected by ECC delayed write"`
		Errors_Corrected_by_rereads_rewrites_write                 string `json:"Errors Corrected by rewrites/rewrites write"`
		Total_errors_corrected_write                               string `json:"Total errors corrected write"`
		Correction_algorithm_invocations_write                     string `json:"Correction algorithm invocations write"`
		Gigabytes_processed_write                                  string `json:"Gigabytes processed [10^9 bytes] write"`
		Total_uncorrected_errors_write                             string `json:"Total uncorrected errors write"`
	} `json:"data"`
}

var plu_name = "plugin_smartctl"
var id, value, thresh, raw_value, endpoint string
var vendor, product string
var sasinfo Sasinfo
var satainfo Satainfo

func main() {
	Getsmartinfo()
}
func Ifraid() string {
	p := bytes.NewBuffer(nil)
	command := exec.Command("/bin/sh", "-c", "lspci | grep RAID")
	command.Stdout = p
	command.Run()
	f := string(p.Bytes())
	r := strings.TrimSpace(f)
	n := len(r)
	if n != 0 {
		//return "sudo /home/v-wxbroot/agent_new/smartctl --scan |grep megaraid"
		return "smartctl --scan |grep megaraid"
	} else {
		//return "sudo /home/v-wxbroot/agent_new/smartctl --scan |grep -v megaraid"
		return "smartctl --scan |grep -v megaraid"
	}
}
func Getsmartinfo() {
	endpoint = Getip()
	sasinfo.Endpoint = endpoint
	satainfo.Endpoint = endpoint
	sasinfo.Version = "1.0"
	satainfo.Version = "1.0"
	command := Ifraid()
	cmd := exec.Command("/bin/sh", "-c", command)
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error=>", err.Error())
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		v := strings.Fields(line)
		var smartcmd, tags string
		if len(v) == 7 {
			//smartcmd = "sudo /home/v-wxbroot/agent_new/smartctl -a " + v[0]
			smartcmd = "smartctl -a " + v[0]
			tags = v[4]
			sasinfo.Id = tags
			satainfo.Id = tags
		}
		if len(v) == 8 {
			//smartcmd = "sudo /home/v-wxbroot/agent_new/smartctl -a -d " + v[2] + " " + v[0]
			smartcmd = "smartctl -a -d " + v[2] + " " + v[0]
			tags = v[5]
			sasinfo.Id = tags
			satainfo.Id = tags
		}
		//sasinfo.Data.Sasdata = make(map[string]string)
		satainfo.Data = make(map[string][4]float64)
		command := exec.Command("/bin/sh", "-c", smartcmd)
		fmt.Println(command.Args)
		stdout, err := command.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error=>", err.Error())
		}
		command.Start()
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			t := time.Now().Unix()
			//timestamp := fmt.Sprintf("%d", t)
			//b, error := strconv.Atoi(t)
			sasinfo.Time = t
			if strings.Contains(line, "SAS") {
				//slice1 := strings.Fields(line)
				//v1 := slice1[2]
				sasinfo.Type = "SAS"
			}
			if strings.Contains(line, "SATA") {
				//slice1 := strings.Fields(line)
				//v1 := slice1[2]
				satainfo.Type = "SATA"
			}
			if strings.Contains(line, "Device Model") {
				slice1 := strings.Fields(line)
				v1 := slice1[2]
				satainfo.DeviceModel = v1
			}
			if strings.Contains(line, "Vendor:") {
				slice1 := strings.Fields(line)
				v1 := slice1[1]
				vendor = v1
			}
			if strings.Contains(line, "Product:") {
				slice1 := strings.Fields(line)
				v1 := slice1[1]
				product = v1
			}
			sasinfo.VendorProduct = vendor + product
			if strings.Contains(line, "Current Drive Temperature") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//fmt.Println(v1, timestamp, tag, endpoint)
				//sasinfo.Data.Sasdata["Current Drive Temperature"] = f
				sasinfo.Data.Current_Drive_Temperature = v1
				//pushIt(v1, timestamp, "Current Drive Temperature", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Drive Trip Temperature") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				sasinfo.Data.Drive_Trip_Temperature = v1
				//pushIt(v1, timestamp, "Drive Trip Temperature", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Manufactured in week") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				v2 := slice1[6]
				date := v2 + v1
				//f, _ := strconv.ParseFloat(date, 64)
				//tag := "disk = " + tags
				fmt.Println(date)
				sasinfo.Data.Manufactured_in_week = date
				//pushIt(date, timestamp, "Manufacture date", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Specified cycle count over device lifetime") {
				slice1 := strings.Fields(line)
				v1 := slice1[6]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				sasinfo.Data.Specified_cycle_count_over_device_lifetime = v1
				//pushIt(v1, timestamp, "Specified cycle count over device lifetime", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Accumulated start-stop cycles") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Accumulated start-stop cycles"] = f
				sasinfo.Data.Accumulated_start_stop_cycles = v1
				//pushIt(v1, timestamp, "Accumulated start-stop cycles", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Specified load-unload count over device lifetime") {
				slice1 := strings.Fields(line)
				v1 := slice1[6]
				//f, _ := strconv.ParseFloat(v1, 64)
				//sasinfo.Data.Sasdata["Specified load-unload count over device lifetime"] = f
				sasinfo.Data.Specified_load_unload_count_over_device_lifetime = v1
				//tag := "disk = " + tags
				//pushIt(v1, timestamp, "Specified load-unload count over device lifetime", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Accumulated load-unload cycles") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Accumulated load-unload cycles"] = f
				sasinfo.Data.Accumulated_load_unload_cycles = v1
				//pushIt(v1, timestamp, "Accumulated load-unload cycles", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Elements in grown defect list") {
				slice1 := strings.Fields(line)
				v1 := slice1[5]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Elements in grown defect list"] = f
				sasinfo.Data.Elements_in_grown_defect_list = v1
				//pushIt(v1, timestamp, "Elements in grown defect list", tag, "", "GAUGE", endpoint)
			}
			/*if strings.Contains(line, "Non-medium error count") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				sasinfo.Data.Sasdata["Non-medium error count"] = f
				//pushIt(v1, timestamp, "Non-medium error count", tag, "", "GAUGE", endpoint)
			}*/
			if strings.Contains(line, "Blocks sent to initiator") {
				slice1 := strings.Fields(line)
				v1 := slice1[5]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Blocks sent to initiator"] = f
				sasinfo.Data.Blocks_sent_to_initiator = v1
				//pushIt(v1, timestamp, "Blocks sent to initiator", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Blocks received from initiator") {
				slice1 := strings.Fields(line)
				v1 := slice1[5]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Blocks received from initiator"] = f
				sasinfo.Data.Blocks_received_from_initiator = v1
				//pushIt(v1, timestamp, "Blocks received from initiator", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Blocks read from cache and sent to initiator") {
				slice1 := strings.Fields(line)
				v1 := slice1[9]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Blocks read from cache and sent to initiator"] = f
				sasinfo.Data.Blocks_read_from_cache_and_sent_to_initiator = v1
				//pushIt(v1, timestamp, "Blocks read from cache and sent to initiator", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Number of read and write commands whose size <= segment size") {
				slice1 := strings.Fields(line)
				v1 := slice1[12]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Number of read and write commands whose size <= segment size"] = f
				sasinfo.Data.Number_of_read_and_write_commands_whose_size__segment_size = v1
				//pushIt(v1, timestamp, "Number of read and write commands whose size <= segment size", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "Number of read and write commands whose size > segment size") {
				slice1 := strings.Fields(line)
				v1 := slice1[12]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["Number of read and write commands whose size > segment size"] = f
				sasinfo.Data.Number_of_read_and_write_commands_whose_size_segment_size = v1
				//pushIt(v1, timestamp, "Number of read and write commands whose size > segment size", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "number of hours powered up") {
				slice1 := strings.Fields(line)
				v1 := slice1[6]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["number of hours powered up"] = f
				sasinfo.Data.number_of_hours_powered_up = v1
				//pushIt(v1, timestamp, "number of hours powered up", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "number of minutes until next internal SMART test") {
				slice1 := strings.Fields(line)
				v1 := slice1[9]
				//f, _ := strconv.ParseFloat(v1, 64)
				//tag := "disk = " + tags
				//sasinfo.Data.Sasdata["number of minutes until next internal SMART test"] = f
				sasinfo.Data.number_of_minutes_until_next_internal_SMART_test = v1
				//pushIt(v1, timestamp, "number of minutes until next internal SMART test", tag, "", "GAUGE", endpoint)
			}

			if strings.Contains(line, "read:") {
				slice1 := strings.Fields(line)
				//tag := "disk = " + tagsf, _ := strconv.ParseFloat(v1, 64)
				/*f1, _ := strconv.ParseFloat(slice1[1], 64)
				f2, _ := strconv.ParseFloat(slice1[2], 64)
				f3, _ := strconv.ParseFloat(slice1[3], 64)
				f4, _ := strconv.ParseFloat(slice1[4], 64)
				f5, _ := strconv.ParseFloat(slice1[5], 64)
				f6, _ := strconv.ParseFloat(slice1[6], 64)
				f7, _ := strconv.ParseFloat(slice1[7], 64)*/
				sasinfo.Data.Errors_Corrected_by_ECC_fast_read = slice1[1]
				sasinfo.Data.Errors_Corrected_by_ECC_delayed_read = slice1[2]
				sasinfo.Data.Errors_Corrected_by_rereads_rewrites_read = slice1[3]
				sasinfo.Data.Total_errors_corrected_read = slice1[4]
				sasinfo.Data.Correction_algorithm_invocations_read = slice1[5]
				sasinfo.Data.Gigabytes_processed_read = slice1[6]
				sasinfo.Data.Total_uncorrected_errors_read = slice1[7]
				/*sasinfo.Data.Sasdata["Errors Corrected by ECC fast read"] = f1
				sasinfo.Data.Sasdata["Errors Corrected by ECC delayed read"] = f2
				sasinfo.Data.Sasdata["Errors Corrected by rereads/rewrites read"] = f3
				sasinfo.Data.Sasdata["Total errors corrected read"] = f4
				sasinfo.Data.Sasdata["Correction algorithm invocations read"] = f5
				sasinfo.Data.Sasdata["Gigabytes processed [10^9 bytes] read"] = f6
				sasinfo.Data.Sasdata["Total uncorrected errors read"] = f7

					pushIt(slice1[1], timestamp, "Errors Corrected by ECC fast read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[2], timestamp, "Errors Corrected by ECC delayed read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[3], timestamp, "Errors Corrected by ECC rewrites read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[4], timestamp, "Errors Corrected by ECC corrected read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[5], timestamp, "Errors Corrected by ECC invocations read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[6], timestamp, "Errors Corrected by ECC [10^9 bytes] read", tag, "", "GAUGE", endpoint)
					pushIt(slice1[7], timestamp, "Errors Corrected by ECC errors read", tag, "", "GAUGE", endpoint)*/
			}
			if strings.Contains(line, "write:") {
				slice1 := strings.Fields(line)
				//tag := "disk = " + tags
				/*f1, _ := strconv.ParseFloat(slice1[1], 64)
				f2, _ := strconv.ParseFloat(slice1[2], 64)
				f3, _ := strconv.ParseFloat(slice1[3], 64)
				f4, _ := strconv.ParseFloat(slice1[4], 64)
				f5, _ := strconv.ParseFloat(slice1[5], 64)
				f6, _ := strconv.ParseFloat(slice1[6], 64)
				f7, _ := strconv.ParseFloat(slice1[7], 64)*/
				sasinfo.Data.Errors_Corrected_by_ECC_fast_write = slice1[1]
				sasinfo.Data.Errors_Corrected_by_ECC_delayed_write = slice1[2]
				sasinfo.Data.Errors_Corrected_by_rereads_rewrites_write = slice1[3]
				sasinfo.Data.Total_errors_corrected_write = slice1[4]
				sasinfo.Data.Correction_algorithm_invocations_write = slice1[5]
				sasinfo.Data.Gigabytes_processed_write = slice1[6]
				sasinfo.Data.Total_uncorrected_errors_write = slice1[7]
				/*sasinfo.Data.Sasdata["Errors Corrected by ECC fast write"] = f1
				sasinfo.Data.Sasdata["Errors Corrected by ECC delayed write"] = f2
				sasinfo.Data.Sasdata["Errors Corrected by ECC rewrites write"] = f3
				sasinfo.Data.Sasdata["Total errors corrected  write"] = f4
				sasinfo.Data.Sasdata["Correction algorithm invocations write"] = f5
				sasinfo.Data.Sasdata["Gigabytes processed [10^9 bytes] write"] = f6
				sasinfo.Data.Sasdata["Total uncorrected errors write"] = f7
				pushIt(slice1[1], timestamp, "Errors Corrected  by ECC fast write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[2], timestamp, "Errors Corrected  by ECC delayed write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[3], timestamp, "Errors Corrected  by ECC rewrites write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[4], timestamp, "Errors Corrected  by ECC corrected write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[5], timestamp, "Errors Corrected  by ECC invocations write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[6], timestamp, "Errors Corrected  by ECC [10^9 bytes] write", tag, "", "GAUGE", endpoint)
				pushIt(slice1[7], timestamp, "Errors Corrected  by ECC errors write", tag, "", "GAUGE", endpoint)*/
			}
			/*if strings.Contains(line, "Non-medium error count") {
				slice1 := strings.Fields(line)
				v1 := slice1[3]
				tag := "disk = " + tags
				pushIt(v1, timestamp, "Non-medium error count", tag, "", "GAUGE", endpoint)
			}
			if strings.Contains(line, "ATTRIBUTE_NAME") {
				v := strings.Fields(line)
				id = v[0]
				value = v[3]
				thresh = v[5]
				raw_value = v[9]
			}*/
			if strings.Contains(line, "0x00") && strings.Contains(line, "-") {
				val := strings.Fields(line)
				f1, _ := strconv.ParseFloat(val[0], 64)
				f2, _ := strconv.ParseFloat(val[3], 64)
				f3, _ := strconv.ParseFloat(val[5], 64)
				f4, _ := strconv.ParseFloat(val[9], 64)
				//LogRun(plu_name + "*****" + "smartkey: " + smartkey)
				//LogRun(plu_name + "*****" + "smartvalue : " + smartvalue)
				t := time.Now().Unix()
				//timestamp := fmt.Sprintf("%d", t)
				satainfo.Time = t
				//tag1 := "type = " + id + "," + " disk = " + tags
				//tag2 := "type = " + value + "," + " disk = " + tags
				//tag3 := "type = " + thresh + "," + " disk = " + tags
				//tag4 := "type = " + raw_value + "," + " disk = " + tags
				//pushIt(val[0], timestamp, val[1], tag1, "", "GAUGE", endpoint)
				//pushIt(val[3], timestamp, val[1], tag2, "", "GAUGE", endpoint)
				//pushIt(val[5], timestamp, val[1], tag3, "", "GAUGE", endpoint)
				//pushIt(val[9], timestamp, val[1], tag4, "", "GAUGE", endpoint)
				satainfo.Data[val[1]] = [4]float64{f1, f2, f3, f4}
			}
		}
		if sasinfo.Type != "" && sasinfo.Type == "SAS" {
			jsonsasinfo, _ := json.Marshal(sasinfo)
			//pushIt(value, timestamp, metric, tags, containerId, counterType, endpoint)
			fmt.Println(string(jsonsasinfo))
		}
		if satainfo.Type != "" && satainfo.Type == "SATA" {
			jsonsatainfo, _ := json.Marshal(satainfo)
			fmt.Println(string(jsonsatainfo))
		}
		//t := time.Now().Unix()
		//timestamp := fmt.Sprintf("%d", t)
		//pushIt(string(jsonsasinfo), timestamp, "sasinfo", "", "", "GAUGE", endpoint)
		command.Wait()
	}
	cmd.Wait()
}

func Getip() string {
	address, err := net.InterfaceByName("enp4s0f0")
	if err != nil {
		fmt.Println("failed to query ip")
		os.Exit(2)
	}
	ip_info, err := address.Addrs()
	ip := strings.Split(ip_info[0].String(), "/")
	fmt.Println(ip[0])
	return ip[0]
}
