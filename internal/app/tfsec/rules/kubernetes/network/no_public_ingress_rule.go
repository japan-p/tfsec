package network

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

import (
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/cidr"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/hclcontext"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/provider"
	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/rule"
	"github.com/aquasecurity/tfsec/pkg/severity"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Provider:  provider.KubernetesProvider,
		Service:   "network",
		ShortCode: "no-public-ingress",
		Documentation: rule.RuleDocumentation{
			Summary:     "Public ingress should not be allowed via network policies",
			Explanation: `You should not expose infrastructure to the public internet except where explicitly required`,
			Impact:      "Exposure of infrastructure to the public internet",
			Resolution:  "Remove public access except where explicitly required",
			BadExample: []string{`
resource "kubernetes_network_policy" "bad_example" {
  metadata {
    name      = "terraform-example-network-policy"
    namespace = "default"
  }

  spec {
    pod_selector {
      match_expressions {
        key      = "name"
        operator = "In"
        values   = ["webfront", "api"]
      }
    }

    ingress {
      ports {
        port     = "http"
        protocol = "TCP"
      }
      ports {
        port     = "8125"
        protocol = "UDP"
      }

      from {
        ip_block {
          cidr = "0.0.0.0/0"
          except = [
            "10.0.0.0/24",
            "10.0.1.0/24",
          ]
        }
      }
    }

    egress {
      ports {
        port     = "http"
        protocol = "TCP"
      }
      ports {
        port     = "8125"
        protocol = "UDP"
      }

      to {
        ip_block {
          cidr = "0.0.0.0/0"
          except = [
            "10.0.0.0/24",
            "10.0.1.0/24",
          ]
        }
      }
    }

    policy_types = ["Ingress", "Egress"]
  }
}
`},
			GoodExample: []string{`
resource "kubernetes_network_policy" "good_example" {
  metadata {
    name      = "terraform-example-network-policy"
    namespace = "default"
  }

  spec {
    pod_selector {
      match_expressions {
        key      = "name"
        operator = "In"
        values   = ["webfront", "api"]
      }
    }

    ingress {
      ports {
        port     = "http"
        protocol = "TCP"
      }
      ports {
        port     = "8125"
        protocol = "UDP"
      }

      from {
        ip_block {
          cidr = "10.0.0.0/16"
          except = [
            "10.0.0.0/24",
            "10.0.1.0/24",
          ]
        }
      }
    }

    egress {
      ports {
        port     = "http"
        protocol = "TCP"
      }
      ports {
        port     = "8125"
        protocol = "UDP"
      }

      to {
        ip_block {
          cidr = "0.0.0.0/0"
          except = [
            "10.0.0.0/24",
            "10.0.1.0/24",
          ]
        }
      }
    }

    policy_types = ["Ingress", "Egress"]
  }
}
`},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/network_policy#spec.ingress.from.ip_block.cidr",
			},
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"kubernetes_network_policy",
		},
		DefaultSeverity: severity.High,
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context) {
			ingressBlock := resourceBlock.GetBlock("spec").GetBlock("ingress")
			if ingressBlock.IsNil() || len(ingressBlock.GetBlocks("from")) == 0 {
				set.AddResult().
					WithDescription("Resource '%s' allows all ingress traffic by default", resourceBlock.FullName())
				return
			}

			for _, from := range ingressBlock.GetBlocks("from") {
				if cidrAttr := from.GetBlock("ip_block").GetAttribute("cidr"); cidrAttr.IsString() {
					if cidr.IsOpen(cidrAttr) {
						set.AddResult().
							WithDescription("Resource '%s' allows ingress traffic from the internet", resourceBlock.FullName()).
							WithAttribute(cidrAttr)
					}
				}
			}
		},
	})
}