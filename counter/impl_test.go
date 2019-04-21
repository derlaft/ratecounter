package counter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("Impl", func() {

	const bigNumber = 1000000

	const testFilename = "testdata/state.json"

	var i *counter

	BeforeEach(func() {

		// create new instance
		i = &counter{
			WindowSize: time.Millisecond * 100,
			Accuracy:   time.Millisecond * 10,
		}
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

	Measure("Benchmarks of i/o", func(b Benchmarker) {

		b.Time("Save performance", func() {
			err := i.Save("testdata/write_test.json")
			Ω(err).ToNot(HaveOccurred())
		})

		b.Time("Load performance", func() {
			var ii = &counter{
				WindowSize: time.Millisecond * 100,
				Accuracy:   time.Millisecond * 10,
			}
			err := ii.Load("testdata/write_test.json")
			Ω(err).ToNot(HaveOccurred())
		})

	}, 100)

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
		err := i.Save(testFilename)
		Ω(err).ToNot(HaveOccurred())

		By("Creating another instance")
		var ii = &counter{
			WindowSize: i.WindowSize,
			Accuracy:   i.Accuracy,
		}

		By("Restoring values")
		err = ii.Load(testFilename)
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
