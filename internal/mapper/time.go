package mapper

import "time"

const dateLayout = "2006-01-02"

func DateStringPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(dateLayout)
}

func DateTimeString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func ParseDate(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	t, err := time.Parse(dateLayout, input)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ParseDateOrNow(input string) (time.Time, error) {
	if input == "" {
		return time.Now().UTC(), nil
	}
	t, err := time.Parse(dateLayout, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
