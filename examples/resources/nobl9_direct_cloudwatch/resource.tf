resource "nobl9_direct_cloudwatch" "test-cloudwatch" {
  name                   = "test-cloudwatch"
  project                = "terraform"
  description            = "desc"
  role_arn               = "secret"
  log_collection_enabled = true
  historical_data_retrieval {
    default_duration {
      unit  = "Day"
      value = 0
    }
    max_duration {
      unit  = "Day"
      value = 15
    }
    triggered_by_slo_creation {
      unit  = "Day"
      value = 10
    }
    triggered_by_slo_edit {
      unit  = "Day"
      value = 10
    }
  }
}