package nobl9

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	n9api "github.com/nobl9/nobl9-go"
)

func TestAcc_Nobl9SLO(t *testing.T) {
	cases := []struct {
		name       string
		configFunc func(string) string
	}{
		{"test-prometheus", testPrometheusSLO},
		{"test-prom-with-ap", testPrometheusSLOWithAlerPolicy},
		{"test-prom-with-countmetrics", testPrometheusSLOWithCountMetrics},
		{"test-prom-with-multiple-objectives", testPrometheusSLOWithMultipleObjectives},
		{"test-newrelic", testNewRelicSLO},
		{"test-appdynamics", testAppdynamicsSLO},
		{"test-splunk", testSplunkSLO},
		{"test-lightstep", testLightstepSLO},
		{"test-splunk-observability", testSplunkObservabilitySLO},
		{"test-dynatrace", testDynatraceSLO},
		{"test-elasticsearch", testElasticsearchSLO},
		{"test-thousandeyes", testThousandeyesSLO},
		{"test-graphite", testGraphiteSLO},
		{"test-bigquery", testBigQuerySLO},
		{"test-opentsdb", testOpenTSDBSLO},
		{"test-grafana-loki", testGrafanaLokiSLO},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: ProviderFactory(),
				CheckDestroy:      DestroyFunc("nobl9_slo", n9api.ObjectSLO),
				Steps: []resource.TestStep{
					{
						Config: tc.configFunc(tc.name),
						Check:  CheckObjectCreated("nobl9_slo." + tc.name),
					},
				},
			})
		})
	}
}

func testPrometheusSLO(name string) string {
	config := testService(name+"-service") +
		testPrometheusConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      prometheus {
        promql = "1.0"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testPrometheusSLOWithAlerPolicy(name string) string {
	config := testService(name+"-service") +
		testPrometheusConfig(name+"-agent") +
		testAlertPolicyWithoutIntegration(name+"-ap") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      prometheus {
        promql = "1.0"
      }
    }
  }

  alert_policies = [ nobl9_alert_policy.:name-ap.name ]
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testPrometheusSLOWithCountMetrics(name string) string {
	config := testService(name+"-service") +
		testPrometheusConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
	count_metrics {
	  incremental = true
	  good {
		prometheus {
		  promql = "1.0"
		}
	  }
	  total {
		prometheus {
		  promql = "1.0"
		}
	  }
	}
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}
func testPrometheusSLOWithMultipleObjectives(name string) string {
	config := testService(name+"-service") +
		testPrometheusConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  objective {
    display_name = "obj2"
    target       = 0.5
    value        = 10
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      prometheus {
        promql = "1.0"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testDatadogSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      datadog {
        query = "avg:system.cpu.user{cluster_name:main}"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testNewRelicSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      newrelic {
        nrql = "SELECT average(duration * 1000) FROM Transaction TIMESERIES"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testAppdynamicsSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      appdynamics {
		application_name = "polakpotrafi"
        metric_path = "End User Experience|App|End User Response Time 95th percentile (ms)"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testSplunkSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      splunk {
        query = "TODO"
		field_name = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testLightstepSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      lightstep {
        stream_id = DzpxcSRh
		type_of_data = latency
		percentile = 95
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testSplunkObservabilitySLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      splunk_observability {
        program = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testDynatraceSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      dynatrace {
        metric_selector = "builtin:synthetic.http.duration.geo:filter(and(in("dt.entity.http_check",entitySelector("type(http_check),entityName(~"API Sample~")")),in("dt.entity.synthetic_location",entitySelector("type(synthetic_location),entityName(~"N. California~")")))):splitBy("dt.entity.http_check","dt.entity.synthetic_location"):avg:auto:sort(value(avg,descending)):limit(20)"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testThousandeyesSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      thousandeyes {
        test_id = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testGraphiteSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      graphite {
        metric_path = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testBigQuerySLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      bigquery {
        projectId = "bdwtest-256112"
        location = "EU"
        query = "SELECT response_time AS n9value, created AS n9date"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testOpenTSDBSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      opentsdb {
        query = "m=none:{{.N9RESOLUTION}}-avg-zero:cpu{cpu.usage=core.1}"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testGrafanaLokiSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      grafana_loki {
        logql = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}

func testElasticsearchSLO(name string) string {
	config := testService(name+"-service") +
		testDatadogConfig(name+"-agent") + `
resource "nobl9_slo" ":name" {
  name         = ":name"
  display_name = ":name"
  project      = "terraform"
  service      = nobl9_service.:name-service.name

  budgeting_method = "Occurrences"

  objective {
    display_name = "obj1"
    target       = 0.7
    value        = 1
    op           = "lt"
  }

  time_window {
    count      = 10
    is_rolling = true
    unit       = "Minute"
  }

  indicator {
    name    = nobl9_agent.:name-agent.name
    project = ":project"
	kind    = "Agent"
    raw_metric {
      elasticsearch {
		index = "TODO"
        query = "TODO"
      }
    }
  }
}
`
	config = strings.ReplaceAll(config, ":name", name)
	config = strings.ReplaceAll(config, ":project", testProject)

	return config
}