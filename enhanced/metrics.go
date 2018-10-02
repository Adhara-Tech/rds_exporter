package enhanced

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// osMetrics represents available Enhanced Monitoring OS metrics from CloudWatch Logs.
//
// See https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.html#USER_Monitoring.OS.CloudWatchLogs
//
// Confusingly, some of the present metrics (for example, "writeLatency") are documented at
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/rds-metricscollected.html
//
//nolint:lll
type osMetrics struct {
	Engine             string    `json:"engine"             help:"The database engine for the DB instance."`
	InstanceID         string    `json:"instanceID"         help:"The DB instance identifier."`
	InstanceResourceID string    `json:"instanceResourceID" help:"A region-unique, immutable identifier for the DB instance, also used as the log stream identifier."`
	NumVCPUs           int       `json:"numVCPUs"           help:"The number of virtual CPUs for the DB instance."`
	Timestamp          time.Time `json:"timestamp"          help:"The time at which the metrics were taken."`
	Uptime             string    `json:"uptime"             help:"The amount of time that the DB instance has been active."`
	Version            float64   `json:"version"            help:"The version of the OS metrics' stream JSON format."`

	CPUUtilization    cpuUtilization    `json:"cpuUtilization"`
	DiskIO            []diskIO          `json:"diskIO"`
	FileSys           []fileSys         `json:"fileSys"`
	LoadAverageMinute loadAverageMinute `json:"loadAverageMinute"`
	Memory            memory            `json:"memory"`
	Network           []network         `json:"network"`
	ProcessList       []processList     `json:"processList"`
	Swap              swap              `json:"swap"`
	Tasks             tasks             `json:"tasks"`
}

type cpuUtilization struct {
	Guest  float64 `json:"guest"  help:"The percentage of CPU in use by guest programs."`
	Idle   float64 `json:"idle"   help:"The percentage of CPU that is idle."`
	Irq    float64 `json:"irq"    help:"The percentage of CPU in use by software interrupts."`
	Nice   float64 `json:"nice"   help:"The percentage of CPU in use by programs running at lowest priority."`
	Steal  float64 `json:"steal"  help:"The percentage of CPU in use by other virtual machines."`
	System float64 `json:"system" help:"The percentage of CPU in use by the kernel."`
	Total  float64 `json:"total"  help:"The total percentage of the CPU in use. This value includes the nice value."`
	User   float64 `json:"user"   help:"The percentage of CPU in use by user programs."`
	Wait   float64 `json:"wait"   help:"The percentage of CPU unused while waiting for I/O access."`
}

//nolint:lll
type diskIO struct {
	// common
	ReadIOsPS  float64 `json:"readIOsPS"  help:"The number of read operations per second."`
	WriteIOsPS float64 `json:"writeIOsPS" help:"The number of write operations per second."`
	Device     string  `json:"device"     help:"The identifier of the disk device in use."`

	// non-Aurora
	AvgQueueLen *float64 `json:"avgQueueLen" help:"The number of requests waiting in the I/O device's queue."`
	AvgReqSz    *float64 `json:"avgReqSz"    help:"The average request size, in kilobytes."`
	Await       *float64 `json:"await"       help:"The number of milliseconds required to respond to requests, including queue time and service time."`
	ReadKb      *int     `json:"readKb"      help:"The total number of kilobytes read."`
	ReadKbPS    *float64 `json:"readKbPS"    help:"The number of kilobytes read per second."`
	RrqmPS      *float64 `json:"rrqmPS"      help:"The number of merged read requests queued per second."`
	TPS         *float64 `json:"tps"         help:"The number of I/O transactions per second."`
	Util        *float64 `json:"util"        help:"The percentage of CPU time during which requests were issued."`
	WriteKb     *int     `json:"writeKb"     help:"The total number of kilobytes written."`
	WriteKbPS   *float64 `json:"writeKbPS"   help:"The number of kilobytes written per second."`
	WrqmPS      *float64 `json:"wrqmPS"      help:"The number of merged write requests queued per second."`

	// Aurora
	DiskQueueDepth  *float64 `json:"diskQueueDepth"  help:"The number of outstanding IOs (read/write requests) waiting to access the disk."`
	ReadLatency     *float64 `json:"readLatency"     help:"The average amount of time taken per disk I/O operation."`
	ReadThroughput  *float64 `json:"readThroughput"  help:"The average number of bytes read from disk per second."`
	WriteLatency    *float64 `json:"writeLatency"    help:"The average amount of time taken per disk I/O operation."`
	WriteThroughput *float64 `json:"writeThroughput" help:"The average number of bytes written to disk per second."`
}

//nolint:lll
type fileSys struct {
	MaxFiles        int     `json:"maxFiles"        help:"The maximum number of files that can be created for the file system."`
	MountPoint      string  `json:"mountPoint"      help:"The path to the file system."`
	Name            string  `json:"name"            help:"The name of the file system."`
	Total           int     `json:"total"           help:"The total number of disk space available for the file system, in kilobytes."`
	Used            int     `json:"used"            help:"The amount of disk space used by files in the file system, in kilobytes."`
	UsedFilePercent float64 `json:"usedFilePercent" help:"The percentage of available files in use."`
	UsedFiles       int     `json:"usedFiles"       help:"The number of files in the file system."`
	UsedPercent     float64 `json:"usedPercent"     help:"The percentage of the file-system disk space in use."`
}

type loadAverageMinute struct {
	Fifteen float64 `json:"fifteen" help:"The number of processes requesting CPU time over the last 15 minutes."`
	Five    float64 `json:"five"    help:"The number of processes requesting CPU time over the last 5 minutes."`
	One     float64 `json:"one"     help:"The number of processes requesting CPU time over the last minute."`
}

//nolint:lll
type memory struct {
	Active         int `json:"active"         node:"Active"          m:"1024" help:"The amount of assigned memory, in kilobytes."`
	Buffers        int `json:"buffers"        node:"Buffers"         m:"1024" help:"The amount of memory used for buffering I/O requests prior to writing to the storage device, in kilobytes."`
	Cached         int `json:"cached"         node:"Cached"          m:"1024" help:"The amount of memory used for caching file system–based I/O."`
	Dirty          int `json:"dirty"          node:"Dirty"           m:"1024" help:"The amount of memory pages in RAM that have been modified but not written to their related data block in storage, in kilobytes."`
	Free           int `json:"free"           node:"MemFree"         m:"1024" help:"The amount of unassigned memory, in kilobytes."`
	HugePagesFree  int `json:"hugePagesFree"  node:"HugePages_Free"  m:"1"    help:"The number of free huge pages. Huge pages are a feature of the Linux kernel."`
	HugePagesRsvd  int `json:"hugePagesRsvd"  node:"HugePages_Rsvd"  m:"1"    help:"The number of committed huge pages."`
	HugePagesSize  int `json:"hugePagesSize"  node:"Hugepagesize"    m:"1024" help:"The size for each huge pages unit, in kilobytes."`
	HugePagesSurp  int `json:"hugePagesSurp"  node:"HugePages_Surp"  m:"1"    help:"The number of available surplus huge pages over the total."`
	HugePagesTotal int `json:"hugePagesTotal" node:"HugePages_Total" m:"1"    help:"The total number of huge pages for the system."`
	Inactive       int `json:"inactive"       node:"Inactive"        m:"1024" help:"The amount of least-frequently used memory pages, in kilobytes."`
	Mapped         int `json:"mapped"         node:"Mapped"          m:"1024" help:"The total amount of file-system contents that is memory mapped inside a process address space, in kilobytes."`
	PageTables     int `json:"pageTables"     node:"PageTables"      m:"1024" help:"The amount of memory used by page tables, in kilobytes."`
	Slab           int `json:"slab"           node:"Slab"            m:"1024" help:"The amount of reusable kernel data structures, in kilobytes."`
	Total          int `json:"total"          node:"MemTotal"        m:"1024" help:"The total amount of memory, in kilobytes."`
	Writeback      int `json:"writeback"      node:"Writeback"       m:"1024" help:"The amount of dirty pages in RAM that are still being written to the backing storage, in kilobytes."`
}

type network struct {
	Interface string  `json:"interface" help:"The identifier for the network interface being used for the DB instance."`
	Rx        float64 `json:"rx"        help:"The number of bytes received per second."`
	Tx        float64 `json:"tx"        help:"The number of bytes uploaded per second."`
}

//nolint:lll
type processList struct {
	CPUUsedPC    float64 `json:"cpuUsedPc"    help:"The percentage of CPU used by the process."`
	ID           int     `json:"id"           help:"The identifier of the process."`
	MemoryUsedPC float64 `json:"memoryUsedPc" help:"The amount of memory used by the process, in kilobytes."`
	Name         string  `json:"name"         help:"The name of the process."`
	ParentID     int     `json:"parentID"     help:"The process identifier for the parent process of the process."`
	RSS          int     `json:"rss"          help:"The amount of RAM allocated to the process, in kilobytes."`
	TGID         int     `json:"tgid"         help:"The thread group identifier, which is a number representing the process ID to which a thread belongs. This identifier is used to group threads from the same process."`
	VSS          int     `json:"vss"          help:"The amount of virtual memory allocated to the process, in kilobytes."`
}

type swap struct {
	Cached int `json:"cached" help:"The amount of swap memory, in kilobytes, used as cache memory."`
	Free   int `json:"free"   help:"The total amount of swap memory free, in kilobytes."`
	Total  int `json:"total"  help:"The total amount of swap memory available, in kilobytes."`
}

type tasks struct {
	Blocked  int `json:"blocked"  help:"The number of tasks that are blocked."`
	Running  int `json:"running"  help:"The number of tasks that are running."`
	Sleeping int `json:"sleeping" help:"The number of tasks that are sleeping."`
	Stopped  int `json:"stopped"  help:"The number of tasks that are stopped."`
	Total    int `json:"total"    help:"The total number of tasks."`
	Zombie   int `json:"zombie"   help:"The number of child tasks that are inactive with an active parent task."`
}

// parseOSMetrics parses OS metrics from given JSON data.
func parseOSMetrics(b []byte) (*osMetrics, error) {
	var m osMetrics
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func makeMetric(desc *prometheus.Desc, labelValues []string, value reflect.Value) prometheus.Metric {
	// skip nil fields
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	var f float64
	switch kind := value.Kind(); kind {
	case reflect.Float64:
		f = value.Float()
	case reflect.Int, reflect.Int64:
		f = float64(value.Int())
	default:
		panic(fmt.Errorf("can't make a metric value from %v (%s)", value, kind))
	}

	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, f, labelValues...)
}

func makeGenericMetrics(s interface{}, namePrefix string, constLabels prometheus.Labels) []prometheus.Metric {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		desc := prometheus.NewDesc(namePrefix+name, help, nil, constLabels)
		m := makeMetric(desc, nil, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func makeNodeCPUMetrics(s *cpuUtilization, constLabels prometheus.Labels) []prometheus.Metric {
	// add "cpu" constant label
	labels := prometheus.Labels{"cpu": "All"}
	for k, v := range constLabels {
		labels[k] = v
	}

	// move mode to label
	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	desc := prometheus.NewDesc("node_cpu_average", "The percentage of CPU utilization.", []string{"mode"}, labels)
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		mode := tags.Get("json")
		m := makeMetric(desc, []string{mode}, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func makeDiskIOMetrics(s *diskIO, constLabels prometheus.Labels, device string) []prometheus.Metric {
	// move device name to label
	labelKeys := []string{"device"}
	labelValues := []string{device}

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		if name == "device" {
			continue
		}
		desc := prometheus.NewDesc("rdsosmetrics_diskIO_"+name, help, labelKeys, constLabels)
		m := makeMetric(desc, labelValues, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func makeFileSysMetrics(s *fileSys, constLabels prometheus.Labels, name, mountPoint string) []prometheus.Metric {
	// move name and mount point to labels
	labelKeys := []string{"name", "mount_point"}
	labelValues := []string{name, mountPoint}

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		switch name {
		case "name", "mountPoint":
			continue
		}
		desc := prometheus.NewDesc("rdsosmetrics_fileSys_"+name, help, labelKeys, constLabels)
		m := makeMetric(desc, labelValues, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func makeNodeLoadMetrics(s *loadAverageMinute, constLabels prometheus.Labels) []prometheus.Metric {
	desc := prometheus.NewDesc("node_load1", "The number of processes requesting CPU time over the last minute.", nil, constLabels)
	m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, s.One)
	return []prometheus.Metric{m}
}

func makeNodeMemoryMetrics(s *memory, constLabels prometheus.Labels) []prometheus.Metric {
	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		suffix, multiplierS := tags.Get("node"), tags.Get("m")
		multiplier, err := strconv.ParseInt(multiplierS, 10, 64)
		if err != nil {
			panic(err)
		}
		desc := prometheus.NewDesc("node_memory_"+suffix, "Memory information field "+suffix+".", nil, constLabels)
		m := makeMetric(desc, nil, reflect.ValueOf(v.Field(i).Int()*multiplier))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func makeNetworkMetrics(s *network, constLabels prometheus.Labels, iface string) []prometheus.Metric {
	// move interface name to label
	labelKeys := []string{"interface"}
	labelValues := []string{iface}

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		if name == "interface" {
			continue
		}
		desc := prometheus.NewDesc("rdsosmetrics_network_"+name, help, labelKeys, constLabels)
		m := makeMetric(desc, labelValues, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

//nolint:lll
func makeProcessListMetrics(s *processList, constLabels prometheus.Labels, name string, id, parentID, TGID int) []prometheus.Metric {
	// move process name, ID, parent ID, thread ID to labels
	labelKeys := []string{"name", "id", "parentID", "tgid"}
	labelValues := []string{name, strconv.Itoa(id), strconv.Itoa(parentID), strconv.Itoa(TGID)}

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]prometheus.Metric, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		switch name {
		case "name", "id", "parentID", "tgid":
			continue
		}
		desc := prometheus.NewDesc("rdsosmetrics_processList_"+name, help, labelKeys, constLabels)
		m := makeMetric(desc, labelValues, v.Field(i))
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

func (m *osMetrics) originalMetrics(region string) []prometheus.Metric {
	res := make([]prometheus.Metric, 0, 100)
	constLabels := prometheus.Labels{
		"instance": m.InstanceID,
		"region":   region,
	}

	res = append(res, prometheus.MustNewConstMetric(
		prometheus.NewDesc("rdsosmetrics_timestamp", "Metrics timestamp (UNIX seconds).", nil, constLabels),
		prometheus.CounterValue,
		float64(m.Timestamp.Unix()),
	))

	// make both generic and node_exporter-like metrics
	metrics := makeGenericMetrics(m.CPUUtilization, "rdsosmetrics_cpuUtilization_", constLabels)
	res = append(res, metrics...)
	metrics = makeNodeCPUMetrics(&m.CPUUtilization, constLabels)
	res = append(res, metrics...)

	for _, disk := range m.DiskIO {
		metrics = makeDiskIOMetrics(&disk, constLabels, disk.Device)
		res = append(res, metrics...)
	}

	for _, fs := range m.FileSys {
		metrics = makeFileSysMetrics(&fs, constLabels, fs.Name, fs.MountPoint)
		res = append(res, metrics...)
	}

	// make both generic and node_exporter-like metrics
	metrics = makeGenericMetrics(m.LoadAverageMinute, "rdsosmetrics_loadAverageMinute_", constLabels)
	res = append(res, metrics...)
	metrics = makeNodeLoadMetrics(&m.LoadAverageMinute, constLabels)
	res = append(res, metrics...)

	// make both generic and node_exporter-like metrics
	metrics = makeGenericMetrics(m.Memory, "rdsosmetrics_memory_", constLabels)
	res = append(res, metrics...)
	metrics = makeNodeMemoryMetrics(&m.Memory, constLabels)
	res = append(res, metrics...)

	for _, n := range m.Network {
		metrics = makeNetworkMetrics(&n, constLabels, n.Interface)
		res = append(res, metrics...)
	}

	for _, p := range m.ProcessList {
		metrics = makeProcessListMetrics(&p, constLabels, p.Name, p.ID, p.ParentID, p.TGID)
		res = append(res, metrics...)
	}

	metrics = makeGenericMetrics(m.Swap, "rdsosmetrics_swap_", constLabels)
	res = append(res, metrics...)

	metrics = makeGenericMetrics(m.Tasks, "rdsosmetrics_tasks_", constLabels)
	res = append(res, metrics...)

	return res
}
