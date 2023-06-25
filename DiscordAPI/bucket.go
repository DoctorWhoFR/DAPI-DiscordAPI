package DiscordAPI

import (
	"azginfr/dapi/DiscordInternal"
	"bytes"
	"net/http"
	"strings"
	"time"
)

type BucketLists map[string]Bucket

var GlobalDiscordBuckets = make(BucketLists, 0)

func (bucketsList BucketLists) findBucket(Route, Method string) Bucket {
	splitted := strings.Split(Route, "/")

	bucket, found := bucketsList[splitted[0]+Method]

	if !found {
		_bucket := Bucket{Route: Route, Method: Method, Requests: map[*chan BucketRequestAnswer]BucketRequest{}, Key: splitted[0] + Method}
		_bucket.Index = len(GlobalDiscordBuckets) + 1
		GlobalDiscordBuckets[Route+Method] = _bucket

		return _bucket
	}

	return bucket
}

func (bucketsList BucketLists) saveBucket(b Bucket) {
	bucketsList[b.Key] = b
}

// BucketRequestAnswer is a specif type used in a chan, to make easier to the client to get his response
// Containing :
//   - Body []byte -> response byte from discord api
//   - Res http.Response --> basic http response
type BucketRequestAnswer struct {
	Body []byte
	Res  http.Response
	Err  error
}

type BucketRequest struct {
	Url           string
	Methode       string
	Payload       []byte
	AnswerQueue   chan BucketRequestAnswer
	BucketName    string
	WantAnswer    bool
	FormHTTP      bool
	FormWriter    *bytes.Buffer
	ContentHeader string
}

type Bucket struct {
	Blocked   bool
	Method    string
	Index     int
	Key       string
	ResetTime time.Time
	Route     string
	Remaining int64
	Requests  map[*chan BucketRequestAnswer]BucketRequest
	IsHandled bool
	BucketID  string
}

func (b *Bucket) save(bl BucketLists) {
	bl[b.Key] = *b
}

func (b *Bucket) unLockBucket() {
	b.Blocked = false
	DiscordInternal.LogDebug("UNLOCKING BUCKET", b.Method)
}

func (b *Bucket) lockBucket() {
	b.Blocked = true
	DiscordInternal.LogDebug("LOCKING BUCKET", b.Method)
}

func (b *Bucket) handled() {
	b.IsHandled = true
}

// Add a discord api request into the discord handler
// should send an BucketRequest containing a valid chan BucketRequestAnswer
func addRequest(br BucketRequest) {
	_bucket := GlobalDiscordBuckets.findBucket(br.BucketName, br.Methode)
	_bucket.Requests[&br.AnswerQueue] = br
	GlobalDiscordBuckets.saveBucket(_bucket)
}

/*
HandleBucketRequest

Thread routine, launched by NewBucketHandler.

Used to handle new discord request, send by library call.

Every Bucket have its own thread.
*/
func HandleBucketRequest(b Bucket) {
	for {
		for c, v := range b.Requests {
			start := time.Now()
			if v.FormHTTP {
				body, res, err := httpDiscordCallFormData(v.Url, v.Methode, v.WantAnswer, v.FormWriter, v.ContentHeader)
				answer := BucketRequestAnswer{Body: body, Res: **res, Err: err}
				v.AnswerQueue <- answer
				delete(b.Requests, c)

				DiscordInternal.HandlingTimeBot = append(DiscordInternal.HandlingTimeBot, time.Since(start))
			} else {
				body, res, err := httpDiscordCallJson(v.Url, v.Methode, v.Payload, v.WantAnswer)
				answer := BucketRequestAnswer{Body: body, Res: **res, Err: err}
				v.AnswerQueue <- answer
				delete(b.Requests, c)

				DiscordInternal.HandlingTimeBot = append(DiscordInternal.HandlingTimeBot, time.Since(start))
			}

			DiscordInternal.LogInfo("handled total request and answered in", time.Since(start).Milliseconds(), "ms, medium: ", DiscordInternal.MediumValueBOT(), "ms")
			DiscordInternal.LogInfo("handled without discord medium:", DiscordInternal.MediumValueWithoutAPI(), "ms")
		}
		DiscordInternal.SimpleSleep()
	}

}
