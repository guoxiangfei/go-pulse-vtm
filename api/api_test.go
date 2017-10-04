package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestConnect(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}
	log.Printf("Supported versions:\n%+v\n", client.VersionsSupported)
}

func TestTraverseAllConfigurationResources(t *testing.T) {
	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	// work with an environment
	client.WorkWithConfigurationResources()

	log.Println("Root Path: ", client.RootPath)

	// Traversing the tree!!
	resources := make(map[string]interface{})
	err = client.TraverseTree(client.RootPath, resources)
	if err != nil {
		t.Fatal("Error traversing tree: ", err)
	}
	for url := range resources {
		log.Println("Found Resource URL: ", url)
	}
}

func TestGetAllResourceTypes(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	// work with an environment
	client.WorkWithConfigurationResources()

	res, err := client.GetAllResourceTypes()
	if err != nil {
		t.Fatal("Error getting all resource types:", err)
	}
	log.Println("Found resources:\n", res)
}

func TestGet(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	// work with an environment
	client.WorkWithConfigurationResources()

	//Get all resource types
	_, err = client.GetAllResourceTypes()

	objs, err := client.GetAllResources("virtual_servers")
	if err != nil {
		t.Fatal("Error getting all configuration resources: ", err)
	}

	for _, obj := range objs {
		//Get a resource by name
		objByName := make(map[string]interface{})
		err := client.GetByName("virtual_servers", obj["name"].(string), &objByName)
		if err != nil {
			t.Fatal("Error getting object by name: ", obj["name"].(string))
		}

		log.Printf("Retrieved resource:\n%+v\n", objByName)

		// ...or get it by URL
		objByURL := make(map[string]interface{})
		err = client.GetByURL(obj["href"].(string), &objByURL)
		if err != nil {
			t.Fatal("Error getting object by URL: ", obj["href"].(string))
		}
		log.Printf("Retrieved resource:\n%+v\n", objByURL)
	}
}

func TestSetAndDelete(t *testing.T) {
	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	profile := make(map[string]interface{})
	name := "test_vs" + strconv.Itoa(rand.Int())
	template := getJSONProfile()
	err = json.Unmarshal(template, &profile)

	log.Println("Going to create virtual server: ", name)
	newRes := make(map[string]interface{})
	err = client.Set("virtual_servers", name, profile, &newRes)
	if err != nil {
		t.Fatal("Error creating a resource:", err)
	}
	properties := newRes["properties"].(map[string]interface{})
	basic := properties["basic"].(map[string]interface{})
	assert.Equal(t, true, basic["add_cluster_ip"])

	//update the same resource
	log.Println("Going to update virtual server: ", name)
	template = getJSONUpdatedProfile()
	err = json.Unmarshal(template, &profile)
	updatedRes := make(map[string]interface{})
	err = client.Set("virtual_servers", name, profile, &updatedRes)
	if err != nil {
		t.Fatal("Error updating a resource", err)
	}
	properties = updatedRes["properties"].(map[string]interface{})
	basic = properties["basic"].(map[string]interface{})
	assert.Equal(t, false, basic["add_cluster_ip"])

	err = client.Delete("virtual_servers", name)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Resource %s deleted", name)
	}
}

func TestTraverseStatus(t *testing.T) {
	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	client.WorkWithStatus()
	resources := make(map[string]interface{})
	err = client.TraverseTree(client.RootPath, resources)
	if err != nil {
		t.Fatal("Error traversing tree: ", err)
	}
	log.Println("Traversed status: ", resources)
	for url := range resources {
		log.Println("Found Resource URL: ", url)
	}
}

func TestGetStatistics(t *testing.T) {

	server := getServer()
	if server == "" {
		t.Fatal("Error getting a valid server")
	}

	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	stats, err := client.GetStatistics(server)
	if err != nil {
		t.Fatal("Error getting statistics")
	}
	for key := range stats {
		log.Println("Found stat: ", key, stats[key])
	}
}

func TestGetInformation(t *testing.T) {

	server := getServer()
	if server == "" {
		t.Fatal("Error getting a valid server")
	}

	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}

	info, err := client.GetInformation(server)
	if err != nil {
		t.Fatal("Error getting information")
	}
	for key := range info {
		log.Println("Found stat: ", info[key])
	}
}

func getServer() string {

	// get a client
	client, err := GetClient()
	if err != nil {
		return ""
	}

	// traverse the status section...
	client.WorkWithStatus()
	resources := make(map[string]interface{})
	err = client.TraverseTree(client.RootPath, resources)
	if err != nil {
		return ""
	}
	for url := range resources {
		tokens := strings.Split(url, "/")
		server := tokens[5]
		log.Println("Server: ", server)
		return server
	}
	return ""
}

func TestGetState(t *testing.T) {
	server := getServer()
	if server == "" {
		t.Fatal("Error getting a valid server")
	}

	// get a client
	client, err := GetClient()
	if err != nil {
		t.Fatal("Error getting a client:", err)
	}
	state, err := client.GetState(server)
	if err != nil {
		t.Fatal("Error getting information")
	}
	log.Println("Node Status:\n", state)

}

func getJSONProfile() []byte {
	return []byte(`{"properties":{"basic":{"add_cluster_ip":true,"add_x_forwarded_for":false,"add_x_forwarded_proto":false,"autodetect_upgrade_headers":true,"bandwidth_class":"","close_with_rst":false,"completionrules":[],"connect_timeout":10,"enabled":false,"ftp_force_server_secure":true,"glb_services":[],"listen_on_any":true,"listen_on_hosts":[],"listen_on_traffic_ips":[],"note":"","pool":"pool_test_rui","port":90,"protection_class":"","protocol":"http","request_rules":[],"response_rules":[],"slm_class":"","so_nagle":false,"ssl_client_cert_headers":"none","ssl_decrypt":false,"ssl_honor_fallback_scsv":"use_default","transparent":false},"aptimizer":{"enabled":false,"profile":[]},"connection":{"keepalive":true,"keepalive_timeout":10,"max_client_buffer":65536,"max_server_buffer":65536,"max_transaction_duration":0,"server_first_banner":"","timeout":300},"connection_errors":{"error_file":"Default"},"cookie":{"domain":"no_rewrite","new_domain":"","path_regex":"","path_replace":"","secure":"no_modify"},"dns":{"edns_client_subnet":true,"edns_udpsize":4096,"max_udpsize":4096,"rrset_order":"fixed","verbose":false,"zones":[]},"ftp":{"data_source_port":0,"force_client_secure":true,"port_range_high":0,"port_range_low":0,"ssl_data":true},"gzip":{"compress_level":1,"enabled":false,"etag_rewrite":"wrap","include_mime":["text/html","text/plain"],"max_size":10000000,"min_size":1000,"no_size":true},"http":{"chunk_overhead_forwarding":"lazy","location_regex":"","location_replace":"","location_rewrite":"if_host_matches","mime_default":"text/plain","mime_detect":false},"http2":{"connect_timeout":0,"data_frame_size":4096,"enabled":true,"header_table_size":4096,"headers_index_blacklist":[],"headers_index_default":true,"headers_index_whitelist":[],"idle_timeout_no_streams":120,"idle_timeout_open_streams":600,"max_concurrent_streams":200,"max_frame_size":16384,"max_header_padding":0,"merge_cookie_headers":true,"stream_window_size":65535},"kerberos_protocol_transition":{"enabled":false,"principal":"","target":""},"log":{"client_connection_failures":false,"enabled":false,"filename":"%zeushome%/zxtm/log/%v.log","format":"%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"","save_all":true,"server_connection_failures":false,"session_persistence_verbose":false,"ssl_failures":false},"recent_connections":{"enabled":true,"save_all":false},"request_tracing":{"enabled":false,"trace_io":false},"rtsp":{"streaming_port_range_high":0,"streaming_port_range_low":0,"streaming_timeout":30},"sip":{"dangerous_requests":"node","follow_route":true,"max_connection_mem":65536,"mode":"sip_gateway","rewrite_uri":false,"streaming_port_range_high":0,"streaming_port_range_low":0,"streaming_timeout":60,"timeout_messages":true,"transaction_timeout":30},"smtp":{"expect_starttls":true},"ssl":{"add_http_headers":false,"client_cert_cas":[],"elliptic_curves":[],"issued_certs_never_expire":[],"ocsp_enable":false,"ocsp_issuers":[],"ocsp_max_response_age":0,"ocsp_stapling":false,"ocsp_time_tolerance":50,"ocsp_timeout":50,"prefer_sslv3":false,"request_client_cert":"dont_request","send_close_alerts":true,"server_cert_alt_certificates":[],"server_cert_default":"","server_cert_host_mapping":[],"signature_algorithms":"","ssl_ciphers":"","ssl_support_ssl2":"use_default","ssl_support_ssl3":"use_default","ssl_support_tls1":"use_default","ssl_support_tls1_1":"use_default","ssl_support_tls1_2":"use_default","trust_magic":false},"syslog":{"enabled":false,"format":"%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"","ip_end_point":"","msg_len_limit":1024},"tcp":{"proxy_close":false},"udp":{"end_point_persistence":true,"port_smp":false,"response_datagrams_expected":1,"timeout":7},"web_cache":{"control_out":"","enabled":false,"error_page_time":30,"max_time":600,"refresh_time":2}}}`)
}

func getJSONUpdatedProfile() []byte {
	return []byte(`{"properties":{"basic":{"add_cluster_ip":false,"add_x_forwarded_for":false,"add_x_forwarded_proto":false,"autodetect_upgrade_headers":true,"bandwidth_class":"","close_with_rst":false,"completionrules":[],"connect_timeout":10,"enabled":false,"ftp_force_server_secure":true,"glb_services":[],"listen_on_any":true,"listen_on_hosts":[],"listen_on_traffic_ips":[],"note":"","pool":"pool_test_rui","port":90,"protection_class":"","protocol":"http","request_rules":[],"response_rules":[],"slm_class":"","so_nagle":false,"ssl_client_cert_headers":"none","ssl_decrypt":false,"ssl_honor_fallback_scsv":"use_default","transparent":false},"aptimizer":{"enabled":false,"profile":[]},"connection":{"keepalive":true,"keepalive_timeout":10,"max_client_buffer":65536,"max_server_buffer":65536,"max_transaction_duration":0,"server_first_banner":"","timeout":300},"connection_errors":{"error_file":"Default"},"cookie":{"domain":"no_rewrite","new_domain":"","path_regex":"","path_replace":"","secure":"no_modify"},"dns":{"edns_client_subnet":true,"edns_udpsize":4096,"max_udpsize":4096,"rrset_order":"fixed","verbose":false,"zones":[]},"ftp":{"data_source_port":0,"force_client_secure":true,"port_range_high":0,"port_range_low":0,"ssl_data":true},"gzip":{"compress_level":1,"enabled":false,"etag_rewrite":"wrap","include_mime":["text/html","text/plain"],"max_size":10000000,"min_size":1000,"no_size":true},"http":{"chunk_overhead_forwarding":"lazy","location_regex":"","location_replace":"","location_rewrite":"if_host_matches","mime_default":"text/plain","mime_detect":false},"http2":{"connect_timeout":0,"data_frame_size":4096,"enabled":true,"header_table_size":4096,"headers_index_blacklist":[],"headers_index_default":true,"headers_index_whitelist":[],"idle_timeout_no_streams":120,"idle_timeout_open_streams":600,"max_concurrent_streams":200,"max_frame_size":16384,"max_header_padding":0,"merge_cookie_headers":true,"stream_window_size":65535},"kerberos_protocol_transition":{"enabled":false,"principal":"","target":""},"log":{"client_connection_failures":false,"enabled":false,"filename":"%zeushome%/zxtm/log/%v.log","format":"%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"","save_all":true,"server_connection_failures":false,"session_persistence_verbose":false,"ssl_failures":false},"recent_connections":{"enabled":true,"save_all":false},"request_tracing":{"enabled":false,"trace_io":false},"rtsp":{"streaming_port_range_high":0,"streaming_port_range_low":0,"streaming_timeout":30},"sip":{"dangerous_requests":"node","follow_route":true,"max_connection_mem":65536,"mode":"sip_gateway","rewrite_uri":false,"streaming_port_range_high":0,"streaming_port_range_low":0,"streaming_timeout":60,"timeout_messages":true,"transaction_timeout":30},"smtp":{"expect_starttls":true},"ssl":{"add_http_headers":false,"client_cert_cas":[],"elliptic_curves":[],"issued_certs_never_expire":[],"ocsp_enable":false,"ocsp_issuers":[],"ocsp_max_response_age":0,"ocsp_stapling":false,"ocsp_time_tolerance":50,"ocsp_timeout":50,"prefer_sslv3":false,"request_client_cert":"dont_request","send_close_alerts":true,"server_cert_alt_certificates":[],"server_cert_default":"","server_cert_host_mapping":[],"signature_algorithms":"","ssl_ciphers":"","ssl_support_ssl2":"use_default","ssl_support_ssl3":"use_default","ssl_support_tls1":"use_default","ssl_support_tls1_1":"use_default","ssl_support_tls1_2":"use_default","trust_magic":false},"syslog":{"enabled":false,"format":"%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"","ip_end_point":"","msg_len_limit":1024},"tcp":{"proxy_close":false},"udp":{"end_point_persistence":true,"port_smp":false,"response_datagrams_expected":1,"timeout":7},"web_cache":{"control_out":"","enabled":false,"error_page_time":30,"max_time":600,"refresh_time":2}}}`)
}