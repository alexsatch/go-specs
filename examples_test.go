package specs_test

import (
	"fmt"

	"github.com/alexsatch/go-specs"
)

func ExampleSpec_Eval() {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	sp := specs.New(Employee.IsLegalAge)

	fmt.Println("bob:", sp.Eval(bob))
	fmt.Println("alice:", sp.Eval(alice))

	// Output:
	// bob: true
	// alice: false
}

func ExampleSpec_AndNotFunc() {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	sp := specs.New(Employee.IsLegalAge).
		AndNotFunc(Employee.IsMale)

	fmt.Println("spec:", sp.String())
	fmt.Println("bob:", sp.Eval(bob))
	fmt.Println("alice:", sp.Eval(alice))

	// Output:
	// spec: all(.IsLegalAge, not(.IsMale))
	// bob: false
	// alice: false
}

func ExampleSpec_AndFunc() {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	sp := specs.New(Employee.IsLegalAge).
		AndFunc(Employee.IsMale)

	fmt.Println("spec:", sp.String())
	fmt.Println("bob:", sp.Eval(bob))
	fmt.Println("alice:", sp.Eval(alice))

	// Output:
	// spec: all(.IsLegalAge, .IsMale)
	// bob: true
	// alice: false
}
