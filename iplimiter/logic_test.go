package iplimiter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"time"

	"github.com/derlaft/ratecounter/iface"
	"github.com/derlaft/ratecounter/mocks"
	"github.com/golang/mock/gomock"
)

var _ = Describe("Logic", func() {

	var (
		sampleIP      = net.ParseIP("8.8.8.8")
		anotherIP     = net.ParseIP("1.1.1.1")
		yetAnotherIP  = net.ParseIP("127.0.0.1")
		_             = yetAnotherIP
		_             = anotherIP
		testInterval  = time.Second
		testAccuracy  = time.Millisecond * 100
		testNumPerIp  = 40
		sameNumber    = 12309
		encodedValue1 = []byte(`"counter1"`)
		encodedValue2 = []byte(`"counter2"`)
		encodedValue3 = []byte(`"counter3"`)
		pv1           = `{"Global":"counter3","Values":{"1.1.1.1":"counter2","8.8.8.8":"counter1"}}`
		pv2           = `{"Global":"counter3","Values":{"8.8.8.8":"counter1","1.1.1.1":"counter2"}}`
	)

	var (
		mockCtl *gomock.Controller

		// general request counter
		mockCounter        *counter_mocks.MockCounter
		mockAnotherCounter *counter_mocks.MockCounter
		mockGlobalCounter  *counter_mocks.MockCounter
		mockCounterFactory *counter_mocks.MockCounterFactory

		c  *limiterImpl
		lf *limiterFactoryImpl
	)

	BeforeEach(func() {
		// initialize mocking
		mockCtl = gomock.NewController(GinkgoT())

		mockCounter = counter_mocks.NewMockCounter(mockCtl)
		mockAnotherCounter = counter_mocks.NewMockCounter(mockCtl)
		mockGlobalCounter = counter_mocks.NewMockCounter(mockCtl)

		mockCounterFactory = counter_mocks.NewMockCounterFactory(mockCtl)

		c = &limiterImpl{
			Interval:       testInterval,
			Accuracy:       testAccuracy,
			MaxNumberPerIp: testNumPerIp,
			Counters:       make(map[string]iface.Counter),

			CounterFactory: mockCounterFactory,
			GlobalCounter:  mockGlobalCounter,
		}

		lf = &limiterFactoryImpl{
			Interval:       testInterval,
			Accuracy:       testAccuracy,
			MaxNumberPerIp: testNumPerIp,
			CounterFactory: mockCounterFactory,
		}
	})

	AfterEach(func() {
		mockCtl.Finish()
	})

	It("Proper calling of TotalRequests", func() {

		mockGlobalCounter.EXPECT().Count().Return(sameNumber)
		count := c.TotalRequests()
		Ω(count).To(Equal(sameNumber))
	})

	It("Some calls of OnRequest", func() {

		// only one request is going to be allowed during this test
		mockGlobalCounter.EXPECT().Incr().Times(1)

		// a new counter must be created
		mockCounterFactory.
			EXPECT().
			New(testInterval, testAccuracy).
			Return(mockCounter)

		// and then incremented
		mockCounter.EXPECT().Incr()

		// and then queried for the new value
		mockCounter.EXPECT().Count().Return(40)

		// and then incremented again
		mockCounter.EXPECT().Incr()

		// adn then queried for the new value
		mockCounter.EXPECT().Count().Return(41)

		// do the call
		reject := c.OnRequest(sampleIP)
		Ω(reject).Should(BeFalse())

		// query for the same ip again
		reject = c.OnRequest(sampleIP)
		Ω(reject).Should(BeTrue())
	})

	It("Cleanup test - stale data must be removed", func() {

		c.Counters = map[string]iface.Counter{
			"8.8.8.8": mockCounter,
			"1.1.1.1": mockAnotherCounter,
		}

		mockCounter.EXPECT().Count().Return(1)
		mockAnotherCounter.EXPECT().Count().Return(0)

		c.Cleanup()

		Ω(c.Counters).To(HaveLen(1))
		Ω(c.Counters).To(HaveKey("8.8.8.8"))
		Ω(c.Counters["8.8.8.8"]).To(Equal(mockCounter))

	})

	It("Save test", func() {

		c.Counters = map[string]iface.Counter{
			"8.8.8.8": mockCounter,
			"1.1.1.1": mockAnotherCounter,
		}

		mockCounter.EXPECT().Save().Return(encodedValue1, nil)
		mockAnotherCounter.EXPECT().Save().Return(encodedValue2, nil)
		mockGlobalCounter.EXPECT().Save().Return(encodedValue3, nil)

		result, err := c.SaveState()
		Ω(err).ToNot(HaveOccurred())

		// two different possible values (thanks go map shuffling)
		Ω(string(result)).To(Or(Equal(pv1), Equal(pv2)))

	})

	It("Restore test", func() {

		mockCounterFactory.EXPECT().
			Load(testInterval, testAccuracy, encodedValue1).
			Return(mockCounter, nil)

		mockCounterFactory.EXPECT().
			Load(testInterval, testAccuracy, encodedValue2).
			Return(mockAnotherCounter, nil)

		mockCounterFactory.EXPECT().
			Load(testInterval, testAccuracy, encodedValue3).
			Return(mockGlobalCounter, nil)

		limiter, err := lf.Restore([]byte(pv1))
		Ω(err).ToNot(HaveOccurred())
		Ω(limiter).ToNot(BeNil())

	})

})
