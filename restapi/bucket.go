package restapi

import (
	"net/http"
	"strings"
	"test/dapi/internal"
	"time"
)

type BucketLists map[string]Bucket

var Buckets = make(BucketLists, 0)

func (bucketsList BucketLists) findBucket(Route, Method string) Bucket {
	splitted := strings.Split(Route, "/")

	bucket, found := bucketsList[splitted[0]+Method]

	if !found {
		_bucket := Bucket{Route: Route, Method: Method, Requests: map[*chan BucketRequestAnswer]BucketRequest{}, Key: splitted[0] + Method}
		_bucket.Index = len(Buckets) + 1
		Buckets[Route+Method] = _bucket

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
	Url         string
	Methode     string
	Payload     []byte
	AnswerQueue chan BucketRequestAnswer
	BucketName  string
	WantAnswer  bool
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
	Handled   bool
	BucketID  string
}

func (b *Bucket) save(bl BucketLists) {
	bl[b.Key] = *b
}

func (b *Bucket) unLockBucket() {
	b.Blocked = false
	internal.LogDebug("UNLOCKING BUCKET", b.Method)
}

func (b *Bucket) lockBucket() {
	b.Blocked = true
	internal.LogDebug("LOCKING BUCKET", b.Method)
}

func (b *Bucket) handled() {
	b.Handled = true
}

// Add a discord api request into the discord handler
// should send an BucketRequest containing a valid chan BucketRequestAnswer
func addRequest(br BucketRequest) {
	_bucket := Buckets.findBucket(br.BucketName, br.Methode)
	_bucket.Requests[&br.AnswerQueue] = br
	Buckets.saveBucket(_bucket)
}
