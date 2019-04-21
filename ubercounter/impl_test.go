package nanocounter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("Impl", func() {

	const bigNumber = 1000000

	var i *counter

	BeforeEach(func() {

		// create new instance
		i = &counter{
			WindowSize: time.Millisecond * 100,
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
		}, "101ms", "1ms").Should(Equal(0))

	})

	Measure("Benchmarks of basic functionality", func(b Benchmarker) {

		b.Time("Lots of calls to insert", func() {
			for j := 0; j < bigNumber; j++ {
				i.Incr()
			}
		})

		b.Time("Lots of calls to count", func() {
			for j := 0; j < bigNumber; j++ {
				i.Count()
			}
		})

	}, 10)

	Measure("Benchmarks of i/o", func(b Benchmarker) {

		i.ProbeVals = make([]timestamp, 0, bigNumber)
		for j := 0; j < bigNumber; j++ {
			i.ProbeVals = append(i.ProbeVals, Now())
		}

		b.Time("Lots of calls to save", func() {
			err := i.Save("testdata/write_test.txt")
			Ω(err).ToNot(HaveOccurred())
		})

		b.Time("Lots of calls to load", func() {
			var ii = &counter{
				WindowSize: time.Millisecond * 100,
			}
			err := ii.Load("testdata/write_test.txt")
			Ω(err).ToNot(HaveOccurred())
		})

	}, 100)

	It("Saving a counter into a file and then restoring it", func() {

		By("Setting large window size")
		i.WindowSize = time.Second
		var saveloadFname = "testdata/write_test.txt"

		By("Some incrementing")
		for j := 0; j < bigNumber; j++ {
			i.Incr()
		}
		Ω(i.Count()).ToNot(BeZero())

		By("Saving file")
		err := i.Save("testdata/write_test.txt")
		Ω(err).ToNot(HaveOccurred())

		By("Creating another instance")
		var ii = &counter{
			WindowSize: i.WindowSize,
		}

		By("Restoring values")
		err = ii.Load(saveloadFname)
		Ω(err).ToNot(HaveOccurred())

		By("Comparing old and new counters")
		Consistently(func() bool {
			// compare that counters decay in the same way
			// @TODO shitty test
			Ω(ii.Count()).To(Equal(i.Count()))
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
