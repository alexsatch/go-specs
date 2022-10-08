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

func (e Employee) IsLegalAge() bool { return e.Age > 18 }
func (e Employee) IsMale() bool     { return e.Gender == Male }

func TestEval(t *testing.T) {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	t.Run("nil", func(t *testing.T) {
		sp := specs.New[Employee](nil)
		require.True(t, sp.Eval(bob))
	})

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

func TestAPI(t *testing.T) {
	bob := Employee{Age: 19, Gender: Male}
	alice := Employee{Age: 17, Gender: Female}

	sp := specs.
		New(Employee.IsLegalAge).
		AndFunc(Employee.IsMale)

	require.True(t, sp.Eval(bob))
	require.False(t, sp.Eval(alice))
}

func is30(e Employee) bool { return e.Age == 30 }

type EmployeeService struct{}

func (es EmployeeService) IsAllowed(e Employee) bool {
	return e.IsLegalAge()
}

func TestString(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		sp := specs.New[Employee](nil)
		require.Equal(t, "nil", sp.String())
	})

	t.Run("term", func(t *testing.T) {
		sp := specs.New(Employee.IsLegalAge)
		require.Equal(t, ".IsLegalAge", sp.String())
	})

	t.Run("anonymous", func(t *testing.T) {
		sp := specs.New(func(e Employee) bool { return e.Gender == Female })
		require.Regexp(t, `<anonymous: (.*)/ddd/specs/spec_test.go:\d+>`, sp.String())
	})

	t.Run("is30", func(t *testing.T) {
		sp := specs.New(is30)
		require.Equal(t, `github.com/alexsatch/ddd/specs_test.is30`, sp.String())
	})

	t.Run("other type", func(t *testing.T) {
		sp := specs.New(EmployeeService{}.IsAllowed)
		require.Equal(t, `github.com/alexsatch/ddd/specs_test.EmployeeService.IsAllowed`, sp.String())
	})

	t.Run("not", func(t *testing.T) {
		sp := specs.Not[Employee](specs.New(Employee.IsLegalAge))
		require.Equal(t, "not(.IsLegalAge)", sp.String())
	})

	t.Run("all", func(t *testing.T) {
		var (
			isLegalAge = specs.New(Employee.IsLegalAge)
			isMale     = specs.New(Employee.IsMale)
			allSpecs   = specs.All[Employee](isLegalAge, isMale)
		)

		require.Equal(t, "all(.IsLegalAge, .IsMale)", allSpecs.String())
	})

	t.Run("any", func(t *testing.T) {
		var (
			isLegalAge = specs.New(Employee.IsLegalAge)
			isMale     = specs.New(Employee.IsMale)
			allSpecs   = specs.Any[Employee](isLegalAge, isMale)
		)

		require.Equal(t, "any(.IsLegalAge, .IsMale)", allSpecs.String())
	})
}
