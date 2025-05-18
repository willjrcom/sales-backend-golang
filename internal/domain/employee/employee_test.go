package employeeentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewEmployee(t *testing.T) {
   userID := uuid.New()
   emp := NewEmployee(userID)
   assert.NotEqual(t, uuid.Nil, emp.ID)
   assert.Equal(t, userID, emp.UserID)
   assert.Nil(t, emp.User)
}