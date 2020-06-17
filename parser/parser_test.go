package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAlive(t *testing.T) {
	/*
		V3 | Dev.: mgi | Seq.: 244 | Bat.: 17 | Temp.: 21 | Err.: 0 | Type: alive, none
		| RSSI.: -110 | Resv.: 100

		ALIVE | Acc Conf.: 2G, 1100mg | Acc Intr.: 1, 1 | Acc Vals.: 59, -249, -29
		| Log: EN(XYZ), 1s, 12blks | Inst.: Y | Per: 360m | App FW: 2.6.3 | LoRa FW: 1.2.2

		RAW | 0302f411150092100064003bff07ffe3016801044c01000101010c0102060301020223
	*/
	raw := "0302f411150092100064003bff07ffe3016801044c01000101010c0102060301020223"

	parser, _ := NewFrameParser(raw)
	header, _ := parser.Header()
	alive, _ := parser.Alive()

	fmt.Printf("%+v %+v", header, alive)

	assert.Equal(t, uint32(3), header.Version)
	assert.Equal(t, InoVibe, header.DevType)
	assert.Equal(t, uint32(244), header.Seq)
	assert.Equal(t, uint32(17), header.Battery)
	assert.Equal(t, int32(21), header.Temperature)
	assert.Equal(t, uint32(0), header.LoRaErr)
	assert.Equal(t, int32(-110), header.RSSI)
	assert.Equal(t, AliveType, header.Payload.Type)
	assert.Equal(t, uint32(0), header.Payload.Request)
	assert.Equal(t, uint32(100), header.Resv)

	assert.Equal(t, 59, alive.X)
	assert.Equal(t, -249, alive.Y)
	assert.Equal(t, -29, alive.Z)
	assert.Equal(t, uint(360), alive.AlivePeriod)
	assert.Equal(t, AccSensitivity2G, alive.Sensitivity)
	assert.Equal(t, uint(1100), alive.Threshold)
	assert.Equal(t, uint(1), alive.AccIntNo)
	assert.Equal(t, uint(1), alive.AccIntData)
	assert.Equal(t, uint(1), alive.LogEnable)
	assert.Equal(t, uint(1), alive.LogInterval)
	assert.Equal(t, uint(12), alive.LogBlocks)
	assert.Equal(t, DeviceSetupInstalled, alive.Setup)
	assert.Equal(t, uint(2), alive.AppFwMajor)
	assert.Equal(t, uint(6), alive.AppFwMinor)
	assert.Equal(t, uint(3), alive.AppFwRev)
	assert.Equal(t, uint(1), alive.LoRaFwMajor)
	assert.Equal(t, uint(2), alive.LoRaFwMinor)
	assert.Equal(t, uint(2), alive.LoRaFwRev)
}

func TestParseInvalidShortFrame(t *testing.T) {
	raw := "0302f411150" // Short frame
	parser, err := NewFrameParser(raw)

	assert.Nil(t, parser)
	assert.Equal(t, ErrInvalidFrame, err)
}

func TestParseWrongAlive(t *testing.T) {
	raw := "0302f411150092100064003bff07ffe3016801044c01000101010c010206030102"
	parser, _ := NewFrameParser(raw)
	parser.Header()
	alive, err := parser.Alive()

	assert.Nil(t, alive)
	assert.Equal(t, ErrInvalidFrame, err)
}

func TestParseWrongVersion(t *testing.T) {
	raw := "0402f411150092100064003bff07ffe3016801044c01000101010c010206030102" // Ver 4 frame.
	parser, err := NewFrameParser(raw)

	assert.Nil(t, parser)
	assert.Equal(t, ErrFrameVersion, err)
}

// TODO: Wrong format

func TestParseWave(t *testing.T) {
	/*
		V3 | Dev.: mgi | Seq.: 227 | Bat.: 1 | Temp.: 19 | Err.: 0 | Type: acc_wave, none | RSSI.: -99 | Resv.: 13
		WAVE | Rng: 2G | Axis: x | ID: 12 | Seq: 1
		RAW | 0302e30113009d80000d1c01ffdfffd800e8ffedffbc00e4fff6ffd300d0fff2ffe400e3fff3ffd000effff2ffe200f4ffe7fff400effff8ffee00d3ae



		V3 | Dev.: mgi | Seq.: 73 | Bat.: 1 | Temp.: 20 | Err.: 0 | Type: acc_wave, none | RSSI.: -116 | Resv.: 15
		WAVE | Rng: 2G | Axis: y | ID: 3 | Seq: 255
		RAW | 0302490114008c80000f23ffffd4fff600ebffd0ffde00f6ffd2ffd600edffd8ffe300faffcbffd400fbffcdffd900e5ffd7ffde00daffe0ffe400d47f

		V3 | Dev.: mgi | Seq.: 241 | Bat.: 96 | Temp.: 22 | Err.: 0 | Type: acc_wave, none | RSSI.: -122 | Resv.: 2661
		WAVE | Rng: 2G | Axis: z | ID: 4 | Seq: 255
		RAW | 0302f160160086800a6534ff001b0002fffa001bfffffffe001cfff600020016fff5000a0013fff800110010fff500120009fffb00140005fffe0015c6


		V3 | Dev.: mgi_100n | Seq.: 219 | Bat.: 92 | Temp.: 29 | Err.: 0 | Type: acc_wave, none | RSSI.: -94 | Resv.: 2598
		WAVE | Rng: 2G | Axis: x y z | ID: 4 | Seq: 255
		RAW | 0303db5c1d00a2800a2604ff0fc2fe6dfef80fdafe8afee10fe3fe6afefb100ffe6cfee3100ffe80fee51007fe6cfeea0ff3fe78feec0fdefe72fed367

		V3 | Dev.: mgi_100n | Seq.: 15 | Bat.: 87 | Temp.: 28 | Err.: 0 | Type: acc_wave, none | RSSI.: -44 | Resv.: 2970
		WAVE | Rng: 2G | Axis: z | ID: 5 | Seq: 0
		RAW |03030f571c00d4800b9a35000d9c164411f9135d10bb111b10fa1110105210190fdd0fde1095110f125e12881193101d59

	*/

	fixtures := []map[string]interface{}{
		{
			"raw":       "0302e30113009d80000d1c01ffdfffd800e8ffedffbc00e4fff6ffd300d0fff2ffe400e3fff3ffd000effff2ffe200f4ffe7fff400effff8ffee00d3ae",
			"range":     0,
			"axis":      1,
			"id":        12,
			"pack_type": WavePack16,
			"pos":       1,
			"x_cnt":     24,
			"y_cnt":     0,
			"z_cnt":     0,
		},
		{
			"raw":       "0302490114008c80000f23ffffd4fff600ebffd0ffde00f6ffd2ffd600edffd8ffe300faffcbffd400fbffcdffd900e5ffd7ffde00daffe0ffe400d47f",
			"range":     0,
			"axis":      2,
			"id":        3,
			"pack_type": WavePack16Finish,
			"pos":       0xF,
			"x_cnt":     0,
			"y_cnt":     24,
			"z_cnt":     0,
		},
		{
			"raw":       "0302f160160086800a6534ff001b0002fffa001bfffffffe001cfff600020016fff5000a0013fff800110010fff500120009fffb00140005fffe0015c6",
			"range":     0,
			"axis":      3,
			"id":        4,
			"pack_type": WavePack16Finish,
			"pos":       0xF,
			"x_cnt":     0,
			"y_cnt":     0,
			"z_cnt":     24,
		},
		{
			"raw":       "0303db5c1d00a2800a2604ff0fc2fe6dfef80fdafe8afee10fe3fe6afefb100ffe6cfee3100ffe80fee51007fe6cfeea0ff3fe78feec0fdefe72fed367",
			"range":     0,
			"axis":      0,
			"id":        4,
			"pack_type": WavePack16Finish,
			"pos":       0xF,
			"x_cnt":     8,
			"y_cnt":     8,
			"z_cnt":     8,
		},

		{
			"raw":       "03030f571c00d4800b9a35000d9c164411f9135d10bb111b10fa1110105210190fdd0fde1095110f125e12881193101d59",
			"range":     0,
			"axis":      3,
			"id":        5,
			"pack_type": WavePack16,
			"pos":       0,
			"x_cnt":     0,
			"y_cnt":     0,
			"z_cnt":     18,
		},
		{
			"raw":       "0302f2642300a1800b13f40001b5fe00fe00010d0032012afe00fe0001ff0188015dffdeff53ffb000c100210067ffa7ff2a005000a800e6ffb0ff8cbf",
			"range":     3,
			"axis":      3,
			"id":        4,
			"pack_type": WavePack16,
			"pos":       0,
			"x_cnt":     0,
			"y_cnt":     0,
			"z_cnt":     24,
		},
		// New typed. 1 packet frame.
		{
			"raw":       "030335641900e1800bfa311f3824b544dd2d1a74a041d44f3364d44b81f837c5733c03e349ba757ff4c748e4db38e5766b",
			"range":     0,
			"axis":      3,
			"id":        1,
			"pack_type": WavePack12,
			"pos":       0xF,
			"x_cnt":     0,
			"y_cnt":     0,
			"z_cnt":     24,
		},
	}

	for _, fixture := range fixtures {
		parser, _ := NewFrameParser(fixture["raw"].(string))
		header, _ := parser.Header()
		wave, _ := parser.Wave()

		assert.Equal(t, WaveType, header.Payload.Type)
		assert.Equal(t, WaveBMARangeType(fixture["range"].(int)), wave.Control.BMARange)
		assert.Equal(t, WaveAxisType(fixture["axis"].(int)), wave.Control.Axis)
		assert.Equal(t, uint(fixture["id"].(int)), wave.Control.ID)
		assert.Equal(t, fixture["pack_type"].(WavePackType), wave.PackType)
		assert.Equal(t, uint(fixture["pos"].(int)), wave.Position)

		assert.Equal(t, fixture["x_cnt"], len(wave.X))
		assert.Equal(t, fixture["y_cnt"], len(wave.Y))
		assert.Equal(t, fixture["z_cnt"], len(wave.Z))
	}
}

func TestParserWaveData(t *testing.T) {
	fixtures := []struct {
		Raw          string
		ExpectValues []int
	}{
		{
			Raw: "0303a530190090800ae03000105310561054104e1050106310671054105d105f105a10541066104a104e105a105f1055105b105f10661061105e105905",
			ExpectValues: []int{4179, 4182, 4180, 4174, 4176, 4195, 4199, 4180,
				4189, 4191, 4186, 4180, 4198, 4170, 4174, 4186,
				4191, 4181, 4187, 4191, 4198, 4193, 4190, 4185},
		},
		{
			Raw: "030337641900e6800bfa361f4193ffc3c53244d4935515e84b60951e66d64283ff49630138e4903c14613d9408412416c0",
			ExpectValues: []int{4196, 4092, -3856, 5320, 4404, 4684, 5444, 6048,
				4824, 596, 1944, 7000, 4256, 4092, 4696, 3076,
				3640, 4672, 3844, 4484, 3940, 4128, 4168, 4184},
		},
	}

	for _, fixture := range fixtures {
		parser, _ := NewFrameParser(fixture.Raw)
		_, _ = parser.Header()
		wave, _ := parser.Wave()

		assert.Equal(t, fixture.ExpectValues, wave.Z)
	}
}

func TestParserNoticePowerUp(t *testing.T) {
	/*
		V3 | Dev.: mgi_100n | Seq.: 22 | Bat.: 100 | Temp.: 30 | Err.: 0 | Type: notice, config | RSSI.: -99 | Resv.: 3010
		ON | Reset(0): -- | TurnOnCnt: 4 | Err(12): PowefOffUninstallCommand
		RAW | 030316641e009d520bc201040004000c6f
	*/
	raw := "030316641e009d520bc201040004000c6f"

	parser, _ := NewFrameParser(raw)
	parser.Header()

	notice, err := parser.Notice()
	powerup, ok := notice.(PowerUp)

	assert.True(t, ok)
	assert.Nil(t, err)

	assert.Equal(t, 0, powerup.ResetReason)
	assert.Equal(t, 4, powerup.Count)
	assert.Equal(t, 12, powerup.OffReason)
}

func TestParserNoticeSetup(t *testing.T) {
	/*
		V3 | Dev.: mgi_100n | Seq.: 27 | Bat.: 100 | Temp.: 30 | Err.: 0 | Type: notice, none | RSSI.: -101 | Resv.: 2994
		SETUP | 2 -> 1
		RAW | 03031b641e009b500bb20402010254
	*/
	raw := "03031b641e009b500bb20402010254"

	parser, _ := NewFrameParser(raw)
	parser.Header()

	notice, err := parser.Notice()
	setup, ok := notice.(Setup)

	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, DeviceSetupPrepareInstall, setup.Previous)
	assert.Equal(t, DeviceSetupInstalled, setup.Current)
}

func TestParserNoticeRejection(t *testing.T) {
	/*
		V3 | Dev.: mgi_100n | Seq.: 42 | Bat.: 100 | Temp.: 28 | Err.: 0 | Type: notice, none | RSSI.: -86 | Resv.: 3002
		REJECT | Count: 77 | Period: 24 | Threshold: 0
		RAW | 03032a641c00aa500bba0604004d1800de
	*/

	raw := "03032a641c00aa500bba0604004d1800de"

	parser, _ := NewFrameParser(raw)
	parser.Header()

	notice, err := parser.Notice()
	reject, ok := notice.(RejectCount)

	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, 77, reject.Count)
	assert.Equal(t, 24, reject.Period)
	assert.Equal(t, 0, reject.Threshold)
}

func TestParseNoticeConfig(t *testing.T) {
	/*
		V3 | Dev.: mgi | Seq.: 211 | Bat.: 100 | Temp.: 26 | Err.: 0 | Type: notice, none | RSSI.: -54 | Resv.: 3207
		APPCFG | Mode: impact | BMA RNG: 16G | HighG THR: 1900mg | # Sub INTVL: 4 | nRF REJ THR: 12mg | nRF IMPACT THR: 320mg | INCL CK PER: 0M | # TX SKIP: 0 | Base: Raw(0, 0, 0) / Deg(--, --, --) | MRMT State: commission | MRMT OP THR: 0mg | MRMT SHOCK THR: 0mg
		RAW | 0302d3641a00ca500c87071c0204076c040c01400000000000000000000000000000000000000000f0
	*/
	raw := "0302d3641a00ca500c87071c0204076c040c01400000000000000000000000000000000000000000f0"

	parser, _ := NewFrameParser(raw)

	parser.Header()
	notice, err := parser.Notice()
	appConfig := notice.(ApplicationConfig)

	assert.Nil(t, err)
	assert.NotNil(t, appConfig)
	assert.Equal(t, AppModeImpact, appConfig.AppMode)
	assert.Equal(t, ConfigBMARange16G, appConfig.BMARange)
	assert.Equal(t, 1900, appConfig.BMAHighGThresholdMg)
	assert.Equal(t, 4, appConfig.NoSubInterval)
	assert.Equal(t, 12, appConfig.NRFRejectThresholdMg)
	assert.Equal(t, 320, appConfig.NRFImpactThresholdMg)
	assert.Equal(t, 0, appConfig.InclinationCheckPeriod)
	assert.Equal(t, 0, appConfig.TxSkipNo)
	assert.Equal(t, 0, appConfig.BaseX)
	assert.Equal(t, 0, appConfig.BaseY)
	assert.Equal(t, 0, appConfig.BaseZ)
	assert.Equal(t, MRMTStateCommission, appConfig.MRMTState)
}
