package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

type S3Redundancy string
type S3ServerSideEncryption string

const (
	S3RedundancyStandard      S3Redundancy           = "standard"
	S3RedundancyReduced       S3Redundancy           = "reduced_redundancy"
	S3ServerSideEncryptionAES S3ServerSideEncryption = "AES256"
	S3ServerSideEncryptionKMS S3ServerSideEncryption = "aws:kms"
)

// S3 represents a S3 response from the Fastly API.
type S3 struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name                         string                 `mapstructure:"name"`
	BucketName                   string                 `mapstructure:"bucket_name"`
	Domain                       string                 `mapstructure:"domain"`
	AccessKey                    string                 `mapstructure:"access_key"`
	SecretKey                    string                 `mapstructure:"secret_key"`
	Path                         string                 `mapstructure:"path"`
	Period                       uint                   `mapstructure:"period"`
	GzipLevel                    uint                   `mapstructure:"gzip_level"`
	Format                       string                 `mapstructure:"format"`
	FormatVersion                uint                   `mapstructure:"format_version"`
	ResponseCondition            string                 `mapstructure:"response_condition"`
	MessageType                  string                 `mapstructure:"message_type"`
	TimestampFormat              string                 `mapstructure:"timestamp_format"`
	Placement                    string                 `mapstructure:"placement"`
	PublicKey                    string                 `mapstructure:"public_key"`
	Redundancy                   S3Redundancy           `mapstructure:"redundancy"`
	ServerSideEncryptionKMSKeyID string                 `mapstructure:"server_side_encryption_kms_key_id"`
	ServerSideEncryption         S3ServerSideEncryption `mapstructure:"server_side_encryption"`
	CreatedAt                    *time.Time             `mapstructure:"created_at"`
	UpdatedAt                    *time.Time             `mapstructure:"updated_at"`
	DeletedAt                    *time.Time             `mapstructure:"deleted_at"`
}

// s3sByName is a sortable list of S3s.
type s3sByName []*S3

// Len, Swap, and Less implement the sortable interface.
func (s s3sByName) Len() int      { return len(s) }
func (s s3sByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s s3sByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListS3sInput is used as input to the ListS3s function.
type ListS3sInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListS3s returns the list of S3s for the configuration version.
func (c *Client) ListS3s(i *ListS3sInput) ([]*S3, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s3s []*S3
	if err := decodeBodyMap(resp.Body, &s3s); err != nil {
		return nil, err
	}
	sort.Stable(s3sByName(s3s))
	return s3s, nil
}

// CreateS3Input is used as input to the CreateS3 function.
type CreateS3Input struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name                         string                 `form:"name,omitempty"`
	BucketName                   string                 `form:"bucket_name,omitempty"`
	Domain                       string                 `form:"domain,omitempty"`
	AccessKey                    string                 `form:"access_key,omitempty"`
	SecretKey                    string                 `form:"secret_key,omitempty"`
	Path                         string                 `form:"path,omitempty"`
	Period                       uint                   `form:"period,omitempty"`
	GzipLevel                    uint                   `form:"gzip_level,omitempty"`
	Format                       string                 `form:"format,omitempty"`
	MessageType                  string                 `form:"message_type,omitempty"`
	FormatVersion                uint                   `form:"format_version,omitempty"`
	ResponseCondition            string                 `form:"response_condition,omitempty"`
	TimestampFormat              string                 `form:"timestamp_format,omitempty"`
	Redundancy                   S3Redundancy           `form:"redundancy,omitempty"`
	Placement                    string                 `form:"placement,omitempty"`
	PublicKey                    string                 `form:"public_key,omitempty"`
	ServerSideEncryptionKMSKeyID string                 `form:"server_side_encryption_kms_key_id,omitempty"`
	ServerSideEncryption         S3ServerSideEncryption `form:"server_side_encryption,omitempty"`
}

// CreateS3 creates a new Fastly S3.
func (c *Client) CreateS3(i *CreateS3Input) (*S3, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.ServerSideEncryption == S3ServerSideEncryptionKMS && i.ServerSideEncryptionKMSKeyID == "" {
		return nil, ErrMissingKMSKeyID
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// GetS3Input is used as input to the GetS3 function.
type GetS3Input struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the S3 to fetch.
	Name string
}

// GetS3 gets the S3 configuration with the given parameters.
func (c *Client) GetS3(i *GetS3Input) (*S3, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// UpdateS3Input is used as input to the UpdateS3 function.
type UpdateS3Input struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the S3 to update.
	Name string

	NewName                      string                 `form:"name,omitempty"`
	BucketName                   string                 `form:"bucket_name,omitempty"`
	Domain                       string                 `form:"domain,omitempty"`
	AccessKey                    string                 `form:"access_key,omitempty"`
	SecretKey                    string                 `form:"secret_key,omitempty"`
	Path                         string                 `form:"path,omitempty"`
	Period                       uint                   `form:"period,omitempty"`
	GzipLevel                    uint                   `form:"gzip_level,omitempty"`
	Format                       string                 `form:"format,omitempty"`
	FormatVersion                uint                   `form:"format_version,omitempty"`
	ResponseCondition            string                 `form:"response_condition,omitempty"`
	MessageType                  string                 `form:"message_type,omitempty"`
	TimestampFormat              string                 `form:"timestamp_format,omitempty"`
	Redundancy                   S3Redundancy           `form:"redundancy,omitempty"`
	Placement                    string                 `form:"placement,omitempty"`
	PublicKey                    string                 `form:"public_key,omitempty"`
	ServerSideEncryptionKMSKeyID string                 `form:"server_side_encryption_kms_key_id,omitempty"`
	ServerSideEncryption         S3ServerSideEncryption `form:"server_side_encryption,omitempty"`
}

// UpdateS3 updates a specific S3.
func (c *Client) UpdateS3(i *UpdateS3Input) (*S3, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.ServerSideEncryption == S3ServerSideEncryptionKMS && i.ServerSideEncryptionKMSKeyID == "" {
		return nil, ErrMissingKMSKeyID
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// DeleteS3Input is the input parameter to DeleteS3.
type DeleteS3Input struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the S3 to delete (required).
	Name string
}

// DeleteS3 deletes the given S3 version.
func (c *Client) DeleteS3(i *DeleteS3Input) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
