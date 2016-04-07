//go:generate hel
package datastructures_test

import (
	"hermes/internal/datastructures"
	"sync"
	"unsafe"

	. "github.com/apoydence/eachers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type someInterface interface {
}

type someStruct struct {
}

var _ = Describe("LinkedList", func() {
	var (
		callbackValues chan someInterface
		list           *datastructures.LinkedList

		emitter1 someInterface
		emitter2 someInterface
		emitter3 someInterface
		emitter4 someInterface
	)

	var callback = func(value unsafe.Pointer) {
		casted := *(*someInterface)(value)
		callbackValues <- casted
	}

	BeforeEach(func() {
		callbackValues = make(chan someInterface, 100)
		list = datastructures.NewLinkedList()

		emitter1 = new(someStruct)
		emitter2 = new(someStruct)
		emitter3 = new(someStruct)
		emitter4 = new(someStruct)
	})

	Describe("Traverse()", func() {
		Describe("Append()", func() {
			Context("empty list", func() {
				It("does not invoke the callback", func() {
					list.Traverse(callback)

					Expect(callbackValues).To(BeEmpty())
				})

				Context("single entry", func() {
					var (
						expectedValue someInterface
					)

					BeforeEach(func() {
						expectedValue = emitter1
						list.Append(unsafe.Pointer(&expectedValue))
					})

					It("invokes the callback only once", func() {
						list.Traverse(callback)

						Expect(callbackValues).To(HaveLen(1))
					})

					It("invokes the callback with the expected value", func() {
						list.Traverse(callback)

						Expect(callbackValues).To(Receive(Equal(expectedValue)))
					})

					Context("multiple entries", func() {
						var (
							expectedValues []someInterface
						)

						BeforeEach(func() {
							expectedValues = []someInterface{emitter2, emitter3, emitter4}
							for _, value := range expectedValues {
								list.Append(unsafe.Pointer(&value))
							}
						})

						It("invokes the callback for each entry", func() {
							list.Traverse(callback)

							Expect(callbackValues).To(HaveLen(len(expectedValues) + 1))
						})

						It("invokes the callback with the expected value", func() {
							list.Traverse(callback)

							Expect(callbackValues).To(EqualEach(emitter1, emitter2, emitter3, emitter4))
						})
					})
				})
			})
		})

		// Describe("Remove()", func() {
		// 	Context("empty list", func() {
		// 		It("does not panic", func() {
		// 			f := func() { list.Remove(99) }
		// 			Expect(f).ToNot(Panic())
		// 		})

		// 		Context("single entry", func() {
		// 			var (
		// 				expectedValue Emitter
		// 			)

		// 			BeforeEach(func() {
		// 				expectedValue = 99
		// 				list.Append(expectedValue)
		// 			})

		// 			It("removes the root", func() {
		// 				list.Remove(expectedValue)
		// 				list.Traverse(callback)

		// 				Expect(callbackValues).To(BeEmpty())
		// 			})

		// 			It("does not remove invalid value", func() {
		// 				list.Remove(expectedValue + 1)
		// 				list.Traverse(callback)

		// 				Expect(callbackValues).ToNot(BeEmpty())
		// 			})

		// 			Context("multiple entries", func() {
		// 				var (
		// 					expectedValues []Emitter
		// 				)

		// 				BeforeEach(func() {
		// 					expectedValues = []Emitter{101, 103, 105}
		// 					for _, value := range expectedValues {
		// 						list.Append(value)
		// 					}
		// 				})

		// 				Context("removes the root", func() {
		// 					BeforeEach(func() {
		// 						list.Remove(expectedValue)
		// 					})

		// 					JustBeforeEach(func() {
		// 						list.Traverse(callback)
		// 					})

		// 					It("has the expected number of values", func() {
		// 						Expect(callbackValues).To(HaveLen(len(expectedValues)))
		// 					})

		// 					It("has the expected values", func() {
		// 						Expect(callbackValues).To(EqualEach(101, 103, 105))
		// 					})

		// 					Context("additional append", func() {
		// 						var (
		// 							newValue Emitter
		// 						)

		// 						BeforeEach(func() {
		// 							newValue = 107
		// 							list.Append(newValue)
		// 						})

		// 						It("has the expected number of values", func() {
		// 							Expect(callbackValues).To(HaveLen(len(expectedValues) + 1))
		// 						})

		// 						It("has the expected values", func() {
		// 							Expect(callbackValues).To(EqualEach(101, 103, 105, 107))
		// 						})
		// 					})
		// 				})

		// 				Context("removes a non-root value", func() {
		// 					BeforeEach(func() {
		// 						list.Remove(expectedValues[0])
		// 					})

		// 					JustBeforeEach(func() {
		// 						list.Traverse(callback)
		// 					})

		// 					It("has the expected number of values", func() {
		// 						Expect(callbackValues).To(HaveLen(len(expectedValues)))
		// 					})

		// 					It("has the expected values", func() {
		// 						Expect(callbackValues).To(EqualEach(99, 103, 105))
		// 					})

		// 					Context("additional append", func() {
		// 						var (
		// 							newValue Emitter
		// 						)

		// 						BeforeEach(func() {
		// 							newValue = 107
		// 							list.Append(newValue)
		// 						})

		// 						It("has the expected number of values", func() {
		// 							Expect(callbackValues).To(HaveLen(len(expectedValues) + 1))
		// 						})

		// 						It("has the expected values", func() {
		// 							Expect(callbackValues).To(EqualEach(99, 103, 105, 107))
		// 						})
		// 					})
		// 				})
		// 			})
		// 		})
		// 	})
		// })
	})

	Describe("multiple go-routines", func() {
		var (
			wg    sync.WaitGroup
			count int
		)

		var write = func() {
			emitters := []someInterface{
				emitter1,
				emitter2,
				emitter3,
				emitter4,
			}
			defer wg.Done()
			for i := 0; i < count; i++ {
				value := emitters[i%len(emitters)]
				list.Append(unsafe.Pointer(&value))
			}
		}

		var remove = func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				//				list.Remove(i)
			}
		}

		var read = func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				list.Traverse(func(unsafe.Pointer) {})
			}
		}

		BeforeEach(func() {
			count = 10
			wg.Add(3)
			defer wg.Wait()
			go write()
			go remove()
			go read()
		})

		It("survives the race detector", func() {
			// NOP
		})
	})
})
