package utils

import (
	"bytes"
	"gopkg.in/xmlpath.v2"
)

type Network struct {
	Descr       string `json:"description"`
	Vendor      string `json:"vendor"`
	Product     string `json:"product"`
	LogicalName string `json:"logicalname"`
	Bandwidth   string `json:"bandwidth"`
	MacAddr     string `json:"macaddr"`
}

type Raid struct {
	Descr       string `json:"description"`
	Vendor      string `json:"vendor"`
	Product     string `json:"product"`
	LogicalName string `json:"logicalname"`
	Disk        []Disk `json:"disks"`
}

type Nvme struct {
	Descr       string `josn:"description"`
	Vendor      string `json:"vendor"`
	Size        string `json:"size"`
	Serial      string `json:"serial"`
	Product     string `json:"product"`
	Logicalname string `json:"logicalname"`
}

type DiskVolume struct {
	Descr       string `json:"description"`
	LogicalName string `json:"logicalname"`
	Size        string `json:"size"`
}

type Disk struct {
	Descr       string       `json:"description"`
	Vendor      string       `json:"vendor"`
	Serial      string       `json:"serial"`
	Size        string       `json:"size"`
	Product     string       `json:"product"`
	LogicalName string       `json:"logicalname"`
	Volumes     []DiskVolume `json:"volumes"`
}

type Memory struct {
	Descr   string `json:"description"`
	Physid  string `json:"physid"`
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
	Slot    string `json:"slot"`
	Size    string `json:"size"`
}

type Cpu struct {
	Descr   string `json:"description"`
	Version string `json:"version"`
	Product string `json:"product"`
	Vendor  string `json:"vendor"`
	Slot    string `json:"slot"`
}

type Chassis struct {
	Descr   string `json:"description"`
	Vendor  string `json:"vendor"`
	Serial  string `json:"serial"`
	Product string `json:"product"`
}

type Hardware struct {
	Chassis  Chassis   `json:"chassis"`
	Cpus     []Cpu     `json:"cpus"`
	Memories []Memory  `json:"memories"`
	Disks    []Disk    `json:"disks"`
	Nvmes    []Nvme    `json:"nvme"`
	Raids    []Raid    `json:"raids"`
	Networks []Network `json:"networks"`
}

func XmlParser(search string, root *xmlpath.Node) string {
	path := xmlpath.MustCompile(search)
	if value, ok := path.String(root); ok {
		return value
	}
	return ""
}

func ChassisParser(root *xmlpath.Node) Chassis {

	chassis := Chassis{}

	chassis.Vendor = XmlParser("//*/node/vendor", root)
	chassis.Product = XmlParser("//*/node/product", root)
	chassis.Descr = XmlParser("//*/node/description", root)
	chassis.Serial = XmlParser("//*/node/serial", root)

	return chassis

}

func CpuParser(root *xmlpath.Node) []Cpu {

	path := xmlpath.MustCompile("//*/node[@id='core']/node[contains(@id,'cpu')]")
	cpuRoots := path.Iter(root)
	cpus := []Cpu{}

	for cpuRoots.Next() {
		cpuRoot := cpuRoots.Node()
		cpu := Cpu{}
		cpu.Descr = XmlParser("description", cpuRoot)
		cpu.Product = XmlParser("product", cpuRoot)
		cpu.Version = XmlParser("version", cpuRoot)
		cpu.Vendor = XmlParser("vendor", cpuRoot)
		cpu.Slot = XmlParser("slot", cpuRoot)

		cpus = append(cpus, cpu)
	}

	return cpus

}

func MemoryParser(root *xmlpath.Node) []Memory {

	path := xmlpath.MustCompile("//*/node[@id='core']/node[contains(@id,'memory')]/node[contains(@id,'bank')]")
	memoryRoots := path.Iter(root)
	memories := []Memory{}

	for memoryRoots.Next() {
		memoryRoot := memoryRoots.Node()
		memory := Memory{}

		memory.Descr = XmlParser("description", memoryRoot)
		memory.Product = XmlParser("product", memoryRoot)
		memory.Vendor = XmlParser("vendor", memoryRoot)
		memory.Physid = XmlParser("physid", memoryRoot)
		memory.Slot = XmlParser("slot", memoryRoot)
		memory.Size = XmlParser("size", memoryRoot)

		memories = append(memories, memory)
	}

	return memories

}

func DiskParser(root *xmlpath.Node) []Disk {

	path := xmlpath.MustCompile("//*/node[contains(@id,'disk')]")
	diskRoots := path.Iter(root)
	disks := []Disk{}

	for diskRoots.Next() {
		diskRoot := diskRoots.Node()
		disk := Disk{}

		disk.Descr = XmlParser("description", diskRoot)
		disk.Product = XmlParser("product", diskRoot)
		disk.Vendor = XmlParser("vendor", diskRoot)
		disk.Serial = XmlParser("serial", diskRoot)
		disk.Size = XmlParser("size", diskRoot)
		disk.LogicalName = XmlParser("logicalname", diskRoot)

		path = xmlpath.MustCompile("node[contains(@id, 'volume')]")
		volumeRoots := path.Iter(diskRoot)
		volumes := []DiskVolume{}

		for volumeRoots.Next() {
			volumeRoot := volumeRoots.Node()
			volume := DiskVolume{}

			volume.Descr = XmlParser("description", volumeRoot)
			volume.LogicalName = XmlParser("logicalname", volumeRoot)
			volume.Size = XmlParser("capacity", volumeRoot)

			volumes = append(volumes, volume)
		}
		disk.Volumes = volumes

		disks = append(disks, disk)
	}

	return disks

}

func RaidParser(root *xmlpath.Node) []Raid {

	path := xmlpath.MustCompile("//*/node[contains(@id,'raid')]")
	raidRoots := path.Iter(root)
	raids := []Raid{}

	for raidRoots.Next() {
		raidRoot := raidRoots.Node()
		raid := Raid{}

		raid.Descr = XmlParser("description", raidRoot)
		raid.Product = XmlParser("product", raidRoot)
		raid.Vendor = XmlParser("vendor", raidRoot)
		raid.LogicalName = XmlParser("logicalname", raidRoot)

		path = xmlpath.MustCompile("node[contains(@id, 'disk')]")
		disksRoots := path.Iter(raidRoot)
		disks := []Disk{}

		for disksRoots.Next() {
			diskRoot := disksRoots.Node()
			disk := Disk{}

			disk.Descr = XmlParser("description", diskRoot)
			disk.Product = XmlParser("product", diskRoot)
			disk.Vendor = XmlParser("vendor", diskRoot)
			disk.Serial = XmlParser("serial", diskRoot)
			disk.LogicalName = XmlParser("logicalname", diskRoot)
			disk.Size = XmlParser("capacity", diskRoot)

			disks = append(disks, disk)
		}

		raid.Disk = disks
		raids = append(raids, raid)
	}

	return raids

}

func NetworkParser(root *xmlpath.Node) []Network {

	path := xmlpath.MustCompile("//*/node[contains(@id,'network')]")
	networkRoots := path.Iter(root)
	networks := []Network{}

	for networkRoots.Next() {
		networkRoot := networkRoots.Node()
		network := Network{}

		network.Descr = XmlParser("description", networkRoot)
		network.Product = XmlParser("product", networkRoot)
		network.Vendor = XmlParser("vendor", networkRoot)
		network.LogicalName = XmlParser("logicalname", networkRoot)
		network.Bandwidth = XmlParser("capacity", networkRoot)
		network.MacAddr = XmlParser("serial", networkRoot)

		networks = append(networks, network)
	}

	return networks

}

func LshwParser(xml []byte) Hardware {

	var hw Hardware

	root, _ := xmlpath.Parse(bytes.NewReader(xml))

	hw.Chassis = ChassisParser(root)

	hw.Cpus = CpuParser(root)

	hw.Memories = MemoryParser(root)

	hw.Disks = DiskParser(root)

	hw.Raids = RaidParser(root)

	hw.Networks = NetworkParser(root)

	//b, _ := json.Marshal(hw)

	//fmt.Println(string(b))
	//fmt.Printf("%#v", hw)
	//pp.Print(hw)

	return hw
}
