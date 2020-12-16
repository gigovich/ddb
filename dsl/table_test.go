package dsl

import (
	"testing"
)

func TestSchemeDefault(t *testing.T) {
	type User struct {
		ID    int
		Email string
	}

	t.Run("obj should be pointer", func(t *testing.T) {
		user := User{}
		if err := Default(user).Err(); err != ErrObjNotPointer {
			t.Fail()
		}
	})

	t.Run("field should be pointer", func(t *testing.T) {
		user := User{}
		if err := Default(&user).WithFieldName(user.Email, "user_email").Err(); err != ErrFieldNotPointer {
			t.Fail()
		}
	})

	t.Run("check mapped", func(t *testing.T) {
		user := User{}
		s := Default(&user).WithFieldName(&user.Email, "user_email")
		if err := s.Err(); err != nil {
			t.Error(err)
		}
		if s.Remaps["Email"] != "user_email" {
			t.Errorf("not remaped")
		}
	})

	t.Run("no such field error", func(t *testing.T) {
		i := 1
		user := User{}
		s := Default(&user).WithPK(&i)
		if err := s.Err(); err != ErrNoSuchField {
			t.Fail()
		}
	})

	t.Run("check PK", func(t *testing.T) {
		user := User{}
		s := Default(&user).WithPK(&user.Email)
		if err := s.Err(); err != nil {
			t.Error(err)
			return
		}

		if len(s.PK) != 1 || s.PK[0] != "Email" {
			t.Errorf("pk was not redefined")
			return
		}
	})

	t.Run("getter field should be pointer", func(t *testing.T) {
		user := User{}
		if err := Default(&user).WithGetter(user.Email).Err(); err != ErrFieldNotPointer {
			t.Fail()
		}
	})

	t.Run("getter should be added", func(t *testing.T) {
		user := User{}
		s := Default(&user).WithGetter(&user.Email)
		if err := s.Err(); err != nil {
			t.Error(err)
			return
		}

		if _, ok := s.Getters["Email"]; !ok {
			t.Fail()
		}
	})
}
