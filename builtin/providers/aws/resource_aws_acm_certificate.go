package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsAcmCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAcmCertificateCreate,
		Read:   resourceAwsAcmCertificateRead,
		Update: resourceAwsAcmCertificateUpdate,
		Delete: resourceAwsAcmCertificateDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subject_alternative_names": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"route53_zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsAcmCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsAcmCertificateRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsAcmCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsAcmCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// acmHelperValidateDomainRecords validates the DNS records for the domains
// that are being created, or rather, that they do not have any TXT or MX
// records, to ensure that Route53 can add them without affecting regular
// operation.
func acmHelperValidateDomainRecords(meta interface{}) error {
	return nil
}

// acmHelperCreateS3Bucket creates a bucket that we will use to direct the
// certificate requests to, via SES.
func acmHelperCreateS3Bucket(meta interface{}) (string, error) {
	bucket := resource.PrefixedUniqueId("tf-acm-bucket")
	conn := meta.(*AWSClient).s3conn

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("private"),
	}
	_, err := conn.CreateBucket(params)

	if err != nil {
		return "", err
	}

	return bucket, nil
}

// acmHelperCreateSESDomain creates the SES domain, adds the MX records to the
// route 53 hosted zone, and adds the rule sets for admin@ for each domain,
// directing the alias to the bucket.
func acmHelperCreateSESDomain(r53id string, domains []string, meta interface{}) error {
	tokens, err := acmHelperVerifySESDoamins(domains, meta)
	if err != nil {
		return err
	}

	err = acmHelperSESAddTXT(tokens, r53id, meta)
	if err != nil {
		return err
	}

	// Poll SES until all domains are verified.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Success"},
		Refresh:    acmHelperSESPollFunc(domains, meta),
		Timeout:    5 * time.Minute,
		MinTimeout: 15 * time.Second,
		Delay:      10 * time.Minute,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return err
	}

	err = acmHelperSESAddMX(domains, r53id, meta)
	if err != nil {
		return err
	}

	return nil
}

// acmHelperVerifySESDoamins creates the SES domains and gets the tokens.
func acmHelperVerifySESDoamins(domains []string, meta interface{}) (map[string]string, error) {
	sesconn := meta.(*AWSClient).sesconn
	tokens := make(map[string]string)
	for _, v := range domains {
		params := &ses.VerifyDomainIdentityInput{
			Domain: aws.String(v),
		}
		resp, err := sesconn.VerifyDomainIdentity(params)

		if err != nil {
			return nil, err
		}

		tokens[v] = *resp.VerificationToken
	}
	return tokens, nil
}

// a small helper type for acmHelperSESRoute53Helper.
type acmSESRRData struct {
	Type     string
	Value    string
	Priority string
}

// acmHelperSESRoute53Helper performs a Route53 operation.
func acmHelperSESRoute53Helper(action string, records map[string]acmSESRRData, r53id string, meta interface{}) error {
	r53conn := meta.(*AWSClient).r53conn

	// add the Route 53 record sets for verification.
	changes := []*route53.Change{}
	for k, v := range records {
		changes = append(changes, &route53.Change{
			Action: aws.String(action),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(k),
				Type:            aws.String(v.Type),
				ResourceRecords: []*route53.ResourceRecord{{Value: aws.String(fmt.Sprintf("%s%s", v.Priority, v.Value))}},
				TTL:             aws.Int64(5),
			},
		})
	}

	changeParams := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
			Comment: aws.String("Created by Terraform for ACM verification"),
		},
		HostedZoneId: aws.String(r53id),
	}
	resp, err := r53conn.ChangeResourceRecordSets(changeParams)

	if err != nil {
		return err
	}

	changeid := *resp.ChangeInfo.Id

	waitParams := &route53.GetChangeInput{
		Id: aws.String(changeid),
	}

	err = r53conn.WaitUntilResourceRecordSetsChanged(waitParams)

	if err != nil {
		return err
	}

	return nil
}

// acmHelperSESAddTXT adds the Route 53 record sets for SES verification.
func acmHelperSESAddTXT(tokens map[string]string, r53id string, meta interface{}) error {
	records := make(map[string]acmSESRRData)
	for k, v := range tokens {
		records[k] = acmSESRRData{
			Type:  "TXT",
			Value: v,
		}
	}
	err := acmHelperSESRoute53Helper("CREATE", records, r53id, meta)

	if err != nil {
		return err
	}

	return nil
}

// acmHelperSESAddMX adds the Route 53 MX record sets for ACM validation.
func acmHelperSESAddMX(domains []string, r53id string, meta interface{}) error {
	region := meta.(*AWSClient).region
	mxhost := fmt.Sprintf("inbound-smtp.%s.amazonaws.com.", region)

	records := make(map[string]acmSESRRData)
	for _, v := range domains {
		records[v] = acmSESRRData{
			Type:     "MX",
			Value:    mxhost,
			Priority: "10 ",
		}
	}
	err := acmHelperSESRoute53Helper("CREATE", records, r53id, meta)

	if err != nil {
		return err
	}

	return nil
}

// acmHelperSESPollFunc polls SES until all domains have returned as verified.
func acmHelperSESPollFunc(domains []string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn := meta.(*AWSClient).sesconn
		params := &ses.GetIdentityVerificationAttributesInput{
			Identities: aws.StringSlice(domains),
		}

		resp, err := conn.GetIdentityVerificationAttributes(params)
		if err != nil {
			return nil, "Failed", fmt.Errorf("Error on fetching SES verification status for domains: %s", err)
		}

		for _, v := range resp.VerificationAttributes {
			switch *v.VerificationStatus {
			case "Failed":
				return nil, "Failed", fmt.Errorf("Verifying domains through SES failed")
			case "Success":
				return nil, "Success", nil
			default:
				return nil, "Pending", nil
			}
		}
		return nil, "Failed", fmt.Errorf("Unhandled case in SES waiter")
	}
}

func acmHelperRequestCertificate() error {
	return nil
}

func acmHelperPollS3Bucket() error {
	return nil
}

func acmHelperApproveRequest() error {
	return nil
}

func acmHelperCleanup() error {
	return nil
}

func acmHelperDeleteSESDomain() error {
	return nil
}

func acmHelperDeleteS3Bucket() error {
	return nil
}
