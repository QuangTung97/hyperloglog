package hll

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"os"
	"sort"
	"testing"
)

func newCounter() *Counter {
	var c Counter
	InitCounter(&c)
	return &c
}

func TestCounter_Zero(t *testing.T) {
	c := newCounter()
	assert.Equal(t, float64(0), c.Estimate())

	c = newCounter()
	c.Add(0x0)
	assert.InDelta(t, 1, c.Estimate(), 0.1)

	c = newCounter()
	c.Add(0x0)
	c.Add(0xffffffffaaaaaaaa)
	assert.InDelta(t, 2, c.Estimate(), 0.1)

	c = newCounter()
	c.Add(0x0)
	c.Add(0xffffffffaaaaaaaa)
	c.Add(0x444fffffaaaaaaaa)
	assert.InDelta(t, 3, c.Estimate(), 0.1)
}

func TestCounter(t *testing.T) {
	c := newCounter()

	const num = 48
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("NUM:%08d", i)
		hash := murmur3.Sum64([]byte(key))
		c.Add(hash)
	}

	e := c.Estimate()
	err := (e - num) / num
	fmt.Println(e, num, err)
}

func TestCounter_Test(t *testing.T) {
	c := newCounter()

	const num = 48
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("KNUM:%08d", i)
		hash := murmur3.Sum64([]byte(key))
		c.Add(hash)
	}

	e := c.Estimate()
	err := (e - num) / num
	fmt.Println(e, num, err)
}

func TestCounter_Test_2(t *testing.T) {
	c := newCounter()

	const num = 1293
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("KNUM:%08d", i)
		hash := murmur3.Sum64([]byte(key))
		c.Add(hash)
	}

	e := c.Estimate()
	err := (e - num) / num
	fmt.Println(e, num, err)
}

func TestCounter_Test_3(t *testing.T) {
	c := newCounter()

	const num = 5 * 64
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("M2NUM:%08d", i)
		hash := murmur3.Sum64([]byte(key))
		c.Add(hash)
	}

	e := c.Estimate()
	err := (e - num) / num
	fmt.Println(e, num, err)
}

func TestCounter_Test_4(t *testing.T) {
	c := newCounter()

	const num = 5 * 64
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("90NUM:%08d", i)
		hash := murmur3.Sum64([]byte(key))
		c.Add(hash)
	}

	e := c.Estimate()
	err := (e - num) / num
	fmt.Println(e, num, err)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var maxValues []float64

func showErrorRange(num int) {
	const loop = 20

	var errorValues []float64

	for i := 0; i < loop; i++ {
		c := newCounter()
		prefix := randString(9)

		for n := 0; n < num; n++ {
			key := prefix + fmt.Sprintf(":%08d", n)
			hash := murmur3.Sum64([]byte(key))
			c.Add(hash)
		}

		err := math.Abs(c.Estimate()-float64(num)) / float64(num)
		errorValues = append(errorValues, err)
	}

	sort.Float64s(errorValues)
	max := errorValues[len(errorValues)-1]
	// fmt.Println("ERROR:", num, errorValues[0], max)
	fmt.Println(max)

	maxValues = append(maxValues, errorValues...)
}

func TestCounter_Random(t *testing.T) {
	rand.Seed(5589)

	maxValues = nil
	for n := 1500; n < 2000; n++ {
		showErrorRange(n)
	}

	file, err := os.Create("output")
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	for k := 0; k < 5000; k++ {
		i := rand.Intn(len(maxValues))
		j := rand.Intn(len(maxValues))
		_, err := fmt.Fprintln(file, (1+maxValues[j])/(1-maxValues[i]))
		if err != nil {
			panic(err)
		}
	}
}

func TestLowerBound(t *testing.T) {
	i := lowerBound([]float64{}, 123)
	assert.Equal(t, 0, i)

	i = lowerBound([]float64{22}, 22)
	assert.Equal(t, 0, i)

	i = lowerBound([]float64{22}, 23)
	assert.Equal(t, 1, i)

	i = lowerBound([]float64{2, 3, 5, 6}, 4)
	assert.Equal(t, 2, i)

	i = lowerBound([]float64{2, 3, 5, 6, 9, 10}, 22)
	assert.Equal(t, 6, i)

	i = lowerBound([]float64{2, 3, 5, 5, 6, 9, 10}, 5)
	assert.Equal(t, 2, i)
}
