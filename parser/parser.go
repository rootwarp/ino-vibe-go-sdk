package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Errors
var (
	ErrNoHeader      = errors.New("Parse Header first")
	ErrInvalidFormat = errors.New("Invalid Frame format")
	ErrInvalidFrame  = errors.New("Invalid Raw Data")
	ErrFrameVersion  = errors.New("Frame version not supported")
	ErrInvalidType   = errors.New("Invalid Payload type")
)

// Header is header of full frame.
type Header struct {
	Version     uint32     `json:"version"`
	DevType     DeviceType `json:"dev_type"`
	Seq         uint32     `json:"seq"`
	Battery     uint32     `json:"battery"`
	Temperature int32      `json:"temperature"`
	LoRaErr     uint32     `json:"lora_err"`
	RSSI        int32      `json:"rssi"`
	Payload     Payload    `json:"payload"`
	Resv        uint32     `json:"resv"`
}

// DeviceType is InoVibe hardware type.
type DeviceType uint

// DeviceTypes
const (
	InoVibe  DeviceType = 2
	InoVibeS DeviceType = 3
)

// Payload is detail frame values.
type Payload struct {
	Type    PayloadType `json:"type"`
	Request uint32      `json:"request"`
}

// PayloadType is several payload types.
type PayloadType uint

// PayloadTypes
const (
	UnknownType     PayloadType = 0
	AliveType       PayloadType = 1
	EventType       PayloadType = 2
	ErrorType       PayloadType = 3
	AckType         PayloadType = 4
	NoticeType      PayloadType = 5
	DataLogType     PayloadType = 6
	ReportType      PayloadType = 7
	WaveType        PayloadType = 8
	InclinationType PayloadType = 9
	MRMeasureType   PayloadType = 10
	MRReportType    PayloadType = 11
)

// NoticePayloadType is types of notice frame.
type NoticePayloadType uint

// NoticeTypes
const (
	NoticePowerUp           NoticePayloadType = 1
	NoticePowerOff          NoticePayloadType = 2
	NoticeSetup             NoticePayloadType = 4
	NoticeTestResult        NoticePayloadType = 5
	NoticeRejectCount       NoticePayloadType = 6
	NoticeApplicationConfig NoticePayloadType = 7
)

func (v PayloadType) key() string {
	keys := []string{
		"non-exist",
		"alive",
		"event",
		"error",
		"ack",
		"notice",
		"datalog",
		"report",
		"wave",
		"inclination",
		"mr_measure",
		"mr_report",
	}
	return keys[v]
}

// AlivePayload is periodical frame from device.
type AlivePayload struct {
	X           int            `json:"x"`
	Y           int            `json:"y"`
	Z           int            `json:"z"`
	AlivePeriod uint           `json:"alive_period"`
	Sensitivity AccSensitivity `json:"sensitivity"`
	Threshold   uint           `json:"threshold"`
	AccIntNo    uint           `json:"acc_int_no"`
	AccIntResv  uint           `json:"acc_int_resv"`
	AccIntData  uint           `json:"acc_int_data"`
	LogEnable   uint           `json:"log_enable"`
	LogInterval uint           `json:"log_interval"`
	LogBlocks   uint           `json:"log_blocks"`
	Setup       DeviceSetup    `json:"setup"`
	AppFwMajor  uint           `json:"app_fw_major"`
	AppFwMinor  uint           `json:"app_fw_minor"`
	AppFwRev    uint           `json:"app_fw_rev"`
	LoRaFwMajor uint           `json:"lora_fw_major"`
	LoRaFwMinor uint           `json:"lora_fw_minor"`
	LoRaFwRev   uint           `json:"lora_fw_rev"`
}

// AccSensitivity is accelerometer's range.
type AccSensitivity uint

// Accelerometer's ranges.
const (
	AccSensitivity2G  AccSensitivity = 1
	AccSensitivity4G  AccSensitivity = 2
	AccSensitivity8G  AccSensitivity = 3
	AccSensitivity16G AccSensitivity = 4
)

// DeviceSetup describes currrent install status of device.
type DeviceSetup uint

// DeviceSetup values.
const (
	DeviceSetupUninstalled    DeviceSetup = 0
	DeviceSetupInstalled      DeviceSetup = 1
	DeviceSetupPrepareInstall DeviceSetup = 2
)

// WavePayload contains accelerometer raw data.
type WavePayload struct {
	Control  WaveControl  `json:"control"`
	PackType WavePackType `json:"pack_type"`
	Position uint         `json:"pos"`
	X        []int        `json:"-"`
	Y        []int        `json:"-"`
	Z        []int        `json:"-"`
}

// WaveControl is wave frame control spec.
type WaveControl struct {
	BMARange WaveBMARangeType `json:"bma_range"`
	Axis     WaveAxisType     `json:"axis"`
	ID       uint             `json:"id"`
}

// WaveBMARangeType is accelerometer config.
type WaveBMARangeType uint

// WaveBMARangeTypes
const (
	WaveBMARange2G  WaveBMARangeType = 0
	WaveBMARange4G  WaveBMARangeType = 1
	WaveBMARange8G  WaveBMARangeType = 2
	WaveBMARange16G WaveBMARangeType = 3
)

// WaveAxisType is accelerometer axis config.
type WaveAxisType uint

// WaveAxisTypes
const (
	WaveAxisXYZ WaveAxisType = 0
	WaveAxisX   WaveAxisType = 1
	WaveAxisY   WaveAxisType = 2
	WaveAxisZ   WaveAxisType = 3
)

// WavePackType is byte align type of wave.
type WavePackType uint

// WavePackTypes
const (
	WavePack16       WavePackType = 0
	WavePack12       WavePackType = 1
	WavePack16Finish WavePackType = 0xF
)

// ConfigBMARange is accelerometer gravity range.
type ConfigBMARange int

// ConfigBMARanges
const (
	ConfigBMARange2G  ConfigBMARange = 1
	ConfigBMARange4G  ConfigBMARange = 2
	ConfigBMARange8G  ConfigBMARange = 3
	ConfigBMARange16G ConfigBMARange = 4
)

// MRMTState is Machine Runtime stages.
type MRMTState int

// MRMTStates
const (
	MRMTStateCommission MRMTState = 0
	MRMTStateInactive   MRMTState = 1
	MRMTStateActive     MRMTState = 2
)

// Preventing stackoverflow from recursion.
const allowNest = 5

func loadFormat() map[string]interface{} {
	var formatMap map[string]interface{}

	err := json.Unmarshal([]byte(frameSpec), &formatMap)
	if err != nil {
		fmt.Println("JSON parse fail ", err)
		return nil
	}

	return formatMap
}

func parseFrame(frameMap map[string]interface{}, rawFrame string, startOffset int, depth int) (map[string]interface{}, int, error) {
	if depth > allowNest {
		return nil, 0, nil
	}

	frameFields := frameMap["fields"].([]interface{})

	binaryStr := ""

	for _, c := range rawFrame {
		parsed, _ := strconv.ParseUint(fmt.Sprintf("%c", c), 16, 4)
		binaryStr += fmt.Sprintf("%04b", parsed)
	}

	parsedMap := map[string]interface{}{}
	offsetBits := startOffset
	for _, fieldEntry := range frameFields {
		fieldMap := fieldEntry.(map[string]interface{})
		fieldName := fieldMap["name"].(string)

		if _, ok := fieldMap["fields"]; ok {
			// Subfields
			nestedMap, offset, err := parseFrame(fieldMap, rawFrame, offsetBits, depth+1)
			if err != nil {
				return nil, offsetBits, err
			}

			parsedMap[fieldName] = nestedMap

			offsetBits = offset
		} else if _, ok := fieldMap["subtype"]; ok {
			// Subtype field
			supportSubTypes := fieldMap["subtype"].(map[string]interface{})
			subtype := strconv.FormatUint(parsedMap["type"].(uint64), 10)
			subtypes := frameMap["subtypes"].(map[string]interface{})

			subtypeName, ok := supportSubTypes[subtype].(string)
			if !ok {
				return nil, 0, ErrInvalidType
			}
			subtypeMap := subtypes[subtypeName].(map[string]interface{})

			nestedMap, offset, err := parseFrame(subtypeMap, rawFrame, offsetBits, depth+1)
			if err != nil {
				return nil, offset, err
			}

			for key, value := range nestedMap {
				parsedMap[key] = value
			}
		} else {
			// Normal field
			fieldBits := int(fieldMap["bits"].(float64))
			startIdx := offsetBits
			endIdx := (offsetBits + fieldBits)
			if len(binaryStr) < (endIdx - 1) {
				return nil, 0, ErrInvalidFrame
			}

			fieldStr := binaryStr[startIdx:endIdx]

			parsedValue, _ := strconv.ParseUint(fieldStr, 2, fieldBits)

			if fieldMap["signed"].(bool) {
				max := (1<<uint(fieldBits-1) - 1)
				if int64(parsedValue) > int64(max) {
					max := 1 << uint(fieldBits)
					parsedMap[fieldName] = -(int64(max) - int64(parsedValue))
				} else {
					parsedMap[fieldName] = int64(parsedValue)
				}
			} else {
				parsedMap[fieldName] = parsedValue
			}

			offsetBits += fieldBits
		}
	}
	return parsedMap, offsetBits, nil
}

// FrameParser parses raw device frames into structures.
type FrameParser interface {
	Header() (*Header, error)
	Alive() (*AlivePayload, error)
	Wave() (*WavePayload, error)
	Notice() (interface{}, error)
}

type frameParser struct {
	Raw        string
	RawHeader  string
	RawPayload string
	header     *Header
	frameMap   map[string]interface{}
}

func (p *frameParser) Header() (*Header, error) {
	headerMaps := p.frameMap["header"].(map[string]interface{})
	headerFields, _, err := parseFrame(headerMaps, p.Raw, 0, 0)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(headerFields)
	if err != nil {
		fmt.Println(err)
	}

	header := Header{}
	json.Unmarshal(data, &header)
	p.header = &header
	return &header, nil
}

func (p *frameParser) Alive() (*AlivePayload, error) {
	if p.header == nil {
		return nil, ErrNoHeader
	}

	payloadMap := p.frameMap[p.header.Payload.Type.key()].(map[string]interface{})
	aliveMap, _, err := parseFrame(payloadMap, p.RawPayload, 0, 0)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(aliveMap)
	if err != nil {
		fmt.Println(err)
	}

	payload := AlivePayload{}
	json.Unmarshal(data, &payload)
	return &payload, nil
}

func (p *frameParser) Wave() (*WavePayload, error) {
	if p.header == nil {
		return nil, ErrNoHeader
	}

	payloadMap := p.frameMap[p.header.Payload.Type.key()].(map[string]interface{})
	waveMap, _, err := parseFrame(payloadMap, p.RawPayload, 0, 0)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(waveMap)
	if err != nil {
		fmt.Println(err)
	}

	payload := WavePayload{}
	json.Unmarshal(data, &payload)

	accFrame := p.RawPayload[4 : len(p.RawPayload)-2]

	binaryStr := ""
	for _, hex := range accFrame {
		parsed, _ := strconv.ParseUint(fmt.Sprintf("%c", hex), 16, 4)
		binaryStr += fmt.Sprintf("%04b", parsed)
	}

	var (
		bitsPerFrame uint = 16
		scale             = 1
	)

	if payload.PackType == WavePack12 {
		bitsPerFrame = 12
		scale = 4
	}

	var frameCount = len(binaryStr) / int(bitsPerFrame)
	if payload.Control.Axis == WaveAxisXYZ {
		frameCount /= 3
	}

	// TODO: Need define repeat field?
	for i := 0; i < frameCount; i++ {
		if payload.Control.Axis == WaveAxisXYZ {
			valueStr := binaryStr[i*int(bitsPerFrame) : (i+1)*int(bitsPerFrame)]
			accValue, _ := strconv.ParseUint(valueStr, 2, int(bitsPerFrame))
			payload.X = append(payload.X, convertSignedInt(accValue, bitsPerFrame)*scale)

			valueStr = binaryStr[(i+1)*int(bitsPerFrame) : (i+2)*int(bitsPerFrame)]
			accValue, _ = strconv.ParseUint(valueStr, 2, int(bitsPerFrame))
			payload.Y = append(payload.Y, convertSignedInt(accValue, bitsPerFrame)*scale)

			valueStr = binaryStr[(i+2)*int(bitsPerFrame) : (i+3)*int(bitsPerFrame)]
			accValue, _ = strconv.ParseUint(valueStr, 2, int(bitsPerFrame))
			payload.Z = append(payload.Z, convertSignedInt(accValue, bitsPerFrame)*scale)
		} else {
			var slice *[]int
			if payload.Control.Axis == WaveAxisX {
				slice = &payload.X
			} else if payload.Control.Axis == WaveAxisY {
				slice = &payload.Y
			} else {
				slice = &payload.Z
			}

			valueStr := binaryStr[i*int(bitsPerFrame) : (i+1)*int(bitsPerFrame)]
			accValue, _ := strconv.ParseUint(valueStr, 2, int(bitsPerFrame))

			*slice = append(*slice, convertSignedInt(accValue, bitsPerFrame)*scale)
		}
	}

	return &payload, nil
}

func convertSignedInt(value uint64, bitsPerFrame uint) int {
	maxValue := int(1<<(bitsPerFrame-1)) - 1
	if int(value) > maxValue {
		return -(int(1<<(bitsPerFrame)) - int(value))
	}
	return int(value)
}

func (p *frameParser) Notice() (interface{}, error) {
	if p.header == nil {
		return nil, ErrNoHeader
	}

	payloadMap := p.frameMap[p.header.Payload.Type.key()].(map[string]interface{})
	noticeMap, _, err := parseFrame(payloadMap, p.RawPayload, 0, 0)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(noticeMap)
	if err != nil {
		return nil, err
	}

	noticeCommon := struct {
		Type   NoticePayloadType `json:"type"`
		Length int               `json:"length"`
	}{}

	err = json.Unmarshal(data, &noticeCommon)

	switch noticeCommon.Type {
	case NoticePowerUp:
		payload := PowerUp{}
		err = json.Unmarshal(data, &payload)
		return payload, err
	case NoticeSetup:
		payload := Setup{}
		err = json.Unmarshal(data, &payload)
		return payload, err
	case NoticeRejectCount:
		payload := RejectCount{}
		err = json.Unmarshal(data, &payload)
		return payload, err
	case NoticeApplicationConfig:
		payload := ApplicationConfig{}
		err = json.Unmarshal(data, &payload)
		return payload, err
	default:
		return nil, ErrInvalidType
	}
}

// NewFrameParser creates new parser.
func NewFrameParser(raw string) (FrameParser, error) {
	if len(raw) < 20 {
		return nil, ErrInvalidFrame
	}

	if ver, _ := strconv.ParseInt(raw[:2], 16, 8); ver != 3 {
		return nil, ErrFrameVersion
	}

	parser := frameParser{
		Raw:        raw,
		RawHeader:  raw[:20],
		RawPayload: raw[20:],
		frameMap:   loadFormat(),
	}

	return &parser, nil
}
