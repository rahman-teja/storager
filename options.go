package storager

import "time"

type GetOptions struct {
	VersionID string
}

type PutOptions struct {
	ContentType             string
	UserMetadata            map[string]string
	UserTags                map[string]string
	ContentEncoding         string
	ContentDisposition      string
	ContentLanguage         string
	CacheControl            string
	WebsiteRedirectLocation string
	Expires                 time.Time
}

type RemoveOptions struct {
	VersionID string
}
