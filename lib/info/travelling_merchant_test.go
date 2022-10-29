package info

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestTravellingMerchant(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := os.ReadFile("testdata/travelling_merchant.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := TravellingMerchant(client)

	assert.Nil(t, err)
	assert.Len(t, res.Products, 4)

	assert.Equal(t, res.Products[0].Name, "Uncharted island map")
	assert.Equal(t, res.Products[1].Name, "Gift for the Reaper")
	assert.Equal(t, res.Products[2].Name, "Silverhawk down")
	assert.Equal(t, res.Products[3].Name, "Large goebie burial charm")

	assert.Equal(t, res.Products[0].Cost, 800000)
	assert.Equal(t, res.Products[1].Cost, 1250000)
	assert.Equal(t, res.Products[2].Cost, 1500000)
	assert.Equal(t, res.Products[3].Cost, 150000)
}

func TestTravellingMerchantOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := TravellingMerchant(http.DefaultClient)

	assert.Nil(t, err)
}
