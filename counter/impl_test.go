package counter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("Impl", func() {

	const bigNumber = 1000000

	var i *counter
	var fact *factory

	BeforeEach(func() {

		// create new instance
		i = &counter{
			WindowSize: time.Millisecond * 100,
			Accuracy:   time.Millisecond * 10,
		}

		fact = &factory{}
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

	It("Saving a counter into a file and then restoring it", func() {

		By("Setting large window size")
		i.WindowSize = time.Millisecond * 800
		i.Accuracy = time.Millisecond * 200

		By("Some incrementing")
		for j := 0; j < bigNumber; j++ {
			i.Incr()
		}
		Ω(i.Count()).ToNot(BeZero())

		By("Saving file")
		data, err := i.Save()
		Ω(err).ToNot(HaveOccurred())

		By("Loading an instance using factory")
		ii, err := fact.Load(i.WindowSize, i.Accuracy, data)
		Ω(err).ToNot(HaveOccurred())

		By("Comparing old and new counters")
		Consistently(func() bool {
			// compare that counters decay in the same way
			// @TODO shitty test
			return ii.Count() == i.Count()
		}).Should(BeTrue())

	})

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
