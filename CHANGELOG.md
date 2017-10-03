## 1.0.0 (Unreleased)

IMPROVEMENTS:

* resource/dyn_record: Add support for `NS` & `MX` records [GH-15]

BUG FIXES:

* resource/dyn_record: Avoid diff for default TTL [GH-12]
* resource/dyn_record: Ignore trailing dot in FQDN [GH-13]
* resource/dyn_record: Support records for top level domain [GH-14]
* resource/dyn_record: Fix broken record update [GH-17]

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
