## 1.1.1 (Unreleased)
## 1.1.0 (October 23, 2017)

IMPROVEMENTS:

* resource/dyn_record: Allow importing records ([#19](https://github.com/terraform-providers/terraform-provider-dyn/issues/19))

## 1.0.0 (October 03, 2017)

IMPROVEMENTS:

* resource/dyn_record: Add support for `NS` & `MX` records ([#15](https://github.com/terraform-providers/terraform-provider-dyn/issues/15))

BUG FIXES:

* resource/dyn_record: Avoid diff for default TTL ([#12](https://github.com/terraform-providers/terraform-provider-dyn/issues/12))
* resource/dyn_record: Ignore trailing dot in FQDN ([#13](https://github.com/terraform-providers/terraform-provider-dyn/issues/13))
* resource/dyn_record: Support records for top level domain ([#14](https://github.com/terraform-providers/terraform-provider-dyn/issues/14))
* resource/dyn_record: Fix broken record update ([#17](https://github.com/terraform-providers/terraform-provider-dyn/issues/17))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
