package mwapi_test

import "mwapi"
import "testing"

func TestMWAPiStruct(*testing.T) {
	mwapi := mwapi.MWApi{"https", "commons.wikimeida.org", "/w/api.php"}
	fmt.Println(1)
	fmt.Println(mwapi)
}
