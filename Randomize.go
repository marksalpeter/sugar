package sugar

import (
	"math"
	"math/rand"
	"reflect"
	"time"
)

var (
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lettersLen = len(letters)
)

var timeType = reflect.TypeOf(time.Time{})

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Randomize populates any interface with random data, optionally excluding different types of data passed in
//
// Example
// This is an example of how to use `Randomize` in a api endpoint testing scenario
//
//   type Model struct {
//   	ID  uint
//   	Foo string
//   	Bar []byte
//   }
//
//   func TestModel(t *testing.T) {
//   	s := sugar.New(t)
//
//   	var model models.Model
//
//   	s.Must("add models to the database", func(_ sugar.Log) bool {
//   		// make a radom test model
//   		sugar.Randomize(&model)
//
//   		// add it to the database
//   		// ...
//
//   		// start the api server
//   		// ...
//
//   	})
//
//      s.Assert("/model/:id returns the correct model", func(log sugar.Log) bool {
//          // query the endpoint and decode the response
//          var responseModel Model
//          if request, err := http.NewRequest("GET", fmt.Sprintf("url/to/api/server/model/%d", testModel.ID), nil); err != nil {
//              log(err)
//              return false
//          } else if response, err := http.DefaultClient.Do(request); err != nil {
//              log(err)
//              return false
//          } else if err := json.NewDecoder(response.Body).Decode(&responseModel); err != nil {
//              log(err)
//              return false
//          }
//
//          // compare the result with the model we added to the database
//          return log.Compare(model, responseModel)
//      })
//
//   }
func Randomize(i interface{}, excluding ...interface{}) {

	iValue, ok := i.(reflect.Value)
	if !ok {
		iValue = reflect.ValueOf(i)
	}
	iType := iValue.Type()
	iKind := iType.Kind()

	if excluding != nil && len(excluding) > 0 {
		for _, omit := range excluding {
			if iType == reflect.TypeOf(omit) {
				return
			}
		}
	}

	// support for random times
	if iType == timeType {
		iValue.Set(reflect.ValueOf(time.Unix(
			int64(rand.Intn(math.MaxInt32)),
			0,
		)))
		return
	}

	switch iKind {
	case reflect.Ptr, reflect.Interface:
		if !iValue.IsNil() {
			Randomize(iValue.Elem(), excluding...)
		}
	case reflect.Slice, reflect.Array:
		for i, l := 0, iValue.Len(); i < l; i++ {
			Randomize(iValue.Index(i), excluding...)
		}
	case reflect.Struct:
		for i, l := 0, iValue.NumField(); i < l; i++ {
			Randomize(iValue.Field(i), excluding...)
		}
	case reflect.String:
		if iValue.CanSet() {
			runeArray := make([]rune, rand.Intn(255)) // max varchar is 255
			for i := range runeArray {
				rand.Seed(time.Now().UTC().UnixNano())
				runeArray[i] = letters[rand.Intn(lettersLen)]
			}
			iValue.SetString(string(runeArray))
		}
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		if iValue.CanSet() {
			iValue.SetInt(int64(rand.Intn(math.MaxInt32)))
		}
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		if iValue.CanSet() {
			iValue.SetUint(uint64(rand.Intn(math.MaxInt32)))
		}
	case reflect.Float64, reflect.Float32:
		if iValue.CanSet() {
			iValue.SetFloat(float64(rand.Intn(math.MaxInt32)))
		}
	case reflect.Bool:
		if iValue.CanSet() {
			iValue.SetBool(rand.Intn(math.MaxInt64)%2 == 1)
		}
	}
}
