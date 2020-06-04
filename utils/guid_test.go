package utils_test

import (
	"./"
	"testing"
)

func TestGetGUID(t *testing.T) {

	guid := utils.GetGlobalId()

	var want uint64 = 1

	if want != guid {
		t.Errorf("GetGlobalId failed,%d != %d", want, guid)
	}
	want = 3
	guid = utils.GetGlobalId()
	guid = utils.GetGlobalId()
	if want != guid {
		t.Errorf("GetGlobalId failed,%d != %d", want, guid)
	}

	want = 3
	utils.ReStoreGlobalId(guid)
	guid = utils.GetGlobalId()
	if want != guid {
		t.Errorf("GetGlobalId failed,%d != %d", want, guid)
	}

}
