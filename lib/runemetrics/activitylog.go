package runemetrics

import (
	"encoding/json"
	"io"
)

type ActivityIterator interface {
	HandleActivity(activity Activity) error
}

func WriteActivities(w io.Writer, activities []Activity) error {
	encoder := json.NewEncoder(w)

	for _, activity := range activities {
		err := encoder.Encode(activity)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadActivities(r io.Reader) ([]Activity, error) {
	out := []Activity{}
	decoder := json.NewDecoder(r)

	for decoder.More() {
		activity := &Activity{}
		err := decoder.Decode(activity)
		if err != nil {
			if err == io.EOF {
				return out, nil
			}
			return nil, err
		}
		out = append(out, *activity)
	}

	return out, nil
}
