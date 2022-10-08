package specs_test

import (
	"github.com/alexsatch/ddd/specs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Gender byte

const (
	Male Gender = iota
	Female
)

type Employee struct {
	Age    int
	Gender Gender
}

func TestEval(t *testing.T) {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	t.Run("term", func(t *testing.T) {
		sp := specs.New(Employee.IsLegalAge)
		require.True(t, sp.Eval(bob))
		require.False(t, sp.Eval(alice))
	})

	t.Run("not", func(t *testing.T) {
		notLegalAge := specs.Not[Employee](specs.New(Employee.IsLegalAge))

		require.True(t, notLegalAge.Eval(alice))
	})

	t.Run("and", func(t *testing.T) {
		var (
			isLegalAge = specs.New(Employee.IsLegalAge)
			isMale     = specs.New(Employee.IsMale)
			isFemale   = specs.New(func(e Employee) bool { return e.Gender == Female })
		)

		require.True(t, isLegalAge.Eval(bob))
		require.True(t, isMale.Eval(bob))
		require.False(t, isFemale.Eval(bob))

		assert.True(t, specs.All[Employee](isLegalAge, isMale).Eval(bob))
		assert.False(t, specs.All[Employee](isLegalAge, isFemale).Eval(bob))
	})

	t.Run("or", func(t *testing.T) {
		var (
			isLegalAge = specs.New(Employee.IsLegalAge)
			isFemale   = specs.New(Employee.IsMale)
		)

		require.True(t, specs.Any[Employee](isLegalAge, isFemale).Eval(bob))
	})

}

func (e Employee) IsLegalAge() bool { return e.Age > 18 }
func (e Employee) IsMale() bool     { return e.Gender == Male }
