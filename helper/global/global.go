package global

import (
	"crypto/rand"
	"fmt"
	"html"
	"io"
	"marketplace-svc/helper/message"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
)

func HtmlEscape(req interface{}) {
	value := reflect.ValueOf(req).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}

		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))
	}
}

func NormalizePhone(phone string) (string, message.Message) {
	phone = strings.Replace(phone, "+", "", 1)

	regex := regexp.MustCompile(`^(\d)+$`)
	result := regex.MatchString(phone)
	if !result {
		return "", message.InvalidPhoneFormat
	}

	regex = regexp.MustCompile(`^((0|62)[0-9]+$)`)
	result = regex.MatchString(phone)
	if !result {
		phone = "62" + phone
	}

	regex = regexp.MustCompile(`^0`)
	Str := "${1}62$2"
	phone = regex.ReplaceAllString(phone, Str)

	return phone, message.SuccessMsg
}

func GenerateNumberWithLength(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func NewBoolTrue() *bool {
	b := true
	return &b
}

func NewBoolFalse() *bool {
	b := false
	return &b
}

func GenerateUid(length int) (string, error) {
	return gonanoid.Generate(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890",
		length,
	)
}

func BcryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ParseToIndonesianDuration(duration time.Duration) (result string) {
	result = duration.Round(time.Second).String()
	result = strings.Replace(result, "s", " Detik", 1)
	result = strings.Replace(result, "m", " Menit ", 1)
	result = strings.Replace(result, "h", " Jam ", 1)

	return
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	case reflect.Struct:
		return reflect.ValueOf(i).IsZero()
	}
	return false
}

func FormatDate(date int64) string {
	createdAtTime := time.Unix(date, 0)
	createdAtTimeStr := createdAtTime.Format(time.RFC3339)
	return createdAtTimeStr
}

func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	case duration < time.Hour:
		time := duration.Minutes()
		s := fmt.Sprintf("%f", time)
		minute := strings.Split(s, ".")
		minuteInt, _ := strconv.ParseFloat(minute[0], 32)
		timeRest := time - minuteInt
		second := int(timeRest * 60)
		return fmt.Sprintf("%sm %ds", minute[0], second)
	case duration < 24*time.Hour:
		time := duration.Hours()
		s := fmt.Sprintf("%f", time)
		hour := strings.Split(s, ".")
		hourInt, _ := strconv.ParseFloat(hour[0], 32)
		timeRest := time - hourInt
		minute := int(timeRest * 60)
		return fmt.Sprintf("%sh %dm", hour[0], minute)
	case duration < 7*24*time.Hour:
		time := duration.Hours() / 24
		s := fmt.Sprintf("%f", time)
		day := strings.Split(s, ".")
		// dayInt, _ := strconv.ParseFloat(day[0], 32)
		// timeRest := time - dayInt
		// hour := int(timeRest * 24)
		return fmt.Sprintf("%sday", day[0])
	case duration < 30*24*time.Hour:
		time := duration.Hours() / (24 * 7)
		s := fmt.Sprintf("%f", time)
		week := strings.Split(s, ".")
		// weekInt, _ := strconv.ParseFloat(week[0], 32)
		// timeRest := time - weekInt
		// day := int(timeRest * 7)
		return fmt.Sprintf("%sweek", week[0])
	case duration < 365*24*time.Hour:
		time := duration.Hours() / (24 * 30)
		s := fmt.Sprintf("%f", time)
		month := strings.Split(s, ".")
		// monthInt, _ := strconv.ParseFloat(month[0], 32)
		// timeRest := time - monthInt
		// week := int(timeRest * 4)
		return fmt.Sprintf("%smonth", month[0])
	default:
		return fmt.Sprintf("%d year", int(duration.Hours()/(24*365)))
	}
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64, unit string) float64 {
	// Earth radius in kilometers by default
	var earthRadius float64

	switch unit {
	case "km":
		earthRadius = 6371.0
	case "mi":
		earthRadius = 3959.0
	case "m":
		earthRadius = 6371000.0
		// 1 meter = 0.000621371 miles
	}

	// Convert latitude and longitude from degrees to radians
	lat1Rad := toRadians(lat1)
	lon1Rad := toRadians(lon1)
	lat2Rad := toRadians(lat2)
	lon2Rad := toRadians(lon2)

	// Calculate differences
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Pow(math.Sin(deltaLat/2), 2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in the specified unit
	distance := earthRadius * c

	return distance
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
