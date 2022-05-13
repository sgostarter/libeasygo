package confuse

import "go.uber.org/atomic"

var (
	confusionCount = atomic.NewInt64(1)
)

func NewRandConfuse2() Confuse {
	return &randConfuse2Impl{}
}

type randConfuse2Impl struct {
}

func (impl *randConfuse2Impl) Seal(d []byte) ([]byte, error) {
	confusionData := impl.buildConfusionData()

	var result []byte

	for _, v := range d {
		result = append(result, confusionData[0])
		confusionData = append(confusionData[1:], confusionData[0])
		result = append(result, v)
	}

	return result, nil
}

func (impl *randConfuse2Impl) Open(d []byte) ([]byte, error) {
	return impl.separateConfusionData(d), nil
}

func (impl *randConfuse2Impl) buildConfusionData() []byte {
	count := int(confusionCount.Inc())
	number := count % 356

	if number < 10 {
		number = 10
	}

	var data = make([]byte, number)

	for i := 0; i < number; i++ {
		count = int(confusionCount.Inc())
		data[i] = uint8((count % 254) + 1)
	}

	return data
}

func (impl *randConfuse2Impl) separateConfusionData(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	var result = make([]byte, 0, len(data)/2)

	for index, v := range data {
		if index%2 == 0 {
			continue
		}

		result = append(result, v)
	}

	return result
}
