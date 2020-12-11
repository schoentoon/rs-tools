package runemetrics

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteActivities(t *testing.T) {
	expected := `{"date":"01-Jan-1970 01:00","details":"I killed 4 Helwyrs, limb tearing hunters.","text":"I killed 4 Helwyrs."}
{"date":"01-Jan-1970 01:00","details":"I killed 3 Gregorovics, all blade wielding terrors.","text":"I killed 3 Gregorovics."}
`
	zero := time.Unix(0, 0)
	in := []Activity{
		{Date: ActivityTimeFormat{Time: &zero}, Details: "I killed 4 Helwyrs, limb tearing hunters.", Text: "I killed 4 Helwyrs."},
		{Date: ActivityTimeFormat{Time: &zero}, Details: "I killed 3 Gregorovics, all blade wielding terrors.", Text: "I killed 3 Gregorovics."},
	}

	var buf bytes.Buffer

	err := WriteActivities(&buf, in)
	assert.Nil(t, err)

	assert.Equal(t, expected, buf.String())
}

func TestReadActivities(t *testing.T) {
	in := bytes.NewBufferString(`{"date":"01-Jan-1970 01:00","details":"I killed 4 Helwyrs, limb tearing hunters.","text":"I killed 4 Helwyrs."}
{"date":"01-Jan-1970 01:00","details":"I killed 3 Gregorovics, all blade wielding terrors.","text":"I killed 3 Gregorovics."}
`)

	out, err := ReadActivities(in)
	assert.Nil(t, err)

	assert.Len(t, out, 2)
}

type testActivityIterator struct {
	errorAfter int
	count      int
	err        error
}

func (t *testActivityIterator) HandleActivity(activity Activity) error {
	t.count++

	if t.errorAfter > 0 && t.count > t.errorAfter {
		return t.err
	}
	return nil
}

func TestIterateActivities(t *testing.T) {
	test := testActivityIterator{
		errorAfter: 0,
	}

	in := bytes.NewBufferString(`{"date":"01-Jan-1970 01:00","details":"I killed 4 Helwyrs, limb tearing hunters.","text":"I killed 4 Helwyrs."}
{"date":"01-Jan-1970 01:00","details":"I killed 3 Gregorovics, all blade wielding terrors.","text":"I killed 3 Gregorovics."}
`)

	err := IterateActivities(in, &test)

	assert.Nil(t, err)
	assert.Equal(t, 2, test.count)
}

func TestIterateActivitiesError(t *testing.T) {
	test := testActivityIterator{
		errorAfter: 1,
		err:        errors.New("ERROR!"),
	}

	in := bytes.NewBufferString(`{"date":"01-Jan-1970 01:00","details":"I killed 4 Helwyrs, limb tearing hunters.","text":"I killed 4 Helwyrs."}
{"date":"01-Jan-1970 01:00","details":"I killed 3 Gregorovics, all blade wielding terrors.","text":"I killed 3 Gregorovics."}
`)

	err := IterateActivities(in, &test)

	assert.Equal(t, test.err, err)
	assert.Equal(t, 2, test.count)
}
