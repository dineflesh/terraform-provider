// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func TestAccIBMIamTrustedProfileLinkBasic(t *testing.T) {
	var conf iamidentityv1.ProfileLink
	profileID := fmt.Sprintf("tf_profile_id_%d", acctest.RandIntRange(10, 100))
	crType := fmt.Sprintf("tf_cr_type_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIamTrustedProfileLinkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIamTrustedProfileLinkConfigBasic(profileID, crType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIamTrustedProfileLinkExists("ibm_iam_trusted_profile_link.iam_trusted_profile_link", conf),
					resource.TestCheckResourceAttr("ibm_iam_trusted_profile_link.iam_trusted_profile_link", "profile_id", profileID),
					resource.TestCheckResourceAttr("ibm_iam_trusted_profile_link.iam_trusted_profile_link", "cr_type", crType),
				),
			},
		},
	})
}

func TestAccIBMIamTrustedProfileLinkAllArgs(t *testing.T) {
	var conf iamidentityv1.ProfileLink
	profileID := fmt.Sprintf("tf_profile_id_%d", acctest.RandIntRange(10, 100))
	crType := fmt.Sprintf("tf_cr_type_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIamTrustedProfileLinkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIamTrustedProfileLinkConfig(profileID, crType, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIamTrustedProfileLinkExists("ibm_iam_trusted_profile_link.iam_trusted_profile_link", conf),
					resource.TestCheckResourceAttr("ibm_iam_trusted_profile_link.iam_trusted_profile_link", "profile_id", profileID),
					resource.TestCheckResourceAttr("ibm_iam_trusted_profile_link.iam_trusted_profile_link", "cr_type", crType),
					resource.TestCheckResourceAttr("ibm_iam_trusted_profile_link.iam_trusted_profile_link", "name", name),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_iam_trusted_profile_link.iam_trusted_profile_link",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIamTrustedProfileLinkConfigBasic(profileID string, crType string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_trusted_profile_link" "iam_trusted_profile_link" {
			profile_id = "%s"
			cr_type = "%s"
			link {
				crn = "crn"
				namespace = "namespace"
				name = "name"
			}
		}
	`, profileID, crType)
}

func testAccCheckIBMIamTrustedProfileLinkConfig(profileID string, crType string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_trusted_profile_link" "iam_trusted_profile_link" {
			profile_id = "%s"
			cr_type = "%s"
			link {
				crn = "crn"
				namespace = "namespace"
				name = "name"
			}
			name = "%s"
		}
	`, profileID, crType, name)
}

func testAccCheckIBMIamTrustedProfileLinkExists(n string, obj iamidentityv1.ProfileLink) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		iamIdentityClient, err := testAccProvider.Meta().(ClientSession).IAMIdentityV1API()
		if err != nil {
			return err
		}

		getLinkOptions := &iamidentityv1.GetLinkOptions{}

		parts, err := sepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getLinkOptions.SetProfileID(parts[0])
		getLinkOptions.SetLinkID(parts[1])

		profileLink, _, err := iamIdentityClient.GetLink(getLinkOptions)
		if err != nil {
			return err
		}

		obj = *profileLink
		return nil
	}
}

func testAccCheckIBMIamTrustedProfileLinkDestroy(s *terraform.State) error {
	iamIdentityClient, err := testAccProvider.Meta().(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_iam_trusted_profile_link" {
			continue
		}

		getLinkOptions := &iamidentityv1.GetLinkOptions{}

		parts, err := sepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getLinkOptions.SetProfileID(parts[0])
		getLinkOptions.SetLinkID(parts[1])

		// Try to find the key
		_, response, err := iamIdentityClient.GetLink(getLinkOptions)

		if err == nil {
			return fmt.Errorf("iam_trusted_profile_link still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for iam_trusted_profile_link (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
