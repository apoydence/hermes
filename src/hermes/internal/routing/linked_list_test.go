//go:generate hel
package routing_test

import (
	"hermes/internal/routing"
	"sync"

	. "github.com/apoydence/eachers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LinkedList", func() {
	var (
		callbackValues chan routing.Emitter
		list           *routing.LinkedList

		emitter1 *mockEmitter
		emitter2 *mockEmitter
		emitter3 *mockEmitter
		emitter4 *mockEmitter
	)

	var callback = func(value routing.Emitter) {
		callbackValues <- value
	}

	BeforeEach(func() {
		callbackValues = make(chan routing.Emitter, 100)
		list = routing.NewLinkedList()

		emitter1 = newMockEmitter()
		emitter2 = newMockEmitter()
		emitter3 = newMockEmitter()
		emitter4 = newMockEmitter()
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
						expectedValue routing.Emitter
					)

					BeforeEach(func() {
						expectedValue = emitter1
						list.Append(expectedValue)
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
							expectedValues []routing.Emitter
						)

						BeforeEach(func() {
							expectedValues = []routing.Emitter{emitter2, emitter3, emitter4}
							for _, value := range expectedValues {
								list.Append(value)
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
			emitters := []routing.Emitter{
				emitter1,
				emitter2,
				emitter3,
				emitter4,
			}
			defer wg.Done()
			for i := 0; i < count; i++ {
				list.Append(emitters[i%len(emitters)])
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
				list.Traverse(func(routing.Emitter) {})
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
