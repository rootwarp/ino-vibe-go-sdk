package parser

var frameSpec = `
{
    "header": {
        "desc": "Ino-Vibe LoRa Header frame version 3",
        "fields": [
            {
                "name": "version",
                "bits": 8,
                "signed": false
            },
            {
                "name": "dev_type",
                "bits": 8,
                "signed": false
            },
            {
                "name": "seq",
                "bits": 8,
                "signed": false
            },
            {
                "name": "battery",
                "bits": 8, 
                "signed": false
            },
            {
                "name": "temperature",
                "bits": 8,
                "signed": true
            },
            {
                "name": "lora_err",
                "bits": 8,
                "signed": false
            },
            {
                "name": "rssi",
                "bits": 8,
                "signed": true
            },
            {
                "name": "payload",
                "bits": 8,
                "fields":  [
                    {
                        "name": "type",
                        "bits": 4,
                        "signed": false
                    },
                    {
                        "name": "request",
                        "bits": 4,
                        "signed": false
                    }
                ]
            },
            {
                "name": "resv",
                "bits": 16,
                "signed": false
            }
        ]
    }, 
    "alive": {
        "desc": "Alive Frame",
        "fields": [
            {
                "name": "x",
                "bits": 16,
                "signed": true
            },
            {
                "name": "y",
                "bits": 16,
                "signed": true
            },
            {
                "name": "z",
                "bits": 16,
                "signed": true
            },
            {
                "name": "alive_period",
                "bits": 16,
                "signed": false
            },
            {
                "name": "sensitivity",
                "bits": 8,
                "signed": false
            },
            {
                "name": "threshold",
                "bits": 16,
                "signed": false
            },
            {
                "name": "acc_int_no",
                "bits": 8,
                "signed": false
            },
            {
                "name": "acc_int_resv",
                "bits": 8,
                "signed": false
            },
            {
                "name": "acc_int_data",
                "bits": 8,
                "signed": false
            },
            {
                "name": "log_enable",
                "bits": 8,
                "signed": false
            },
            {
                "name": "log_interval",
                "bits": 8,
                "signed": false
            },
            {
                "name": "log_blocks",
                "bits": 8,
                "signed": false
            },
            {
                "name": "setup",
                "bits": 8,
                "signed": false
            },
            {
                "name": "app_fw_major",
                "bits": 8,
                "signed": false
            },
            {
                "name": "app_fw_minor",
                "bits": 8,
                "signed": false
            },
            {
                "name": "app_fw_rev",
                "bits": 8,
                "signed": false
            },
            {
                "name": "lora_fw_major",
                "bits": 8,
                "signed": false
            },
            {
                "name": "lora_fw_minor",
                "bits": 8,
                "signed": false
            },
            {
                "name": "lora_fw_rev",
                "bits": 8,
                "signed": false
            }
        ]
    },
    "wave": {
        "desc": "Wave raw data",
        "fields": [
            {
                "name": "control",
                "bits": 8,
                "fields":  [
                    {
                        "name": "bma_range",
                        "bits": 2,
                        "signed": false
                    },
                    {
                        "name": "axis",
                        "bits": 2,
                        "signed": false
                    },
                    {
                        "name": "id",
                        "bits": 4,
                        "signed": false
                    }
                ]
            },
            {
                "name": "pack_type",
                "bits": 4,
                "signed": false
            },
            {
                "name": "pos",
                "bits": 4,
                "signed": false
            },
            {
                "name": "x",
                "bits": 16,
                "signed": true
            },
            {
                "name": "y",
                "bits": 16,
                "signed": true
            },
            {
                "name": "z",
                "bits": 16,
                "signed": true
            }
        ]
    },
    "notice": {
        "desc": "Notice from device",
        "fields": [
            {
                "name": "type",
                "bits": 8,
                "signed": false
            },
            {
                "name": "length",
                "bits": 8,
                "signed": false
            },
            {
                "name": "payload",
                "subtype": {
                    "1": "power_up",
                    "4": "setup",
                    "6": "reject_count",
                    "7": "notice_config"
                }
            }
        ],
        "subtypes": {
            "power_up": {
                "desc": "PowerUp Notice",
                "fields": [
                    {
                        "name": "reset_reason",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "turnon_count",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "resv",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "off_reason",
                        "bits": 8,
                        "signed": false
                    }
                ]
            },
            "setup": {
                "desc": "Setup",
                "fields": [
                    {
                        "name": "current_state",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "previous_state",
                        "bits": 8,
                        "signed": false
                    }
                ]
            },
            "reject_count": {
                "desc": "Rejection Count",
                "fields": [
                    {
                        "name": "count",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "period",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "threshold",
                        "bits": 8,
                        "signed": false
                    }
                ]
            },
            "notice_config": {
                "desc": "Application config",
                "fields": [
                    {
                        "name": "app_mode",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "bma_g_range",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "bma_high_g_threshold_mg",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "no_sub_interval",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "nrf_reject_threshold_mg",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "nrf_impact_threshold_mg",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "inclination_check_period",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "tx_skip_no",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "base_x",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "base_y",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "base_z",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "resv0",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "mrmt_state",
                        "bits": 8,
                        "signed": false
                    },
                    {
                        "name": "mrmt_operation_threshold_mg",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "mrmt_shock_threshold_mg",
                        "bits": 16,
                        "signed": false
                    },
                    {
                        "name": "resv1",
                        "bits": 48,
                        "signed": false
                    }
                ]
            }
        }
    }
}
`
