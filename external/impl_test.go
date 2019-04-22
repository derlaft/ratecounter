package counter

import (
	"github.com/derlaft/ratecounter/iface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("Impl", func() {

	const bigNumber = 1000000

	const testFilename = "testdata/state.json"

	var i iface.Counter

	BeforeEach(func() {

		// create new instance
		var createErr error
		i, createErr = NewCounter(time.Millisecond * 100)
		Ω(createErr).ToNot(HaveOccurred())
	})

	It("It works on basic values", func() {

		By("Just started, it should return zero")
		Ω(i.Count()).To(Equal(0))

		By("After one increment, it should return one")
		i.Incr()
		Ω(i.Count()).To(Equal(1))

		By("Eventually, that one value should expire")
		Eventually(func() int {
			return i.Count()
		}, "120ms", "5ms").Should(Equal(0))

	})

	Measure("Benchmarks of basic functionality", func(b Benchmarker) {

		b.Time("1kk calls to insert", func() {
			for j := 0; j < bigNumber; j++ {
				i.Incr()
			}
		})

		b.Time("1kk calls to count", func() {
			for j := 0; j < bigNumber; j++ {
				i.Count()
			}
		})

	}, 10)

	It("Trying to detect some races", func() {

		var wg sync.WaitGroup

		wg.Add(2)

		go func() {
			defer GinkgoRecover()

			for j := 0; j < bigNumber; j++ {
				i.Incr()
				i.Count()
			}

			wg.Done()
		}()

		go func() {
			defer GinkgoRecover()

			for j := 0; j < bigNumber; j++ {
				i.Count()
				i.Incr()
			}

			wg.Done()
		}()

		wg.Wait()

	})

})
