package server

import (
	"github.com/stretchr/testify/mock"
)

type MockControllers struct {
	mock.Mock
}

func (m *MockControllers) SomeMethod() {
	// Mock implementation
}

//
// func TestNewServer(t *testing.T) {
// 	mockControllers := new(MockControllers)
// 	mockControllers.On("SomeMethod").Return()
//
// 	server := NewServer(
// 		WithPort("8080"),
// 		WithControllers(mockControllers),
// 	)
//
// 	// Check if the server has the expected properties
// 	assert.Equal(t, "8080", server.port)
// 	assert.NotNil(t, server.controllers)
// 	assert.Implements(t, (*controllers.Controllers)(nil), server.controllers)
//
// 	// Check if routes were set up
// 	// This could be a bit tricky; you might need to check if the routes are registered correctly
// 	// You might need to check internal router configurations or use an HTTP test server
//
// 	// Example: Ensure router is initialized (not directly testable without reflection)
// 	assert.NotNil(t, server.router)
// }
