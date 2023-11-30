package unfoldcpp

// test for Unfold()
import (
	"testing"
)

func TestUnfold(t *testing.T) {
	// test
	path := "/home/comavius/projects/comavius-cpp-library/main.cpp"
	code, err := Unfold(path)
	if err != nil {
		t.Error(err)
	}
	// stdout code
	t.Log(code)
}
