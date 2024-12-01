package testsuite

type MetaData map[string]interface{}

type TestSuite struct {
	// stores useful global data for tests
	MetaData MetaData
}

func NewTestSuite() *TestSuite {
	return &TestSuite{
		MetaData: make(map[string]interface{}),
	}
}
