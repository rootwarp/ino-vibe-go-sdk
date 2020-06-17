package parser

// TODO: Add resson enums.

// PowerUp is power up frame.
type PowerUp struct {
	Type        NoticePayloadType `json:"type"`
	Length      int               `json:"length"`
	ResetReason int               `json:"reset_reason"`
	Count       int               `json:"turnon_count"`
	OffReason   int               `json:"off_reason"`
}

// Setup describes setup frames.
type Setup struct {
	Type     NoticePayloadType `json:"type"`
	Length   int               `json:"length"`
	Previous DeviceSetup       `json:"previous_state"`
	Current  DeviceSetup       `json:"current_state"`
}

// RejectCount describes rejection count frame.
type RejectCount struct {
	Type      NoticePayloadType `json:"type"`
	Length    int               `json:"length"`
	Threshold int               `json:"threshold"`
	Period    int               `json:"period"`
	Count     int               `json:"count"`
}

// ApplicationMode is detection mode of device.
type ApplicationMode int

// ApplicationModes
const (
	AppModeExaInc ApplicationMode = 0
	AppModeMRMT   ApplicationMode = 1
	AppModeImpact ApplicationMode = 2
)

// ApplicationConfig descibes config frame about working application.
type ApplicationConfig struct {
	Type                     int             `json:"type"`
	Length                   int             `json:"length"`
	AppMode                  ApplicationMode `json:"app_mode"`
	BMARange                 ConfigBMARange  `json:"bma_g_range"`
	BMAHighGThresholdMg      int             `json:"bma_high_g_threshold_mg"`
	NoSubInterval            int             `json:"no_sub_interval"`
	NRFRejectThresholdMg     int             `json:"nrf_reject_threshold_mg"`
	NRFImpactThresholdMg     int             `json:"nrf_impact_threshold_mg"`
	InclinationCheckPeriod   int             `json:"inclination_check_period"`
	TxSkipNo                 int             `json:"tx_skip_no"`
	BaseX                    int             `json:"base_x"`
	BaseY                    int             `json:"base_y"`
	BaseZ                    int             `json:"base_z"`
	MRMTState                MRMTState       `json:"mrmt_state"`
	MRMTOperationThresholdMg int             `json:"mrmt_operation_threshold_mg"`
	MRMTShockThresholdMg     int             `json:"mrmt_shock_threshold_mg"`
}
