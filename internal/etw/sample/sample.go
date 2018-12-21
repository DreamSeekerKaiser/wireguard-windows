// Shows a sample usage of the ETW logging package.
package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Microsoft/go-winio/internal/etw"
	"github.com/sirupsen/logrus"

	"golang.org/x/sys/windows"
)

func callback(sourceID *windows.GUID, state etw.ProviderState, level etw.Level, matchAnyKeyword uint64, matchAllKeyword uint64, filterData uintptr) {
	fmt.Printf("Callback: isEnabled=%d, level=%d, matchAnyKeyword=%d\n", state, level, matchAnyKeyword)
}

func guidToString(guid *windows.GUID) string {
	data1 := make([]byte, 4)
	binary.BigEndian.PutUint32(data1, guid.Data1)
	data2 := make([]byte, 2)
	binary.BigEndian.PutUint16(data2, guid.Data2)
	data3 := make([]byte, 2)
	binary.BigEndian.PutUint16(data3, guid.Data3)
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hex.EncodeToString(data1),
		hex.EncodeToString(data2),
		hex.EncodeToString(data3),
		hex.EncodeToString(guid.Data4[:2]),
		hex.EncodeToString(guid.Data4[2:]))
}

func main() {
	provider, err := etw.NewProvider("TestProvider", callback)

	if err != nil {
		logrus.Error(err)
		return
	}
	defer func() {
		if err := provider.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	fmt.Println("Provider ID:", guidToString(provider.ID))

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Press enter to log events")
	reader.ReadString('\n')

	// Write using high-level API.
	if err := provider.WriteEvent(
		"TestEvent",
		etw.WithLevel(etw.LevelInfo),
		etw.WithKeyword(0x140),
		etw.StringField("TestField", "Foo"),
		etw.StringField("TestField2", "Bar"),
		etw.Struct("TestStruct",
			etw.StringField("Field1", "Value1"),
			etw.StringField("Field2", "Value2")),
		etw.StringArray("TestArray", []string{
			"Item1",
			"Item2",
			"Item3",
			"Item4",
			"Item5",
		}),
	); err != nil {
		logrus.Error(err)
		return
	}

	// Write using low-level API.
	descriptor := etw.NewEventDescriptor()
	descriptor.Level = etw.LevelInfo
	descriptor.Keyword = 0x140
	em := &etw.EventMetadata{}
	ed := &etw.EventData{}
	em.WriteEventHeader("TestEvent", 0)
	em.WriteField("TestField", etw.InTypeANSIString, etw.OutTypeUTF8, 0)
	ed.WriteString("Foo")
	em.WriteField("TestField2", etw.InTypeANSIString, etw.OutTypeUTF8, 0)
	ed.WriteString("Bar")
	em.WriteStruct("TestStruct", 2, 0)
	em.WriteField("Field1", etw.InTypeANSIString, etw.OutTypeUTF8, 0)
	ed.WriteString("Value1")
	em.WriteField("Field2", etw.InTypeANSIString, etw.OutTypeUTF8, 0)
	ed.WriteString("Value2")
	em.WriteArray("TestArray", etw.InTypeANSIString, etw.OutTypeDefault, 0)
	ed.WriteUint16(5)
	ed.WriteString("Item1")
	ed.WriteString("Item2")
	ed.WriteString("Item3")
	ed.WriteString("Item4")
	ed.WriteString("Item5")
	if err := provider.WriteEventRaw(descriptor, [][]byte{em.Bytes()}, [][]byte{ed.Bytes()}); err != nil {
		logrus.Error(err)
		return
	}
}
